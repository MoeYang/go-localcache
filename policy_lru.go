package localcache

import (
	"github.com/MoeYang/go-localcache/datastruct/list"
)

type policyLRU struct {
	cap   int
	cache *localCache // the cache obj
	list  *list.List
}

func newPolicyLRU(cap int, cache *localCache) policy {
	return &policyLRU{
		cap:   cap,
		cache: cache,
		list:  list.New(),
	}
}

func (p *policyLRU) add(obj interface{}) {
	ele, ok := obj.(*list.Element)
	if !ok {
		return
	}
	// need to del when list is full
	if p.list.Len() >= p.cap {
		lastEle := p.list.Back()
		if lastEle != nil {
			// del from cache
			p.cache.del(lastEle.Value.(*element).key)
		}
	}
	// push ele to first of list
	p.list.PushElementFront(ele)
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

func (p *policyLRU) pack(ele *element) interface{} {
	return p.list.NewElement(ele)
}
