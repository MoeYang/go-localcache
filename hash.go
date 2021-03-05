package localcache

const (
	// offset64 fnv-1a
	offset64 = 14695981039346656037
	// prime64 fnv-1a
	prime64 = 1099511628211
)

// Sum64 gets the string and returns its uint64 hash value. fnv-1a
func Sum64(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}
	return hash
}
