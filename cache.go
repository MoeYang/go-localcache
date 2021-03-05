package localcache

const (
	defaultShardCnt = 16
	defaultTTL      = 60 // default key expire time 60 sec
	defaultCap      = 1024
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Del(key string) bool
	Flush()
}

type localCache struct {
	shards   []shard
	shardCnt int // shardings count
	size     int // key count
	cap      int // capacity
	ttl      int // Global Keys expire seconds

}

func NewLocalCache(options ...Option) Cache {
	c := &localCache{}
	for _, opt := range options {
		opt(c)
	}
	c.shards = make([]shard, c.shardCnt)
	for i := 0; i < c.shardCnt; i++ {

	}
	return c
}

type Option func(*localCache)

func WithGlobalTTL(expireSecond int) Option {
	if expireSecond <= 0 {
		expireSecond = defaultTTL
	}
	return func(c *localCache) {
		c.ttl = expireSecond
	}
}

func WithCapacity(cap int) Option {
	if cap <= 0 {
		cap = defaultCap
	}
	return func(c *localCache) {
		c.cap = cap
	}
}

func WithShardCount(shardCnt int) Option {
	if shardCnt <= 0 {
		shardCnt = defaultShardCnt
	}
	return func(c *localCache) {
		c.shardCnt = shardCnt
	}
}

func (l localCache) Get(key string) (interface{}, bool) {
	panic("implement me")
}

func (l localCache) Set(key string, value interface{}) {
	panic("implement me")
}

func (l localCache) Del(key string) bool {
	panic("implement me")
}

func (l localCache) Flush() {
	panic("implement me")
}
