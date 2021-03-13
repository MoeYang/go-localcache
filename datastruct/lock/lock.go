package lock

import (
	"sync"

	"github.com/MoeYang/go-localcache/common"
)

type Locker struct {
	locks       []*sync.RWMutex
	lockerCount uint32
}

func NewLocker(lockerCount uint32) *Locker {
	l := &Locker{
		locks:       make([]*sync.RWMutex, lockerCount),
		lockerCount: lockerCount,
	}
	for i := range l.locks {
		l.locks[i] = &sync.RWMutex{}
	}
	return l
}

func (l *Locker) Lock(key string) {
	lock := l.getLock(common.GetShardIndex(key, l.lockerCount))
	lock.Lock()
}

func (l *Locker) Unlock(key string) {
	lock := l.getLock(common.GetShardIndex(key, l.lockerCount))
	lock.Unlock()
}

func (l *Locker) RLock(key string) {
	lock := l.getLock(common.GetShardIndex(key, l.lockerCount))
	lock.RLock()
}

func (l *Locker) RUnlock(key string) {
	lock := l.getLock(common.GetShardIndex(key, l.lockerCount))
	lock.RUnlock()
}

func (l *Locker) getLock(idx uint32) *sync.RWMutex {
	return l.locks[idx]
}
