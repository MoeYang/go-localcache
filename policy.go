package localcache

const (
	PolicyTypeLRU = "lru"
	//PolicyTypeLFU = "lfu"
)

// policy of del useless element
type policy interface {
	// add an element
	add(interface{})
	// hit a key of element
	hit(interface{})
	// del a key
	del(interface{})
	// flush all key
	flush()
	// unpack interface to element
	unpack(interface{}) *element
	// pack element to interface
	pack(*element) interface{}
}

// newPolicy return policy implement by type const
func newPolicy(policyType string, cap int, cache *localCache) policy {
	var p policy
	switch policyType {
	case PolicyTypeLRU:
		p = newPolicyLRU(cap, cache)
	default:
		p = newPolicyLRU(cap, cache)
	}
	return p
}
