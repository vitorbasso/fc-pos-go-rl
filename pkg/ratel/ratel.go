package ratel

import (
	"errors"
	"ratel/pkg/ratel/key"
	"ratel/pkg/ratel/rconfig"
	"ratel/pkg/ratel/request"
	"ratel/pkg/ratel/request/local"
	"sync"
	"time"
)

type Ratel struct {
	configStore  *rconfig.ConfigStore
	requestStore request.RequestStore
	mu           sync.Mutex
}

func NewLimiter(options ...Option) *Ratel {
	ratel := &Ratel{
		configStore:  rconfig.NewConfigStore(),
		requestStore: local.NewLocalStore(),
	}
	for _, option := range options {
		option(ratel)
	}
	return ratel
}

func (r *Ratel) Allow(k string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cfg, rk := r.GetConfig(k)
	if cfg.TimeWindowSecond <= 0 {
		return true, nil
	}

	requestItem, err := r.requestStore.Get(rk)
	if errors.Is(err, request.ErrNotFound) {
		requestItem = request.RequestItem{
			Requests: make(map[int64]int),
		}
	} else if err != nil {
		return false, err
	}
	if requestItem.BlockedAt+cfg.BlockDurationSecond >= time.Now().Unix() {
		return false, nil
	}
	requestItem.Requests[time.Now().Unix()]++
	totalRequests := 0
	for req, count := range requestItem.Requests {
		if req >= (time.Now().Unix() - cfg.TimeWindowSecond) {
			totalRequests += count
		} else {
			delete(requestItem.Requests, req)
		}
	}
	blocked := totalRequests > cfg.Capacity
	if blocked {
		requestItem.BlockedAt = time.Now().Unix()
	}
	err = r.requestStore.Set(rk, requestItem)
	if err != nil {
		return false, err
	}
	return !blocked, nil
}

func (r *Ratel) GetConfig(k string) (rconfig.ConfigItem, key.RatelKey) {
	rk := key.TokenKey(k)
	cfg, ok := r.configStore.GetConfig(rk)
	if !ok {
		rk = key.IPKey(k)
		cfg, ok = r.configStore.GetConfig(rk)
	}
	if !ok {
		cfg = r.configStore.GetDefaultConfig()
	}
	return cfg, rk
}
