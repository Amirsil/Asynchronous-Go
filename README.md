# Asynchronous-Go
Implementation of Async/Await and Promise.All Using Go's primitives for concurrency (goroutines and channels)

No need for special asynchronous functions when every function could be run asynchronously :D

Usage:

```go
func increaseNum(num int) int {
  time.Sleep(1 * time.Duration)
  return num + 1
}

awaitable := Async(increaseNum, 5)

// Do Stuff

result, err := awaitable.Await()

// result == 6
```

You could also use Promise like notation:

```go
result, err := Async(increaseNum, 5).
                  Then(func(num int) (int, error) { return num - 1, errors.New("Wow! Thats an error!") }).
                  Catch(func(err error) { log.Fatal(err) }).
                  Then(func(num int) int { return num }).
                  Await()

// result == 5
```
