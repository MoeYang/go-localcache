package common

import "sync"

// call is an in-flight or completed Do call
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
// like golang.org/x/sync/singleflight
type Group struct {
	lock sync.Mutex
	m    map[string]*call
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.lock.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, has := g.m[key]; has {
		g.lock.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.lock.Unlock()
	// call func
	c.val, c.err = fn()
	c.wg.Done()
	// delete key from map
	g.lock.Lock()
	delete(g.m, key)
	g.lock.Unlock()

	return c.val, c.err
}
