# Asynchronous-Go
Implementation of Async/Await and Promise.All Using Go's primitives for concurrency (goroutines and channels)

No need for special asynchronous functions when every function could be run asynchronously :D

This is the generic version so run/build the program with ```-gcflags=-G=3``` flag

Usage:

```go
func increaseNum(num int) int {
  time.Sleep(1 * time.Duration)
  return num + 1
}

awaitable := async[int](increaseNum, 5)

// Do Stuff

result, err := awaitable.await()

// result == 6
```
