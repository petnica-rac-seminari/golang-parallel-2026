# Block 3 — Reference Sheet

Extend the program to follow discovered links recursively. Instead of fetching a fixed list, start from a single URL, extract links from each page, and crawl those as well. Use a `sync.WaitGroup` to know when all goroutines are done, and a `sync.Mutex` to protect a visited map so the same URL is never crawled twice.

## sync.WaitGroup

```go
import "sync"

var wg sync.WaitGroup

wg.Add(1)      // increment before launching a goroutine
go func() {
    defer wg.Done() // decrement when the goroutine finishes
    // ...
}()

wg.Wait() // block until the counter reaches zero
```

`wg.Add(1)` must be called **before** `go`, not inside the goroutine.

## sync.Mutex

```go
import "sync"

var mu sync.Mutex

mu.Lock()
// only one goroutine can be here at a time
mu.Unlock()

// or with defer
mu.Lock()
defer mu.Unlock()
```

## sync.RWMutex (stretch)

Allows multiple concurrent readers, but only one writer at a time.

```go
var mu sync.RWMutex

mu.RLock()   // acquire read lock
mu.RUnlock()

mu.Lock()    // acquire write lock
mu.Unlock()
```

## Counting active goroutines (stretch)

```go
import "runtime"

runtime.NumGoroutine() // returns int
```
