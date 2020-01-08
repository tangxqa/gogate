package syncmap

import (
	"sync"

	"github.com/tangxqa/gogate/throttle"
)

// 封装sync.map, 提供类型安全的方法调用
type RateLimiterSyncMap struct {
	rlMap			*sync.Map
}

func NewRateLimiterSyncMap() *RateLimiterSyncMap {
	return &RateLimiterSyncMap{
		rlMap: new(sync.Map),
	}
}

func (rsm *RateLimiterSyncMap) Get(key string) (throttle.RateLimiter, bool) {
	val, exist := rsm.rlMap.Load(key)
	if !exist {
		return nil, false
	}

	rl, ok := val.(throttle.RateLimiter)
	if !ok {
		return nil, false
	}

	return rl, true
}

func (rsm *RateLimiterSyncMap) Put(key string, val throttle.RateLimiter) {
	rsm.rlMap.Store(key, val)
}
