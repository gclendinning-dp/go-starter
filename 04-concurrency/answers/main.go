package main

import (
	"fmt"
	"sync"
	"time"
)

// SlowFetch simulates a slow network request. It takes 200ms and returns
// a string based on the URL.
func SlowFetch(url string) string {
	time.Sleep(200 * time.Millisecond)
	return "response from " + url
}

// FetchAll runs SlowFetch concurrently for each URL and returns all results.
func FetchAll(urls []string) []string {
	ch := make(chan string)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			ch <- SlowFetch(u)
		}(url)
	}

	// Close the channel once all goroutines are done.
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results from the channel.
	var results []string
	for result := range ch {
		results = append(results, result)
	}

	return results
}

func main() {
	urls := []string{
		"https://example.com/a",
		"https://example.com/b",
		"https://example.com/c",
	}

	results := FetchAll(urls)
	for _, r := range results {
		fmt.Println(r)
	}
}
