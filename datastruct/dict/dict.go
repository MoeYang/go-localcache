package dict

import (
	"math/rand"
	"sync"

	"github.com/MoeYang/go-localcache/common"
)

type Dict interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Del(key string) bool
	// RandKeys get count rand keys, may return keys repeat!
	RandKeys(count int) []string
	Len() int
	Flush()
}

// ForEachCustom used for Dict.ForEach() to deal k-v in every shard
//  while need to break, this fuch should return true
type ForEachCustom func(k string, v interface{}) (needBreak bool)

type concurrentMap struct {
	shards     []*shard
	shardCount uint32
}

// NewDict return a dict
func NewDict(shardCnt int) Dict {
	m := &concurrentMap{
		shardCount: uint32(shardCnt),
		shards:     make([]*shard, shardCnt),
	}
	for i := 0; i < shardCnt; i++ {
		m.shards[i] = newShard()
	}
	return m
}

func (m *concurrentMap) Get(key string) (interface{}, bool) {
	idx := common.GetShardIndex(key, m.shardCount)
	shard := m.getShard(idx)
	return shard.get(key)
}

func (m *concurrentMap) Set(key string, value interface{}) {
	idx := common.GetShardIndex(key, m.shardCount)
	shard := m.getShard(idx)
	shard.set(key, value)
}

func (m *concurrentMap) Del(key string) bool {
	idx := common.GetShardIndex(key, m.shardCount)
	shard := m.getShard(idx)
	return shard.del(key)
}

func (m *concurrentMap) Flush() {
	for _, shard := range m.shards {
		shard.flush()
	}
}

func (m *concurrentMap) Len() int {
	var l int
	for _, shard := range m.shards {
		l += shard.len()
	}
	return l
}

// RandKeys may return keys repeat!
func (m *concurrentMap) RandKeys(count int) []string {
	if maxCount := m.Len(); maxCount < count {
		count = maxCount
	}
	keys := make([]string, count)
	for i := 0; i < count; {
		shard := m.shards[rand.Intn(int(m.shardCount))]
		// randKey maybe "" if shards has not enough key, so get a key until not ""
		keys[i] = shard.randKey()
		if keys[i] != "" {
			i++
		}
	}
	return keys
}

// getShard get shard by shardIdx
func (m *concurrentMap) getShard(idx uint32) *shard {
	return m.shards[idx]
}

// shard is a concurrent safe map
type shard struct {
	lock  sync.RWMutex
	store map[string]interface{}
}

func newShard() *shard {
	return &shard{store: make(map[string]interface{})}
}

func (m *shard) get(key string) (interface{}, bool) {
	m.lock.RLock()
	v, has := m.store[key]
	m.lock.RUnlock()
	return v, has
}

func (m *shard) set(key string, value interface{}) {
	m.lock.Lock()
	m.store[key] = value
	m.lock.Unlock()
}

func (m *shard) del(key string) bool {
	m.lock.Lock()
	_, has := m.store[key]
	delete(m.store, key)
	m.lock.Unlock()
	return has
}

func (m *shard) flush() {
	m.lock.Lock()
	m.store = make(map[string]interface{})
	m.lock.Unlock()
}

func (m *shard) len() int {
	return len(m.store)
}

// foreach return true if custom func said need to break the loop
func (m *shard) foreach(custom ForEachCustom) bool {
	// need RLock the shard，to protect the write to shard
	m.lock.RLock()
	defer m.lock.RUnlock()
	for k, v := range m.store {
		needBreak := custom(k, v)

		if needBreak {
			return true
		}
	}
	return false
}

func (m *shard) randKey() string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for key := range m.store {
		return key
	}
	// shard is empty
	return ""
}
