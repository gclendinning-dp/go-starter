package main

import "fmt"

// Greet returns a welcome message for the given name.
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s! Welcome to Go.", name)
}

func main() {
	message := Greet("Student")
	fmt.Println(message)
}
