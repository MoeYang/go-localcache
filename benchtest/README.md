# bench_test

1、write keys with 1 goroutine; 
2、get keys with 100 goroutines while write keys with 1 goroutine and del keys with 1 goroutine;


# localcache
```
write local cost = 1.425530094s
read local cost = 10.496183401s
2021/03/14 22:19:26  localcache Pause:93065 1
delete local cost = 534.445934ms
map[hit:100000000 hitRate:100 miss:0]
```

# sync.map
```
write sync.Map cost = 1.281531642s
read sync.Map cost = 9.296012626s
2021/03/14 22:19:37  sync.Map Pause:109340 1
del sync.Map cost = 324.273299ms
```
	 
	 
	 
	 


