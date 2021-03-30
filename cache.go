package localcache

import (
	"sync"
	"time"

	"github.com/MoeYang/go-localcache/common"
	"github.com/MoeYang/go-localcache/datastruct/dict"
)

const (
	defaultShardCnt = 256 // ShardCnt must be a power of 2
	defaultCap      = 1024

	defaultTTL             = 60  // default key expire time 60 sec
	defaultTTLTick         = 100 // default time.tick 100ms
	defaultTTLCheckCount   = 100 // every time check 100 keys
	defaultTTLCheckPercent = 25  // every check expierd key > 25, check another time
	defaultTTLCheckRunTime = 50  // max run time for a tick

	hitChanLen = 1 << 15 // 32768
	addChanLen = 1 << 15

	opTypeDel = uint8(1)
	opTypeAdd = uint8(2)
)

type Cache interface {
	// Get a key and return the value and if the key exists
	Get(key string) (interface{}, bool)
	// GetOrLoad get a key, while not exists, call f() to load data
	GetOrLoad(key string, f LoadFunc) (interface{}, error)
	// Set a key-value with default seconds to live
	Set(key string, value interface{})
	// SetWithExpire set a key-value with seconds to live
	SetWithExpire(key string, value interface{}, ttl int64)
	// Del delete key
	Del(key string)
	// Len return count of keys in cache
	Len() int
	// Flush clear all keys in chache, should do this when set and del is stop
	Flush()
	// Stop the cacheProcess by close stopChan
	Stop()
	// Statistic return cache Statistic {"hit":1, "miss":1, "hitRate":50.0}
	Statistic() map[string]interface{}
}

// LoadFunc is called to load data from user storage
type LoadFunc func() (interface{}, error)

type localCache struct {
	// elimination policy of keys
	policy     policy
	policyType string

	// data dict
	dict     dict.Dict
	shardCnt int // shardings count
	cap      int // capacity

	// ttl dict
	ttlDict dict.Dict
	ttl     int64 // Global Keys expire seconds

	hitChan  chan interface{} // chan while get a key should put in
	opChan   chan opMsg       // add del and add msg in one chan, so we can do options order by time acs
	stopChan chan struct{}    // chan stop signal

	// cache statist
	statist statist

	// group singleFlight
	group common.Group
}

// NewLocalCache return Cache obj with options
func NewLocalCache(options ...Option) Cache {
	c := &localCache{
		shardCnt: defaultShardCnt,
		cap:      defaultCap,
		ttl:      defaultTTL,
		hitChan:  make(chan interface{}, hitChanLen),
		opChan:   make(chan opMsg, addChanLen),
		statist:  newstatisCaculator(false),
	}
	// set options
	for _, opt := range options {
		opt(c)
	}
	// init dict
	c.dict = dict.NewDict(c.shardCnt)
	// init ttl dict
	c.ttlDict = dict.NewDict(c.shardCnt)
	// init policy
	c.policy = newPolicy(c.policyType, c.cap, c)
	// start goroutine
	c.start()

	return c
}

type Option func(*localCache)

// WithGlobalTTL set all keys default expire time of seconds
func WithGlobalTTL(expireSecond int64) Option {
	if expireSecond <= 0 {
		expireSecond = defaultTTL
	}
	return func(c *localCache) {
		c.ttl = expireSecond
	}
}

// WithShardCount set max Capacity
func WithCapacity(cap int) Option {
	if cap <= 0 {
		cap = defaultCap
	}
	return func(c *localCache) {
		c.cap = cap
	}
}

// WithShardCount shardCnt must be a power of 2
func WithShardCount(shardCnt int) Option {
	if shardCnt <= 0 {
		shardCnt = defaultShardCnt
	}
	return func(c *localCache) {
		c.shardCnt = shardCnt
	}
}

// WithPolicy set the elimination policy of keys
func WithPolicy(policyType string) Option {
	return func(c *localCache) {
		c.policyType = policyType
	}
}

// WithStatist set whether need to caculate the cache`s statist, default false.
//  not need may led performance a very little better ^-^
func WithStatist(needStatistic bool) Option {
	return func(c *localCache) {
		c.statist = newstatisCaculator(needStatistic)
	}
}

func (l *localCache) Get(key string) (interface{}, bool) {
	obj, has := l.dict.Get(key)
	if has {
		element := l.policy.unpack(obj)
		element.lock.RLock()
		value := element.value
		isExpire := element.isExpire()
		element.lock.RUnlock()
		if !isExpire {
			// add hit count, if chan full, skip this signal is ok
			select {
			case l.hitChan <- obj:
			default:
			}
			l.statist.hitIncr()
			return value, true
		} else {
			l.Del(key)
		}
	}
	// not exists or expired
	l.statist.missIncr()
	return nil, false
}

func (l *localCache) GetOrLoad(key string, f LoadFunc) (interface{}, error) {
	res, has := l.Get(key)
	if has {
		return res, nil
	}
	// key not exists, load and set cache
	return l.load(key, f)
}

func (l *localCache) Set(key string, value interface{}) {
	l.SetWithExpire(key, value, l.ttl)
}

