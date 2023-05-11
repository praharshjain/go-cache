go-cache
========== 
A light-weight caching library for Go. It can cache the results of a function with given TTL.

Known caveats
----------------
1. Return value of the function to be cached should be exactly a result (of any type) and an error
2. Function to be cached should not be declared using `var fn = func` syntax but `func fn` instead. (This is because the library uses reflection to fetch the function name)
3. Data in the `InApp` strategy just expires, but is never deleted even after TTL. So it should be used responsibly.

Usage
----------------
1. Init config (only to be done once, generally at the start of your service)
```go
//call Init with config map (refer config.json for sample)
Init(cfg)
````
2. Instantiate a strategy (strategy is meant to be reused across calls to Cache)
```go
//strategy is to instantiated only once and not on every call to Cache()
strategy := NewInApp()
```
3. Use the strategy in Cache function as
```go
res, err := Cache(ctx, key, strategy, fn, args...)
```
4. Then you need to type case the result into the expected return type of the function
```go
s, _ := res.(string)
```

Use any of the available strategies like -
```
InApp,
InRedis,
or implement your own custom strategy.
```
see `examples_test.go` for more