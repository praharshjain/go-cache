package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v9"
)

type InRedis struct {
	client *redis.Client
}

func (i InRedis) GetName() string {
	return "InRedis"
}

func (i InRedis) Read(ctx context.Context, key string, opType reflect.Type) (interface{}, error) {
	val, err := i.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	if len(val) == 0 {
		return nil, errors.New("zero bytes found for redis key:" + key)
	}
	opVal := reflect.New(opType)
	res := opVal.Interface()
	err = json.Unmarshal(val, res)
	return opVal.Elem().Interface(), err
}

func (i InRedis) Write(ctx context.Context, key string, expiration time.Duration, res interface{}, err error) error {
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return i.client.Set(ctx, key, b, expiration).Err()
}

func (i InRedis) Delete(ctx context.Context, key string) error {
	return i.client.Del(ctx, key).Err()
}

func NewInRedis(host, port string) Strategy {
	s := InRedis{}
	redisOpts := &redis.Options{Addr: host + ":" + port}
	s.client = redis.NewClient(redisOpts)
	_, err := s.client.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("can't connect to redis err:%v", err))
	}
	return s
}
