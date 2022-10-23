go-cache
========== 
A light-weight caching library for Go. It can cache the results of a function with given TTL.

Usage
----------------
Use any of the example strategies like -
```
InApp,
InRedis,
or implement your own custom strategy.
```
Init a strategy and use it in Cache function as -
```go
strategy := NewInApp()
res, err := Cache(ctx, key, strategy, fn, args...)
```
see `examples_test.go`

License
----------------
MIT