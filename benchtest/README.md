# bench_test

We set keys with 1 goroutine, and get keys with 100 goroutines.


     go-localcache bench
	 write local cost = 1.355967273s
	 read local cost = 15.260196832s
	 2021/03/09 20:09:50  localcache Pause:59311 1
	 delete local cost = 440.067947ms
	 
	 sync.Map bench
     write sync.Map cost = 1.341951564s
	 read sync.Map cost = 15.626010979s
	 2021/03/09 20:10:08  sync.Map Pause:59965 1
	 del sync.Map cost = 279.28925ms


