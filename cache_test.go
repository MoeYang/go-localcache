package localcache

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	c := NewLocalCache()
	defer c.Stop()
	_, has := c.Get("123")
	if has {
		t.Error("TestGet1 not exists")
	}
	c.Set("123", 1)
	_, has = c.Get("123")
	if !has {
		t.Error("TestGet2 not exists")
	}
	c.Set("123", 3)
	obj, _ := c.Get("123")
	if obj.(int) != 3 {
		t.Errorf("TestGet3 get <> 3 %+v", obj)
	}
	c.SetWithExpire("123", 2, 0)
	time.Sleep(1 * time.Second)
	obj, has = c.Get("123")
	if has {
		t.Errorf("TestGet4 not exists %+v", obj)
	}
}

func TestDel(t *testing.T) {
	c := NewLocalCache()
	defer c.Stop()
	has := c.Del("123")
	if has {
		t.Error("TestDel1 not exists")
	}
	c.Set("123", 1)
	has = c.Del("123")
	if !has {
		t.Error("TestDel2 not exists")
	}
}

func TestLen(t *testing.T) {
	c := NewLocalCache()
	defer c.Stop()
	l := c.Len()
	if l != 0 {
		t.Error("TestLen1 <> 0")
	}
	c.Set("123", 1)
	l = c.Len()
	if l != 1 {
		t.Error("TestLen2 <> 1")
	}
}

func TestFlush(t *testing.T) {
	c := NewLocalCache()
	defer c.Stop()
	c.Set("123", 1)
	l := c.Len()
	if l != 1 {
		t.Error("TestFlush1 <> 1")
	}
	c.Flush()
	l = c.Len()
	if l != 0 {
		t.Error("TestFlush2 <> 0")
	}
}
