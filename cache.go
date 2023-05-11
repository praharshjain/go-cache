package cache

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var (
	cacheConfig map[string]Config
)

// Strategy is a custom data type used to denote the caching strategy (i.e. where to cache).
type Strategy interface {
	GetName() string
	Read(ctx context.Context, key string, opType reflect.Type) (interface{}, error)
	Write(ctx context.Context, key string, expiration time.Duration, res interface{}, err error) error
	Delete(ctx context.Context, key string) error
}

// Config contains parameters for all cache related configs.
type Config struct {
	TtlInSeconds int64 `json:"ttlInSeconds"`
	Enabled      bool  `json:"enabled"`
}

func Init(cfg map[string]Config) {
	cacheConfig = cfg
}

// isCachingEnabled checks if a particular method is cacheable from config
func isCachingEnabled(fnName string) (bool, time.Duration) {
	conf, ok := cacheConfig[fnName]
	if !ok {
		return false, 0
	}
	return conf.Enabled, time.Duration(conf.TtlInSeconds) * time.Second
}

// Cache function takes a func as param along with its arguments, cache key and strategy.
// It calls the passed fn with given args and caches the result with given strategy and returns the result.
// Only works with func having 2 return values, with the second one being an interface.
// Caching will only work if it is enabled in config.
func Cache(ctx context.Context, key string, strategy Strategy, fnc interface{}, args ...interface{}) (interface{}, error) {
	fnType, err := fetchFuncType(fnc)
	if err != nil {
		return nil, err
	}
	fnName := getFuncName(fnc)
	key = fnName + "_" + key
	err = checkForArgsError(args, fnName, fnType)
	if err != nil {
		return nil, err
	}
	if fnType.NumOut() != 2 {
		return nil, errors.New("func: " + fnName + " should have exactly two return values")
	}
	err = checkArgIndexForError(fnType, 1, fnName)
	if err != nil {
		return nil, err
	}
	enabled, expiration := isCachingEnabled(fnName)
	if !enabled {
		return invoke(fnc, args...)
	}
	if strategy == nil {
		panic(fmt.Sprintf("invalid cache strategy:%v", strategy))
	}
	opType := reflect.TypeOf(fnc).Out(0)
	res, err := strategy.Read(ctx, key, opType)
	if err == nil {
		return res, err
	}
	res, err = invoke(fnc, args...)
	strategy.Write(ctx, key, expiration, res, err)
	return res, err
}

// Delete func deletes the given cache key
func Delete(ctx context.Context, key string, strategy Strategy) {
	if strategy == nil {
		panic(fmt.Sprintf("invalid cache strategy:%v", strategy))
	}
	strategy.Delete(ctx, key)
}

// checkArgIndexForError looks for all the values returned by ProcessResponse and checks if the index specified is error.
func checkArgIndexForError(fnType reflect.Type, index int, fnName string) error {
	if !fnType.Out(index).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return fmt.Errorf("second return type should be error, fn:" + fnName)
	}
	return nil
}

// invoke calls the function with relevant args
// Note: This function would panic in case there is a mismatch in the arguments in the source
func invoke(fnc interface{}, args ...interface{}) (interface{}, error) {
	fn := reflect.ValueOf(fnc)
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	res := fn.Call(inputs)
	err, _ := res[1].Interface().(error)
	return res[0].Interface(), err
}

// getFuncName takes function as input and returns the name of the function as string
// This won't work if functions is defined like this var abc = func()
func getFuncName(fnc interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Name()
	endName := filepath.Ext(fullName)
	name := strings.TrimPrefix(endName, ".")
	//this is to remove `-fm` from struct func
	if len(name) > 3 && name[len(name)-3:] == "-fm" {
		name = name[:len(name)-3]
	}
	return name
}

// fetchFuncType checks whether the func passed as a parameter is of type function or not and returns its type
func fetchFuncType(fnc interface{}) (fnType reflect.Type, err error) {
	fn := reflect.ValueOf(fnc)
	fnType = reflect.TypeOf(fnc)
	if fn.Kind() != reflect.Func {
		err = fmt.Errorf("function %v is not of type reflect Func", fnc)
	}
	return
}

// checkForArgsError checks for all possible args error
func checkForArgsError(args []interface{}, fnName string, fnType reflect.Type) error {
	noOfArgs := fnType.NumIn()
	if noOfArgs != len(args) {
		return fmt.Errorf("mismatch in the no of args expected:%v, provided:%v, fn:%v", noOfArgs, len(args), fnName)
	}
	for i := 0; i < noOfArgs; i++ {
		if !reflect.TypeOf(args[i]).AssignableTo(fnType.In(i)) {
			return fmt.Errorf("mismatch in the type of arg no %v, expected:%v, provided:%v, fn:%v", i+1, fnType.In(i), reflect.TypeOf(args[i]), fnName)
		}
	}
	return nil
}
