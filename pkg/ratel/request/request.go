package request

import (
	"errors"
	"ratel/pkg/ratel/key"
)

var (
	ErrNotFound = errors.New("key not found")
)

type RequestItem struct {
	Requests  map[int64]int
	BlockedAt int64
}

type RequestStore interface {
	Set(k key.RatelKey, item RequestItem) error
	Get(k key.RatelKey) (RequestItem, error)
}
