package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessageHandler(t *testing.T) {
	// Create a test server using our handler — no real port needed.
	server := httptest.NewServer(http.HandlerFunc(MessageHandler))
	defer server.Close()

	// Make a request to the test server.
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the Content-Type header.
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	// Read and parse the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	want := "Hello from the backend!"
	if result["message"] != want {
		t.Errorf("message = %q, want %q", result["message"], want)
	}
}
