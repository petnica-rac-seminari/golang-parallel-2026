package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type CrawlResult struct {
	Url    string
	Length int
	Links  []string
}

type GzipJob struct {
	Url  string
	Data []byte
}

var (
	workGroup         sync.WaitGroup
	visited           = make(map[string]bool)
	visitedLock       sync.Mutex
	fetchJobs         = make(chan string, 100)
	gzipJobs          = make(chan GzipJob, 100)
	gzip_worker_count = runtime.NumCPU()
)

const (
	domain_name        = "https://golang-parallel.fly.dev"
	fetch_worker_count = 50
)

func extractLinks(reader io.Reader) []string {
	var links []string

	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("Error in parsing HTML: %v", err)
		return nil
	}

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Data == "a" { //this node is an a tag
			for _, attr := range node.Attr { //range over all attributes
				if attr.Key == "href" { //find href attributes
					var linkResult string
					if strings.HasPrefix(attr.Val, "/") { //relative path
						linkResult = domain_name + attr.Val
					} else { //absolute path
						linkResult = attr.Val
					}
					links = append(links, linkResult)
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(doc)

	return links
}

func crawl(url string, ctx context.Context) {
	defer workGroup.Done()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Error in creating HTTP request %s: %v\n", url, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Error in fetching url %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error in reading body for %s: %v\n", url, err)
		return
	}

	if ctx.Err() != nil {
		return
	}

	links := extractLinks(bytes.NewReader(body))

	for _, link := range links {
		visitedLock.Lock()
		if !visited[link] {
			visited[link] = true
			workGroup.Add(1)
			fetchJobs <- link
		}
		visitedLock.Unlock()
	}

	//pipeline pattern
	gzipJobs <- GzipJob{Url: url, Data: body}

}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start_url := "https://golang-parallel.fly.dev/page/1"

	//download workers
	for range fetch_worker_count {
		go func() {
			for job := range fetchJobs {
				crawl(job, ctx)
			}
		}()
	}

	//gzip workers
	var gzipWorkGroup sync.WaitGroup
	for range gzip_worker_count {
		gzipWorkGroup.Add(1)
		go func() {
			defer gzipWorkGroup.Done()
			for {
				select {
				case toGzip, ok := <-gzipJobs:
					if !ok {
						return
					}
					var buf bytes.Buffer
					w := gzip.NewWriter(&buf)
					w.Write(toGzip.Data)
					w.Close()

					ratio := (float64(buf.Len()) / float64(len(toGzip.Data))) * 100.0
					fmt.Printf("Zipped %s: %d -> %d %.1f%%\n", toGzip.Url, len(toGzip.Data), buf.Len(), ratio)
				case <-ctx.Done():
					return
				}
			}

		}()
	}

	visited[start_url] = true
	workGroup.Add(1)
	fetchJobs <- start_url

	workGroup.Wait()
	close(fetchJobs)
	close(gzipJobs)
	gzipWorkGroup.Wait()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Printf("Timeout reached. Ending gracefully.\n")
	}
	if ctx.Err() == context.Canceled {
		fmt.Printf("User interrupt. Ending gracefully.\n")
	}

	fmt.Printf("Crawled and zipped %d pages with %d cpu workers\n", len(visited), gzip_worker_count)

}
