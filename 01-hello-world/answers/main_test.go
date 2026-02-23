package main

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "Alice", want: "Hello, Alice! Welcome to Go."},
		{name: "World", want: "Hello, World! Welcome to Go."},
		{name: "", want: "Hello, ! Welcome to Go."},
	}

	for _, tt := range tests {
		got := Greet(tt.name)
		if got != tt.want {
			t.Errorf("Greet(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
