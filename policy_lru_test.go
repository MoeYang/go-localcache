package localcache

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	c := NewLocalCache(WithCapacity(2), WithPolicy(PolicyTypeLRU))
	defer c.Stop()
	policy := c.(*localCache).policy.(*policyLRU)
	c.Set("1", 1)
	c.Set("2", 2)
	c.Set("3", 3)
	time.Sleep(10 * time.Millisecond)
	if policy.list.Len() != 2 {
		t.Errorf("TestAdd list len <> 2, len=%d", policy.list.Len())
	}
	if c.Len() != 2 {
		t.Errorf("TestAdd list len <> 2, len=%d", c.Len())
	}
	if policy.list.Front().Value.(*element).key != "3" {
		t.Errorf("TestAdd list front <> 3, %+v", policy.list.Front().Value)
	}
	_, has := c.Get("1")
	if has {
		t.Error("TestAdd cache has 1")
	}
	_, has = c.Get("2")
	if !has {
		t.Error("TestAdd cache has 2")
	}
	_, has = c.Get("3")
	if !has {
		t.Error("TestAdd cache has 3")
	}
}

func TestHit(t *testing.T) {
	c := NewLocalCache(WithCapacity(2), WithPolicy(PolicyTypeLRU))
	defer c.Stop()
	policy := c.(*localCache).policy.(*policyLRU)
	c.Set("2", 2)
	c.Set("1", 1) // root -> 1 -> 2
	time.Sleep(10 * time.Millisecond)
	c.Get("2")
	time.Sleep(10 * time.Millisecond)
	//policy.hit(policy.list.Back())
	if policy.list.Front().Value.(*element).key != "2" {
		t.Errorf("TestHit list front <> 2, %+v", policy.list.Front().Value)
	}
}
