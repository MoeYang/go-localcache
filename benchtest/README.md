# bench_test

1、write keys with 1 goroutine; 
2、get keys with 100 goroutines while write keys with 1 goroutine and del keys with 1 goroutine;


     go-localcache bench
	 write local cost = 2.132928994s
	 read\write\delete local cost = 11.856400333s
	 2021/03/09 20:09:50  localcache Pause:87826 1
	 delete local cost = 404.306562ms
	 
	 sync.Map bench
     	 write sync.Map cost = 1.263461937s
	 read\write\delete sync.Map cost = 9.6866624s
	 2021/03/09 20:10:08  sync.Map Pause:111403 1
	 del sync.Map cost = 325.917375ms


