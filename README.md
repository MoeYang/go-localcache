# go-localcache
A mem cache library by golang.

Support to set TTL for every key and LRU policy to delete useless keys. 


# How to use
```go
	cache := localcache.NewLocalCache(
		localcache.WithCapacity(1024),
		localcache.WithShardCount(256),
		localcache.WithPolicy(localcache.PolicyTypeLRU),
		localcache.WithGlobalTTL(120),
	)
	
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
```
