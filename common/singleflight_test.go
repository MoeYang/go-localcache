package common

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGroupDo(t *testing.T) {
	var g Group
	v, err := g.Do("k", func() (interface{}, error) {
		return "v", nil
	})
	if s, ok := v.(string); !ok || s != "v" || err != nil {
		t.Errorf("Err_TestGroup_Do1 return %+v", v)
	}
	v, err = g.Do("k", func() (interface{}, error) {
		return nil, errors.New("err")
	})
	if err == nil {
		t.Error("Err_TestGroup_Do2 err=nil")
	}
}

func TestGroupDoMulti(t *testing.T) {
	var g Group
	var wg sync.WaitGroup
	ch := make(chan struct{})
	var callCnt int32
	fn := func() (interface{}, error) {
		<-ch
		atomic.AddInt32(&callCnt, 1)
		return "v", nil
	}
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := g.Do("k", fn)
			if err != nil {
				t.Error("TestGroupDoMulti err != nil")
			}
		}()
	}
	// sleep a while to wait all goroutines ready
	time.AfterFunc(100*time.Millisecond, func() {
		close(ch)
	})
	wg.Wait()
	if atomic.LoadInt32(&callCnt) != 1 {
		t.Error("TestGroupDoMulti callcnt != 1")
	}
}
