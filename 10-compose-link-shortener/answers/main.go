package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

// LinkStore manages shortened links using Redis.
type LinkStore struct {
	rdb *redis.Client
	ctx context.Context
}

// NewLinkStore creates a LinkStore connected to the given Redis address.
func NewLinkStore(redisAddr string) *LinkStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &LinkStore{
		rdb: rdb,
		ctx: context.Background(),
	}
}

// Shorten stores a URL in Redis and returns its short key.
func (s *LinkStore) Shorten(url string) (string, error) {
	// INCR gives us a unique integer key (starts at 1)
	next, err := s.rdb.Incr(s.ctx, "link:next").Result()
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%d", next)

	err = s.rdb.HSet(s.ctx, "links", key, url).Err()
	if err != nil {
		return "", err
	}

	return key, nil
}

// Lookup returns the URL for a key and whether it was found.
func (s *LinkStore) Lookup(key string) (string, bool) {
	url, err := s.rdb.HGet(s.ctx, "links", key).Result()
	if err != nil {
		return "", false
	}
	return url, true
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

	key, err := s.Shorten(req.URL)
	if err != nil {
		http.Error(w, "failed to shorten", http.StatusInternalServerError)
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
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	store := NewLinkStore(redisAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", store.ShortenHandler)
	mux.HandleFunc("GET /r/{key}", store.RedirectHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
