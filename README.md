A mem cache library by golang.

# Why choose go-localcache ?
1、Support to set different TTL for every key 

2、LRU policy to delete useless keys

3、Similar performance like sync.Map -> see [bench_test](https://github.com/MoeYang/go-localcache/tree/main/benchtest "bench_test")


# How to use
```go
	cache := localcache.NewLocalCache(
		localcache.WithCapacity(1024), // WithShardCount set max Capacity
		localcache.WithShardCount(256),// WithShardCount shardCnt must be a power of 2
		localcache.WithGlobalTTL(120), // WithGlobalTTL set all keys default expire time of seconds
		localcache.WithStatist(true),  // WithStatist set whether need to caculate the cache stastic
		localcache.WithPolicy(localcache.PolicyTypeLRU), // WithPolicy set the elimination policy of key
	)
	
	// Get a key and return the value and if the key exists
	cache.Get(key string) (interface{}, bool)

	// GetOrLoad get a key, while key not exists, call f() to load data, and will set the load data to cache.
	// Load data process will called singleFlight called 
	cache.GetOrLoad(key string, f LoadFunc) (interface{}, error)

	// Set a key-value with default seconds to live
	cache.Set(key string, value interface{})
	
	// SetWithExpire set a key-value with seconds to live
	cache.SetWithExpire(key string, value interface{}, ttl int64)
	
	// Del delete key and return if the key exists
	cache.Del(key string) bool
	
	// Len return count of keys in cache
	cache.Len() int
	
	// Flush clear all keys in chache, should do this when set and del is stop
	cache.Flush()
	
	// Stop the cacheProcess by close stopChan
	cache.Stop()

	// Statistic return cache Statics {"hit":1, "miss":1, "hitRate":50.0}
	Statistic() map[string]interface{}
```
