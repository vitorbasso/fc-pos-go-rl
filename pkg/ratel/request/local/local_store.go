package local

import (
	"ratel/pkg/ratel/key"
	"ratel/pkg/ratel/request"
)

type LocalStore struct {
	items map[key.RatelKey]request.RequestItem
}

func NewLocalStore() *LocalStore {
	return &LocalStore{
		items: make(map[key.RatelKey]request.RequestItem),
	}
}

func (l *LocalStore) Set(k key.RatelKey, item request.RequestItem) error {
	l.items[k] = item
	return nil
}

func (l *LocalStore) Get(k key.RatelKey) (request.RequestItem, error) {
	item, ok := l.items[k]
	if !ok {
		return request.RequestItem{}, request.ErrNotFound
	}
	return item, nil
}
