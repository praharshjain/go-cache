package cache

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestInApp(t *testing.T) {
	ctx := context.Background()
	key := "test-key"
	inApp := NewInApp()
	res, err := Cache(ctx, key, inApp, fn, "AbcD")
	fmt.Printf("res:%v, err:%v", res, err)
	Delete(ctx, key, inApp)
}

func TestInRedis(t *testing.T) {
	ctx := context.Background()
	key := "test-key"
	inRedis := NewInRedis()
	res, err := Cache(ctx, key, inRedis, fn, "AbcD")
	fmt.Printf("res:%v, err:%v", res, err)
	Delete(ctx, key, inRedis)
}

func fn(s string) (string, error) {
	return strings.ToLower(s), nil
}
