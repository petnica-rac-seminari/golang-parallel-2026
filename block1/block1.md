# Block 1 — Reference Sheet

Write a program that fetches a list of URLs and prints the number of bytes received from each one along with the total time taken. Start by fetching them sequentially and observe how long it takes. Then make all fetches run concurrently using goroutines, and use `time.Sleep` to keep the program alive long enough for the goroutines to finish.

## Fetching a URL

```go
import (
    "io"
    "net/http"
)

resp, err := http.Get("http://example.com")
if err != nil {
    // handle error
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
// body is []byte
fmt.Println(len(body)) // number of bytes received
```

## Launching a goroutine

```go
go someFunction()
go someFunction(arg1, arg2)

// inline with an anonymous function
go func() {
    // ...
}()
```

## Pausing execution

```go
import "time"

time.Sleep(2 * time.Second)
time.Sleep(500 * time.Millisecond)
```

## Measuring elapsed time

```go
import "time"

start := time.Now()
// ... do work ...
fmt.Println(time.Since(start)) // e.g. "1.234s"
```

## Running with the race detector

```
go run -race ./cmd/yourprogram
```

Reports concurrent accesses to shared variables, even when the output looks correct.

## Stretch — controlling the number of OS threads

```go
import "runtime"

runtime.GOMAXPROCS(1) // restrict to a single OS thread
runtime.GOMAXPROCS(runtime.NumCPU()) // default: one thread per CPU core
```
