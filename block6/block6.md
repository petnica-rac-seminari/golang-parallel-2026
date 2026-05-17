# Block 6 — Reference Sheet

Add cancellation to the crawler so it can stop cleanly. Introduce a timeout that automatically stops the crawl after a set duration, and handle Ctrl+C so the user can interrupt it gracefully. Both signals should cancel all in-flight work through the same mechanism — a `context.Context` passed through the program.

## Creating a context with timeout

```go
import (
    "context"
    "time"
)

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

## Creating a context with Ctrl+C handling

```go
import (
    "context"
    "os"
    "os/signal"
)

ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()
```

## Combining timeout and signal cancellation

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
```

Whichever fires first — the timeout or Ctrl+C — cancels the context.

## Making HTTP requests respect a context

```go
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
if err != nil {
    // handle error
}
resp, err := http.DefaultClient.Do(req)
```

The request is cancelled automatically if the context is cancelled.

## Checking if a context has been cancelled

```go
ctx.Err() // returns nil if active, context.Canceled or context.DeadlineExceeded otherwise

if ctx.Err() != nil {
    // context is done
}
```

## Selecting on context cancellation

```go
select {
case v, ok := <-ch:
    if !ok {
        return
    }
    // process v
case <-ctx.Done():
    return // context cancelled or timed out
}
```
