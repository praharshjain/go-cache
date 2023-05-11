//go:build integration
// +build integration

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestInApp(t *testing.T) {
	initConfig()
	ctx := context.Background()
	key := "test-key"
	inApp := NewInApp()
	res, err := Cache(ctx, key, inApp, fn, "AbcD")
	//need to type cast back to the return type of the function (from interface{})
	s := res.(string)
	fmt.Printf("s:%s, err:%v", s, err)
	Delete(ctx, key, inApp)
}

func TestInRedis(t *testing.T) {
	initConfig()
	ctx := context.Background()
	key := "test-key"
	inRedis := NewInRedis("127.0.0.1", "6379")
	res, err := Cache(ctx, key, inRedis, fn2, "AbcD")
	//need to type cast back to the return type of the function (from interface{})
	s := res.(string)
	fmt.Printf("s:%s, err:%v", s, err)
	Delete(ctx, key, inRedis)
}

func fn(s string) (string, error) {
	return strings.ToLower(s), nil
}

func fn2(s string) (string, error) {
	return strings.ToUpper(s), nil
}

func initConfig() {
	cfg := make(map[string]Config)
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("config not found")
	}
	json.Unmarshal(b, &cfg)
	Init(cfg)
}
