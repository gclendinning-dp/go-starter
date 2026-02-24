package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMessageHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /message", MessageHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/message")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	want := "Hello from the backend!"
	if result["message"] != want {
		t.Errorf("message = %q, want %q", result["message"], want)
	}
}

func TestGreetHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /greet", GreetHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	body := strings.NewReader(`{"name": "Alice"}`)
	resp, err := http.Post(server.URL+"/greet", "application/json", body)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var result GreetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	want := "Hello, Alice!"
	if result.Greeting != want {
		t.Errorf("greeting = %q, want %q", result.Greeting, want)
	}
}

func TestGreetHandlerEmptyName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /greet", GreetHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	body := strings.NewReader(`{"name": ""}`)
	resp, err := http.Post(server.URL+"/greet", "application/json", body)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestGreetHandlerBadJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /greet", GreetHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	body := strings.NewReader(`{bad`)
	resp, err := http.Post(server.URL+"/greet", "application/json", body)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}
