package rredis

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"ratel/pkg/ratel/key"
	"ratel/pkg/ratel/request"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	redisClient *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		redisClient: client,
	}
}

func (r *RedisStore) Set(k key.RatelKey, item request.RequestItem) error {
	ctx := context.Background()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(item)
	if err != nil {
		return err
	}
	if err = r.redisClient.Set(ctx, k.String(), buf.Bytes(), 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) Get(k key.RatelKey) (request.RequestItem, error) {
	ctx := context.Background()
	var item request.RequestItem
	by, err := r.redisClient.Get(ctx, k.String()).Bytes()
	if errors.Is(err, redis.Nil) {
		return request.RequestItem{}, request.ErrNotFound
	}
	if err != nil {
		return request.RequestItem{}, err
	}
	dec := gob.NewDecoder(bytes.NewReader(by))
	err = dec.Decode(&item)
	if err != nil {
		return request.RequestItem{}, err
	}
	return item, nil
}
