package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// MessageResponse holds the JSON body for the message endpoint.
type MessageResponse struct {
	Message  string `json:"message"`
	Hostname string `json:"hostname"`
}

// GreetRequest holds the JSON body sent by the client.
type GreetRequest struct {
	Name string `json:"name"`
}

// GreetResponse holds the JSON body sent back to the client.
type GreetResponse struct {
	Greeting string `json:"greeting"`
	Hostname string `json:"hostname"`
}

// MessageHandler writes a JSON response with a greeting and the container hostname.
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{
		Message:  fmt.Sprintf("Hello from %s!", hostname),
		Hostname: hostname,
	})
}

// GreetHandler reads a name from the request body and responds with a greeting.
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	var req GreetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	hostname, _ := os.Hostname()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s! From %s", req.Name, hostname),
		Hostname: hostname,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /message", MessageHandler)
	mux.HandleFunc("POST /greet", GreetHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
