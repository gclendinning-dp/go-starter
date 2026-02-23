package main

import (
	"sort"
	"testing"
	"time"
)

func TestFetchAllCorrectness(t *testing.T) {
	urls := []string{
		"https://example.com/a",
		"https://example.com/b",
		"https://example.com/c",
	}

	results := FetchAll(urls)

	if len(results) != len(urls) {
		t.Fatalf("got %d results, want %d", len(results), len(urls))
	}

	// Sort both so order doesn't matter.
	sort.Strings(results)
	want := []string{
		"response from https://example.com/a",
		"response from https://example.com/b",
		"response from https://example.com/c",
	}
	sort.Strings(want)

	for i := range want {
		if results[i] != want[i] {
			t.Errorf("results[%d] = %q, want %q", i, results[i], want[i])
		}
	}
}

func TestFetchAllConcurrency(t *testing.T) {
	urls := []string{
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
		"https://example.com/4",
		"https://example.com/5",
	}

	start := time.Now()
	results := FetchAll(urls)
	elapsed := time.Since(start)

	if len(results) != 5 {
		t.Fatalf("got %d results, want 5", len(results))
	}

	// Each SlowFetch takes 200ms. If run sequentially, 5 would take >=1s.
	// If run concurrently, they should all finish in ~200ms.
	// We allow up to 500ms to be safe on slow CI machines.
	if elapsed >= 500*time.Millisecond {
		t.Errorf("FetchAll took %v; want < 500ms (are goroutines running concurrently?)", elapsed)
	}
}

func TestFetchAllEmpty(t *testing.T) {
	results := FetchAll([]string{})

	if len(results) != 0 {
		t.Errorf("got %d results for empty input, want 0", len(results))
	}
}