func (l *localCache) SetWithExpire(key string, value interface{}, ttl int64) {
	obj, has := l.dict.Get(key)
	expireTime := time.Now().Add(time.Duration(ttl) * time.Second).Unix()
	if has {
		// update element info
		element := l.policy.unpack(obj)
		element.lock.Lock()
		element.value = value
		element.expireTime = expireTime
		// set ttl surround by lock
		l.ttlDict.Set(key, expireTime)
		element.lock.Unlock()
		// add hit count, if chan full, skip this signal is ok
		select {
		case l.hitChan <- obj:
		default:
		}
	} else {
		element := &element{
			key:        key,
			value:      value,
			expireTime: expireTime,
		}
		// add async by chan
		obj = l.policy.pack(element)
		l.opChan <- opMsg{opType: opTypeAdd, obj: obj}
	}
}

// Del delete key
func (l *localCache) Del(key string) {
	// del async by chan
	l.opChan <- opMsg{opType: opTypeDel, obj: key}
}

// Len return count of keys in cache
func (l *localCache) Len() int {
	return l.dict.Len()
}

// Flush clear all keys in cache
func (l *localCache) Flush() {
	l.dict.Flush()
	l.Stop()

	l.hitChan = make(chan interface{}, hitChanLen)
	l.opChan = make(chan opMsg, addChanLen)
	l.policy.flush()

	l.start()
}

// Stop the cacheProcess by close stopChan
func (l *localCache) Stop() {
	close(l.stopChan)
}

func (l *localCache) Statistic() map[string]interface{} {
	return map[string]interface{}{
		"hit":     l.statist.GetHitCount(),
		"miss":    l.statist.GetMissCount(),
		"hitRate": l.statist.GetHitRate(),
	}
}

// start cacheProcess
func (l *localCache) start() {
	l.stopChan = make(chan struct{})
	// deal chan signals
	go l.cacheProcess()
	//  delete the keys which are expired
	go l.ttlProcess()
}

// cacheProcess run a loop to deal chan signals
//  use a single goroutine to make policy ops safe.
func (l *localCache) cacheProcess() {
	for {
		select {
		case obj := <-l.hitChan:
			l.policy.hit(obj)
		case opMsg := <-l.opChan:
			if opMsg.opType == opTypeAdd {
				l.set(opMsg.obj)
			} else if opMsg.opType == opTypeDel {
				l.del(opMsg.obj.(string))
			}
		case <-l.stopChan:
			return
		}
	}
}

// load use singleFlight to load and set cache
func (l *localCache) load(key string, f LoadFunc) (interface{}, error) {
	loadF := func() (interface{}, error) {
		res, err := f()
		// if no err, set k-v to cache
		if err == nil {
			l.Set(key, res)
		}
		return res, err
	}
	// use singleFlight to load and set cache
	return l.group.Do(key, loadF)
}

// set called by single goroutine cacheProcess() to sync call
func (l *localCache) set(obj interface{}) {
	ele := l.policy.unpack(obj)
	objOld, has := l.dict.Get(ele.key)
	if has { // exists, del objOld from lru list
		l.policy.del(objOld)
	}
	l.dict.Set(ele.key, obj)
	// set ttl
	l.ttlDict.Set(ele.key, ele.expireTime)
	// add policy
	l.policy.add(obj)
}

// del called by single goroutine cacheProcess() to sync call
func (l *localCache) del(key string) {
	obj, has := l.dict.Get(key)
	if !has {
		return
	}
	// need del
	l.dict.Del(key)
	// del ttl
	l.ttlDict.Del(key)
	// del policy list
	l.policy.del(obj)
}

// ttlProcess run a loop to delete the keys which are expired
func (l *localCache) ttlProcess() {
	t := time.NewTicker(defaultTTLTick * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-l.stopChan:
			return
		case <-t.C:
			ti := time.Now()
			var delCount = 100
			// every 100ms, check rand 100 keys;
			// if expired more than 25, check again; like redis.
			// max run 50 ms.
			for delCount > defaultTTLCheckPercent &&
				time.Now().Sub(ti) < defaultTTLCheckRunTime*time.Millisecond {
				delCount = 0
				now := time.Now().Unix()
				keys := l.dict.RandKeys(defaultTTLCheckCount)
				distinctMap := make(map[string]struct{}, defaultTTLCheckCount)
				for _, key := range keys {
					if _, see := distinctMap[key]; see {
						continue
					}
					// add distinct key in map because RandKeys may repeat
					distinctMap[key] = struct{}{}
					v, has := l.ttlDict.Get(key)
					if has {
						// key expired, del it from dict & ttl dict
						expireTime := v.(int64)
						if now > expireTime {
							l.Del(key)
							delCount++
						}
					}
				}
			}
			//fmt.Println(time.Now(), time.Now().Sub(ti), l.ttlDict.Len(), l.Len())
		}
	}
}

// element is what factly save in dict
type element struct {
	lock       sync.RWMutex // element should be multi-safe
	key        string       // need key to del in policy when list is full
	value      interface{}
	expireTime int64
}

// isExpire return whether key is dead
func (e *element) isExpire() bool {
	return time.Now().Unix() > e.expireTime
}

// opMsg is a msg send to opChan when add or del a key
type opMsg struct {
	opType uint8       // type: add || del
	obj    interface{} // policy`s obj when set || key string when del
}
