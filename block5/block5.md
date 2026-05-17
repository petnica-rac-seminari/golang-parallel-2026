# Block 5 — Reference Sheet

Extend the crawler with a second stage that processes each fetched page. Instead of handling everything in the fetching logic, send each page's content to a separate pool of workers that compress it and report the compression ratio. The two pools run independently and are connected by a channel — this is the pipeline pattern.

## gzip compression

```go
import (
    "bytes"
    "compress/gzip"
)

var buf bytes.Buffer
w := gzip.NewWriter(&buf)
w.Write(body)
w.Close()

// buf.Bytes() contains the compressed data
// buf.Len() is the compressed size
```

## Number of CPU cores

```go
import "runtime"

runtime.NumCPU() // returns int
```

Useful for deciding how many CPU-bound workers to run.
