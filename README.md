# go-futures
This is Future[T] implementation in go
## examples
```go
f1 := future.Async(...)
f2 := f1.Map(...).Map(...).Recover(...).GetUnsafe()
```
