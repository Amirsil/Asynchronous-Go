# Asynchronous-Go
Implementation of Async/Await and Promise.All Using Go's primitives for concurrency (goroutines and channels)

No need for special asynchronous functions when every function could be run asynchronously :D

Usage:
```go
func IncreaseNum(num int) {
  time.Sleep(1 * time.Duration)
  return num + 1
}

awaitable := Async(increaseNum, 5)
// Do Stuff

result, err := awaitable.Await()
```
