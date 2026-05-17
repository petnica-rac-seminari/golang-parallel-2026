# Stretch — Word Frequency

Build a command-line tool that counts word frequencies across all `.txt` files in a directory and prints the top 20 most common words.

It walks a directory, finds every `.txt` file, and reads them all concurrently using a worker pool. Each worker reads one file and counts the words in it. When all workers are done, the per-file counts are merged into a single global table and the top 20 entries are printed.

```
$ go run .
the        4821
a          3012
and        2987
...
```

Sample dataset: [here](https://drive.google.com/file/d/1Emu6NlPZqpxkaNx0bSJfHTAr0HE3Lroq/view?usp=drive_link)

## Requirements

- Hardcode the directory path as a constant at the top of the file
- Find all `.txt` files in that directory (non-recursive is fine)
- Use a worker pool to read files concurrently — try 4 workers to start
- Each worker counts words in one file and returns a `map[string]int`
- Merge all per-file maps into a single frequency table
- Print the top 20 words sorted by frequency

## Listing files in a directory

```go
import "os"

entries, err := os.ReadDir(dir)
for _, e := range entries {
    if !e.IsDir() {
        name := e.Name() // just the filename, not the full path
    }
}
```

## Reading a file word by word

```go
import (
    "bufio"
    "os"
    "strings"
)

f, err := os.Open(path)
if err != nil {
    // handle error
}
defer f.Close()

scanner := bufio.NewScanner(f)
scanner.Split(bufio.ScanWords) // scan one word at a time instead of one line
for scanner.Scan() {
    word := scanner.Text()
    // word is one word
}
```
