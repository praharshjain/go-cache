package cache

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"
)

var (
	sm = sync.Map{}
)

type InApp struct {
}

// Item is a struct denoting the format in which any object is stored in In-APP Cache. It contains the object along with an expiration time.
type Item struct {
	object     interface{}
	expiration time.Time
}

func NewInApp() Strategy {
	return InApp{}
}

func (i InApp) GetName() string {
	return "InApp"
}

func (i InApp) Read(ctx context.Context, key string, opType reflect.Type) (interface{}, error) {
	currTime := time.Now()
	result, ok := sm.Load(key)
	if !ok {
		return nil, errors.New("cache miss")
	}
	item, isValid := result.(Item)
	if !isValid {
		return nil, errors.New("cached entry is not of expected type")
	}
	expTime := item.expiration
	if currTime.After(expTime) {
		sm.Delete(key)
		return nil, errors.New("cache expired")
	}
	return item.object, nil
}

func (i InApp) Write(ctx context.Context, key string, expiration time.Duration, res interface{}, err error) error {
	sm.Store(key, Item{object: res, expiration: time.Now().Add(expiration)})
	return nil
}

func (i InApp) Delete(ctx context.Context, key string) error {
	sm.Delete(key)
	return nil
}
