package common

const prime32 = uint32(16777619)

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// GetShardIndex by hash code of key
func GetShardIndex(key string, shardCount uint32) uint32 {
	return fnv32(key) & (shardCount - 1)
}
