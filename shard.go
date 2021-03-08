package localcache

import "sync"

type shard interface {
	get(key string) (interface{}, bool)
	set(key string, value interface{})
	del(key string) bool
	len() int
	flush()
}

// shardMap is implement of shard
type shardMap struct {
	lock  sync.RWMutex
	store map[string]interface{}
}

// newShardMap return a shard
func newShardMap() shard {
	return &shardMap{
		store: make(map[string]interface{}),
	}
}

func (sm *shardMap) get(key string) (interface{}, bool) {
	sm.lock.RLock()
	v, has := sm.store[key]
	sm.lock.RUnlock()
	return v, has
}

func (sm *shardMap) set(key string, value interface{}) {
	sm.lock.Lock()
	sm.store[key] = value
	sm.lock.Unlock()
}

func (sm *shardMap) del(key string) bool {
	sm.lock.Lock()
	_, has := sm.store[key]
	delete(sm.store, key)
	sm.lock.Unlock()
	return has
}

func (sm *shardMap) flush() {
	sm.lock.Lock()
	sm.store = make(map[string]interface{})
	sm.lock.Unlock()
}

func (sm *shardMap) len() int {
	return len(sm.store)
}
