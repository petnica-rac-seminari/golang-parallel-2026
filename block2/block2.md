# Block 2 — Reference Sheet

Replace the `time.Sleep` hack from Block 1 with a channel that collects results from each goroutine. Add a timeout using `select` so the program doesn't hang if a URL is slow to respond. Finally, extend the fetching logic to parse each HTML response and extract the links it contains.

## Channels

```go
ch := make(chan string)       // unbuffered
ch := make(chan string, 10)   // buffered, capacity 10

ch <- "hello"        // send (blocks until receiver is ready)
value := <-ch        // receive (blocks until sender is ready)
```

Directional channel types — used in function signatures:

```go
func producer(ch chan<- string) { ch <- "x" }   // send-only
func consumer(ch <-chan string) { v := <-ch }    // receive-only
```

## Select

```go
select {
case v := <-ch1:
    // received from ch1
case ch2 <- value:
    // sent to ch2
case <-time.After(3 * time.Second):
    // neither happened within 3 seconds
}
```

## Wrapping a []byte as an io.Reader

```go
import "bytes"

reader := bytes.NewReader(body) // body is []byte
```

## Parsing HTML

```go
import "golang.org/x/net/html"

doc, err := html.Parse(r) // r is an io.Reader; doc is *html.Node
```

Relevant fields on `*html.Node`:

| Field            | Type               | Description                                                            |
|------------------|--------------------|------------------------------------------------------------------------|
| `n.Type`         | `html.NodeType`    | Kind of node — compare to `html.ElementNode`, `html.TextNode`, etc.    |
| `n.Data`         | `string`           | Tag name for elements (e.g. `"a"`, `"p"`), text content for text nodes |
| `n.Attr`         | `[]html.Attribute` | Attributes of the element; each has `.Key` and `.Val` string fields    |
| `n.FirstChild`   | `*html.Node`       | First child node (`nil` if none)                                       |
| `n.NextSibling`  | `*html.Node`       | Next sibling node (`nil` if none)                                      |

## Checking a string prefix

```go
import "strings"

strings.HasPrefix("/page/1", "/") // true
```
