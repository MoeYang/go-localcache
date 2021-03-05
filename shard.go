package localcache

import "sync"

type shard interface {
	get(key string) (interface{}, bool)
	set(key string, value interface{})
	del(key string) bool
	flush()
}

type shardMap struct {
	lock  sync.RWMutex
	store map[string]interface{}
}

func newShardMap() shard {
	return &shardMap{
		store: make(map[string]interface{}),
	}
}

func (s shardMap) get(key string) (interface{}, bool) {
	panic("implement me")
}

func (s shardMap) set(key string, value interface{}) {
	panic("implement me")
}

func (s shardMap) del(key string) bool {
	panic("implement me")
}

func (s shardMap) flush() {
	panic("implement me")
}
