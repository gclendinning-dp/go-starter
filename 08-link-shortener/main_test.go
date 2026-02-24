package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	store := NewLinkStore(filepath.Join(t.TempDir(), "links.json"))

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", store.ShortenHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	body := strings.NewReader(`{"url": "https://go.dev"}`)
	resp, err := http.Post(server.URL+"/shorten", "application/json", body)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	var result ShortenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if result.Key == "" {
		t.Error("expected a non-empty key")
	}
}

func TestRedirectHandler(t *testing.T) {
	store := NewLinkStore(filepath.Join(t.TempDir(), "links.json"))
	key := store.Shorten("https://go.dev")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /r/{key}", store.RedirectHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Use a client that does not follow redirects.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(server.URL + "/r/" + key)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusFound)
	}

	location := resp.Header.Get("Location")
	if location != "https://go.dev" {
		t.Errorf("Location = %q, want %q", location, "https://go.dev")
	}
}

func TestRedirectNotFound(t *testing.T) {
	store := NewLinkStore(filepath.Join(t.TempDir(), "links.json"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /r/{key}", store.RedirectHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/r/nonexistent")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestShortenEmptyURL(t *testing.T) {
	store := NewLinkStore(filepath.Join(t.TempDir(), "links.json"))

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", store.ShortenHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	body := strings.NewReader(`{"url": ""}`)
	resp, err := http.Post(server.URL+"/shorten", "application/json", body)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestPersistence(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "links.json")

	// Create a store, add a link, and save.
	store1 := NewLinkStore(file)
	key := store1.Shorten("https://go.dev")

	if err := store1.SaveToFile(); err != nil {
		t.Fatalf("SaveToFile returned error: %v", err)
	}

	// Create a new store from the same file and load.
	store2 := NewLinkStore(file)
	if err := store2.LoadFromFile(); err != nil {
		t.Fatalf("LoadFromFile returned error: %v", err)
	}

	url, ok := store2.Lookup(key)
	if !ok {
		t.Fatalf("key %q not found after loading from file", key)
	}
	if url != "https://go.dev" {
		t.Errorf("url = %q, want %q", url, "https://go.dev")
	}
}
