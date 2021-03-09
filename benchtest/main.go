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
	readTimes = 100
)

func main() {
	test1()
	test2()
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

func test1() {
	cache := localcache.NewLocalCache(localcache.WithCapacity(keys))
	{
		startT := time.Now() //计算当前时间
		for i := 0; i < keys; i++ {
			cache.Set(fmt.Sprintf("%s%010d", kPre, i), []byte(vPre+strconv.Itoa(i)))
		}
		fmt.Printf("write local cost = %v\n", time.Since(startT))
	}
	//读性能测试
	startT := time.Now() //计算当前时间
	var wg sync.WaitGroup
	wg.Add(readTimes)
	for i := 0; i < readTimes; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < keys; j++ {
				v, err := cache.Get(fmt.Sprintf("%s%010d", kPre, j))
				if v == nil || err {
					fmt.Sprintf("%d", 1)
				}
			}

		}()
	}
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
}
func test2() {
	var cache sync.Map
	{
		startT := time.Now() //计算当前时间
		for i := 0; i < keys; i++ {
			cache.Store(fmt.Sprintf("%s%010d", kPre, i), []byte(vPre+strconv.Itoa(i)))
		}
		fmt.Printf("write sync.Map cost = %v\n", time.Since(startT))
	}
	//读性能测试
	startT := time.Now() //计算当前时间
	var wg sync.WaitGroup
	wg.Add(readTimes)
	for i := 0; i < readTimes; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < keys; j++ {
				v, err := cache.Load(fmt.Sprintf("keyprefix%010d", i, j))
				if v == nil || !err {
					fmt.Sprintf("%d", 1)
				}
			}

		}()
	}
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
