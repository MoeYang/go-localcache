package localcache

import "testing"

func TestAdd(t *testing.T) {
	c := NewLocalCache(WithCapacity(2), WithPolicy(PolicyTypeLRU))
	defer c.Stop()
	policy := c.(*localCache).policy.(*policyLRU)
	c.Set("1", 1)
	c.Set("2", 2)
	c.Set("3", 3)
	if policy.list.Len() != 2 {
		t.Errorf("TestAdd list len <> 2, len=%d", policy.list.Len())
	}
	if c.Len() != 2 {
		t.Errorf("TestAdd list len <> 2, len=%d", c.Len())
	}
	_, has := c.Get("1")
	if has {
		t.Error("TestAdd cache has 1")
	}
}

func TestHit(t *testing.T) {
	c := NewLocalCache(WithCapacity(2), WithPolicy(PolicyTypeLRU))
	defer c.Stop()
	policy := c.(*localCache).policy.(*policyLRU)
	c.Set("2", 2)
	c.Set("1", 1) // root -> 1 -> 2
	policy.hit(policy.list.Back())
	if policy.list.Front().Value.(*element).key != "2" {
		t.Errorf("TestHit list front <> 2, %+v", policy.list.Front().Value)
	}
}
