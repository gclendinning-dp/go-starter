package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

// LinkStore holds shortened links in memory with file-based persistence.
type LinkStore struct {
	mu    sync.Mutex
	links map[string]string
	next  int
	file  string
}

// NewLinkStore creates a LinkStore that persists to the given file path.
func NewLinkStore(file string) *LinkStore {
	return &LinkStore{
		links: make(map[string]string),
		file:  file,
	}
}

// Shorten stores a URL and returns its short key.
func (s *LinkStore) Shorten(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%d", s.next)
	s.links[key] = url
	s.next++
	return key
}

// Lookup returns the URL for a key and whether it was found.
func (s *LinkStore) Lookup(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, ok := s.links[key]
	return url, ok
}

// SaveToFile writes the link map to disk as JSON.
func (s *LinkStore) SaveToFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(s.links, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.file, data, 0644)
}

// LoadFromFile reads the link map from disk. If the file doesn't exist, it
// does nothing (a fresh start is fine).
func (s *LinkStore) LoadFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &s.links); err != nil {
		return err
	}

	s.next = len(s.links)
	return nil
}

// ShortenRequest holds the JSON body for the shorten endpoint.
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse holds the JSON response from the shorten endpoint.
type ShortenResponse struct {
	Key      string `json:"key"`
	ShortURL string `json:"short_url"`
}

// ShortenHandler accepts a URL and returns a shortened key.
func (s *LinkStore) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	key := s.Shorten(req.URL)

	if err := s.SaveToFile(); err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenResponse{
		Key:      key,
		ShortURL: fmt.Sprintf("http://localhost:8080/r/%s", key),
	})
}

// RedirectHandler looks up a key and redirects to the original URL.
func (s *LinkStore) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	url, ok := s.Lookup(key)
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	store := NewLinkStore("links.json")

	if err := store.LoadFromFile(); err != nil {
		fmt.Println("Error loading links:", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", store.ShortenHandler)
	mux.HandleFunc("GET /r/{key}", store.RedirectHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
