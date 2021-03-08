package localcache

const (
	defaultShardCnt = 16 // ShardCnt must be a power of 2
	defaultTTL      = 60 // default key expire time 60 sec
	defaultCap      = 1024
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Del(key string) bool
	Len() int
	Flush() // clear all keys in cache
}

type localCache struct {
	shards    []shard
	shardCnt  int // shardings count
	shardMask uint64
	cap       int // capacity
	ttl       int // Global Keys expire seconds

}

// NewLocalCache return Cache obj with options
func NewLocalCache(options ...Option) Cache {
	c := &localCache{
		shardCnt:  defaultShardCnt,
		shardMask: defaultShardCnt - 1,
		cap:       defaultCap,
		ttl:       defaultTTL,
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
	return c
}

type Option func(*localCache)

// WithGlobalTTL set all keys default expire time of seconds
func WithGlobalTTL(expireSecond int) Option {
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

func (l localCache) Get(key string) (interface{}, bool) {
	idx := l.getShardIndex(Sum64(key))
	return l.shards[idx].get(key)
}

func (l localCache) Set(key string, value interface{}) {
	idx := l.getShardIndex(Sum64(key))
	l.shards[idx].set(key, value)
}

// Del delete key and return if the key exists
func (l localCache) Del(key string) bool {
	idx := l.getShardIndex(Sum64(key))
	return l.shards[idx].del(key)
}

// Flush clear all keys in chache
func (l localCache) Flush() {
	for _, shard := range l.shards {
		shard.flush()
	}
}

// Len return count of keys in cache
func (l localCache) Len() int {
	var cnt int
	for _, shard := range l.shards {
		cnt += shard.len()
	}
	return cnt
}

// getShardIndex getShardIndex by hash code of key
func (l localCache) getShardIndex(n uint64) uint64 {
	return n & l.shardMask
}
