package localcache

import "sync/atomic"

type statist interface {
	// hitIncr add hit count
	hitIncr()
	// missIncr add miss count
	missIncr()
	GetHitCount() uint64
	GetMissCount() uint64
	GetHitRate() float64
}

// statisCaculator implement a cache statist
type statisCaculator struct {
	needStastic bool
	hitCount    uint64
	missCount   uint64
}

// newstatisCaculator needs a param whether need to do cache statist
func newstatisCaculator(needStastic bool) statist {
	return &statisCaculator{needStastic: needStastic}
}

func (s *statisCaculator) hitIncr() {
	if !s.needStastic {
		return
	}
	atomic.AddUint64(&s.hitCount, 1)
}

func (s *statisCaculator) missIncr() {
	if !s.needStastic {
		return
	}
	atomic.AddUint64(&s.missCount, 1)
}

func (s *statisCaculator) GetHitCount() uint64 {
	return atomic.LoadUint64(&s.hitCount)
}

func (s *statisCaculator) GetMissCount() uint64 {
	return atomic.LoadUint64(&s.missCount)
}

func (s *statisCaculator) GetHitRate() float64 {
	hit, miss := s.GetHitCount(), s.GetMissCount()
	if total := hit + miss; total != 0 {
		return float64(hit) / float64(hit+miss) * 100
	}
	return 0
}
