package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockMessageHandler replicates the server from Task 02 so this test is
// self-contained — no need to start a separate process.
func mockMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Hello from the backend!"}`)
}

func TestFetchMessage(t *testing.T) {
	// Spin up a test server with the mock handler.
	server := httptest.NewServer(http.HandlerFunc(mockMessageHandler))
	defer server.Close()

	// Call FetchMessage against the test server.
	got, err := FetchMessage(server.URL)
	if err != nil {
		t.Fatalf("FetchMessage returned error: %v", err)
	}

	want := "Hello from the backend!"
	if got != want {
		t.Errorf("FetchMessage() = %q, want %q", got, want)
	}
}

func TestFetchMessageBadURL(t *testing.T) {
	_, err := FetchMessage("http://localhost:1/nope")
	if err == nil {
		t.Error("expected an error for a bad URL, got nil")
	}
}
