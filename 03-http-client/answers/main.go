package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response holds the JSON structure returned by the server.
type Response struct {
	Message string `json:"message"`
}

// FetchMessage makes a GET request to the given URL and returns the message.
func FetchMessage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body failed: %w", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing JSON failed: %w", err)
	}

	return result.Message, nil
}

func main() {
	message, err := FetchMessage("http://localhost:8080/message")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(message)
}
