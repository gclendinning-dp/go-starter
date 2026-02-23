package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadLines reads a file and returns all non-empty lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

// FilterLines returns only the lines that contain the given keyword.
func FilterLines(lines []string, keyword string) []string {
	var filtered []string
	for _, line := range lines {
		if strings.Contains(line, keyword) {
			filtered = append(filtered, line)
		}
	}
	return filtered
}

// WriteLines creates a file and writes each line into it.
func WriteLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		fmt.Fprintln(file, line)
	}

	return nil
}

func main() {
	lines, err := ReadLines("server.log")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	errors := FilterLines(lines, "ERROR")

	if err := WriteLines("errors.txt", errors); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Printf("Found %d errors, written to errors.txt\n", len(errors))
}
