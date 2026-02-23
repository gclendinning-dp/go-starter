package main

import (
	"fmt"
	"net/http"
)

// MessageHandler writes a JSON response with a greeting message.
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "Hello from the backend!"}`)
}

func main() {
	http.HandleFunc("/message", MessageHandler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
