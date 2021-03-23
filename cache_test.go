package localcache

import (
	"sync"
	"sync/atomic"
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
	time.Sleep(1 * time.Millisecond)
	_, has = c.Get("123")
	if !has {
		t.Error("TestGet2 not exists")
	}
	c.Set("123", 3)
	time.Sleep(1 * time.Millisecond)
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

func TestGetOrLoad(t *testing.T) {
	c := NewLocalCache()
	defer c.Stop()
	var wg sync.WaitGroup
	var callCnt int32
	ch := make(chan struct{})
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := c.GetOrLoad("k", func() (interface{}, error) {
				<-ch
				return atomic.AddInt32(&callCnt, 1), nil
			})
			if err != nil || res.(int32) != 1 {
				t.Errorf("TestGetOrLoad err != nil: res=%d", res)
			}
		}()
	}
	time.AfterFunc(100*time.Millisecond, func() {
		close(ch)
	})
	wg.Wait()
	if atomic.LoadInt32(&callCnt) != 1 {
		t.Error("TestGetOrLoad callcnt != 1")
	}
}
func TestDel(t *testing.T) {
	c := NewLocalCache()
	defer c.Stop()
	c.Set("123", 1)
	c.Del("123")
	time.Sleep(1 * time.Millisecond)
	_, has := c.Get("123")
	if has {
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
	time.Sleep(1 * time.Millisecond)
	l = c.Len()
	if l != 1 {
		t.Error("TestLen2 <> 1")
	}
}

func TestFlush(t *testing.T) {
	c := NewLocalCache()
	defer c.Stop()
	c.Set("123", 1)
	time.Sleep(1 * time.Millisecond)
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
