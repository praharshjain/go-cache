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

const (
	host = "127.0.0.1"
	port = "6379"
)

var (
	redisClient *redis.Client
)

func init() {
	redisOpts := &redis.Options{Addr: host + ":" + port}
	redisClient = redis.NewClient(redisOpts)
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("can't connect to redis err:%v", err))
	}
}

type InRedis struct {
}

func (i InRedis) GetName() string {
	return "InRedis"
}

func (i InRedis) Read(ctx context.Context, key string, opType reflect.Type) (interface{}, error) {
	val, err := redisClient.Get(ctx, key).Bytes()
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
	return redisClient.Set(ctx, key, b, expiration).Err()
}

func (i InRedis) Delete(ctx context.Context, key string) error {
	return redisClient.Del(ctx, key).Err()
}

func NewInRedis() Strategy {
	return InRedis{}
}
