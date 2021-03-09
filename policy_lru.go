package localcache

import (
	"container/list"
)

type policyLRU struct {
	cap   int
	cache Cache // the cache obj
	list  *list.List
}

func newPolicyLRU(cap int, cache Cache) policy {
	return &policyLRU{
		cap:   cap,
		cache: cache,
		list:  list.New(),
	}
}

func (p *policyLRU) add(ele *element) interface{} {
	obj := p.list.PushFront(ele)
	// need to del when list is full( ele is already in list, so this should use > )
	if p.list.Len() > p.cap {
		lastEle := p.list.Back()
		if lastEle != nil {
			p.cache.Del(lastEle.Value.(*element).key)
		}
	}
	return obj
}

func (p *policyLRU) hit(obj interface{}) {
	ele, ok := obj.(*list.Element)
	if !ok {
		return
	}
	p.list.MoveToFront(ele)
}

func (p *policyLRU) del(obj interface{}) {
	ele, ok := obj.(*list.Element)
	if !ok {
		return
	}
	p.list.Remove(ele)
}

func (p *policyLRU) flush() {
	p.list = list.New()
}

// unpack decode a *list.Element and return *element
func (p *policyLRU) unpack(obj interface{}) *element {
	ele, ok := obj.(*list.Element)
	if !ok {
		return nil
	}
	return ele.Value.(*element)
}
