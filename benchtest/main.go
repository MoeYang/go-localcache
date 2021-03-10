package main

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/MoeYang/go-localcache"
)

const (
	vPre      = "v"
	kPre      = "1234567890k"
	keys      = 1000000
	readTimes = 50
)

func main() {
	testLocalCache()
	testSyncMap()
}

var ss debug.GCStats

func gcPause(caller string) {
	runtime.GC()
	debug.ReadGCStats(&ss)
}

func gcPauseX(caller string) {
	runtime.GC()
	var ss1 debug.GCStats
	debug.ReadGCStats(&ss1)
	log.Printf(" %s Pause:%d %d\n", caller, ss1.PauseTotal-ss.PauseTotal, ss1.NumGC-ss.NumGC)
}

func testLocalCache() {
	cache := localcache.NewLocalCache(
		localcache.WithCapacity(keys),
		localcache.WithShardCount(256),
		localcache.WithPolicy(localcache.PolicyTypeLRU),
		localcache.WithGlobalTTL(120),
		localcache.WithStatist(true),
	)
	{
		startT := time.Now() //计算当前时间
		for i := 0; i < keys; i++ {
			cache.Set(fmt.Sprintf("%s%d", kPre, i), []byte(vPre+strconv.Itoa(i)))
		}
		fmt.Printf("write local cost = %v\n", time.Since(startT))
	}
	//读性能测试
	startT := time.Now() //计算当前时间
	var wg sync.WaitGroup
	wg.Add(readTimes + 1)
	for i := 0; i < readTimes; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < keys; j++ {
				v, has := cache.Get(fmt.Sprintf("%s%d", kPre, j))
				if v == nil || !has {
					fmt.Println(j)
				}
			}

		}()
	}
	go func() {
		defer wg.Done()
		for i := 0; i < keys; i++ {
			cache.Set(fmt.Sprintf("%s%d", kPre, i), []byte(vPre+strconv.Itoa(i)))
		}
		for i := 0; i < keys; i++ {
			cache.Set(fmt.Sprintf("%s%d", kPre, i), []byte(vPre+strconv.Itoa(i)))
		}
	}()
	wg.Wait()
	fmt.Printf("read local cost = %v\n", time.Since(startT))
	gcPause("localcache")
	gcPauseX("localcache")
	{
		startT1 := time.Now() //计算当前时间
		for j := 0; j < keys; j++ {
			cache.Del(fmt.Sprintf("%s%d", kPre, j))
		}
		fmt.Printf("delete local cost = %v\n", time.Since(startT1))
	}
	fmt.Println(cache.Statics())
}
func testSyncMap() {
	var cache sync.Map
	{
		startT := time.Now() //计算当前时间
		for i := 0; i < keys; i++ {
			cache.Store(fmt.Sprintf("%s%d", kPre, i), []byte(vPre+strconv.Itoa(i)))
		}
		fmt.Printf("write sync.Map cost = %v\n", time.Since(startT))
	}
	//读性能测试
	startT := time.Now() //计算当前时间
	var wg sync.WaitGroup
	wg.Add(readTimes + 1)
	for i := 0; i < readTimes; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < keys; j++ {
				v, has := cache.Load(fmt.Sprintf("%s%d", kPre, j))
				if v == nil || !has {
					fmt.Println(j)
				}
			}

		}()
	}
	go func() {
		defer wg.Done()
		for i := 0; i < keys; i++ {
			cache.Store(fmt.Sprintf("%s%d", kPre, i), []byte(vPre+strconv.Itoa(i)))
		}
		for i := 0; i < keys; i++ {
			cache.Store(fmt.Sprintf("%s%d", kPre, i), []byte(vPre+strconv.Itoa(i)))
		}
	}()
	wg.Wait()
	fmt.Printf("read sync.Map cost = %v\n", time.Since(startT))

	gcPause("sync.Map")
	gcPauseX("sync.Map")
	//del性能测试
	{
		startT1 := time.Now() //计算当前时间
		for j := 0; j < keys; j++ {
			cache.Delete(fmt.Sprintf("%s%d", kPre, j))
		}
		fmt.Printf("del sync.Map cost = %v\n", time.Since(startT1))
	}
}
