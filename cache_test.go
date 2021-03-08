package localcache

import (
	"testing"
)

func TestGet(t *testing.T) {
	c := NewLocalCache()
	_, has := c.Get("123")
	if has {
		t.Error("TestGet1 not exists")
	}
	c.Set("123", 1)
	_, has = c.Get("123")
	if !has {
		t.Error("TestGet2 not exists")
	}
}

func TestDel(t *testing.T) {
	c := NewLocalCache()
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
