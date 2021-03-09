# what is go-localcache ?
A mem cache library by golang.

# why choose go-localcache ?
1、Support to set TTL for every key 
2、LRU policy to delete useless keys
3、Similar performance like sync.Map[bench_test](https://github.com/MoeYang/go-localcache/tree/main/benchtest "bench_test")


# How to use
```go
	cache := localcache.NewLocalCache(
		localcache.WithCapacity(1024),
		localcache.WithShardCount(256),
		localcache.WithPolicy(localcache.PolicyTypeLRU),
		localcache.WithGlobalTTL(120),
	)
	
	// Get a key and return the value and if the key exists
	cache.Get(key string) (interface{}, bool)
	
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
```
