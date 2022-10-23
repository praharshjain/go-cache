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
}

func TestInRedis(t *testing.T) {
	ctx := context.Background()
	key := "test-key"
	inApp := NewInRedis()
	res, err := Cache(ctx, key, inApp, fn, "AbcD")
	fmt.Printf("res:%v, err:%v", res, err)
}

func fn(s string) (string, error) {
	return strings.ToLower(s), nil
}
