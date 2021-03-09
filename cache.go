package localcache

import (
	"sync"
	"time"
)

const (
	defaultShardCnt = 256 // ShardCnt must be a power of 2
	defaultTTL      = 60  // default key expire time 60 sec
	defaultCap      = 1024

	hitChanLen = 1 << 15
)

type Cache interface {
	// Get a key and return the value and if the key exists
	Get(key string) (interface{}, bool)
	// Set a key-value with default seconds to live
	Set(key string, value interface{})
	// SetWithExpire set a key-value with seconds to live
	SetWithExpire(key string, value interface{}, ttl int64)
	// Del delete key and return if the key exists
	Del(key string) bool
	// Len return count of keys in cache
	Len() int
	// Flush clear all keys in chache, should do this when set and del is stop
	Flush()
	// Stop the cacheProcess by close stopChan
	Stop()
}

type localCache struct {
	policy     policy // elimination policy of keys
	policyType string
	shards     []shard
	shardCnt   int // shardings count
	shardMask  uint64
	cap        int   // capacity
	ttl        int64 // Global Keys expire seconds

	hitChan  chan interface{} // chan while get a key should put in
	stopChan chan struct{}    // chan stop signal
}

// NewLocalCache return Cache obj with options
func NewLocalCache(options ...Option) Cache {
	c := &localCache{
		shardCnt:  defaultShardCnt,
		shardMask: defaultShardCnt - 1,
		cap:       defaultCap,
		ttl:       defaultTTL,
		hitChan:   make(chan interface{}, hitChanLen),
	}
	// set options
	for _, opt := range options {
		opt(c)
	}
	// init shardings
	c.shards = make([]shard, c.shardCnt)
	for i := 0; i < c.shardCnt; i++ {
		c.shards[i] = newShardMap()
	}
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
		c.shardMask = uint64(shardCnt - 1)
	}
}

// WithPolicy set the elimination policy of keys
func WithPolicy(policyType string) Option {
	return func(c *localCache) {
		c.policyType = policyType
	}
}

func (l *localCache) Get(key string) (interface{}, bool) {
	idx := l.getShardIndex(sum64(key))
	obj, has := l.shards[idx].get(key)
	if has {
		element := l.policy.unpack(obj)
		element.lock.RLock()
		defer element.lock.RUnlock()
		if element.isExpire() {
			// add hit count, if chan full, skip this signal is ok
			select {
			case l.hitChan <- obj:
			default:
			}
			return element.value, true
		} else {
			// out of ttl, need del
			l.Del(key)
		}
	}
	return nil, false
}

func (l *localCache) Set(key string, value interface{}) {
	l.SetWithExpire(key, value, l.ttl)
}

func (l *localCache) SetWithExpire(key string, value interface{}, ttl int64) {
	idx := l.getShardIndex(sum64(key))
	obj, has := l.shards[idx].get(key)
	if has {
		// update element info
		element := l.policy.unpack(obj)
		element.lock.Lock()
		element.value = value
		element.expireTime = time.Now().Add(time.Duration(ttl) * time.Second).Unix()
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
			expireTime: time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
		}
		// add new key
		obj := l.policy.add(element)
		l.shards[idx].set(key, obj)
	}
}

// Del delete key and return if the key exists
func (l *localCache) Del(key string) bool {
	idx := l.getShardIndex(sum64(key))
	obj, has := l.shards[idx].get(key)
	if has {
		// need del
		l.policy.del(obj)
		l.shards[idx].del(key)
		return true
	}
	return false
}

// Len return count of keys in cache
func (l *localCache) Len() int {
	var cnt int
	for _, shard := range l.shards {
		cnt += shard.len()
	}
	return cnt
}

// Flush clear all keys in chache
func (l *localCache) Flush() {
	for _, shard := range l.shards {
		shard.flush()
	}
	l.Stop()

	l.hitChan = make(chan interface{}, hitChanLen)
	l.policy.flush()

	l.start()
}

// Stop the cacheProcess by close stopChan
func (l *localCache) Stop() {
	close(l.stopChan)
}

// start cacheProcess
func (l *localCache) start() {
	l.stopChan = make(chan struct{})
	go l.cacheProcess()
}

// cacheProcess run a loop to deal chan signals
func (l *localCache) cacheProcess() {
	for {
		select {
		case obj := <-l.hitChan:
			l.policy.hit(obj)
		case <-l.stopChan:
			return
		}
	}

}

// getShardIndex getShardIndex by hash code of key
func (l *localCache) getShardIndex(n uint64) uint64 {
	return n & l.shardMask
}

// element is what factly save in []shard
type element struct {
	lock       sync.RWMutex // element should be multi-safe
	key        string
	value      interface{}
	expireTime int64
}

// isExpire return whether element in ttl
func (e *element) isExpire() bool {
	return time.Now().Unix() <= e.expireTime
}
