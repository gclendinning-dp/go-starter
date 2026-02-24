package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadLines(t *testing.T) {
	// Write a temp file to read back.
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	content := "line one\n\nline two\nline three\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	lines, err := ReadLines(path)
	if err != nil {
		t.Fatalf("ReadLines returned error: %v", err)
	}

	want := []string{"line one", "line two", "line three"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d", len(lines), len(want))
	}

	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("lines[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestFilterLines(t *testing.T) {
	lines := []string{
		"2024-01-15 10:00:01 INFO: Server started",
		"2024-01-15 10:01:12 ERROR: Connection failed",
		"2024-01-15 10:01:14 INFO: Connection restored",
		"2024-01-15 10:02:30 ERROR: Request timeout",
	}

	filtered := FilterLines(lines, "ERROR")

	if len(filtered) != 2 {
		t.Fatalf("got %d filtered lines, want 2", len(filtered))
	}

	if filtered[0] != lines[1] {
		t.Errorf("filtered[0] = %q, want %q", filtered[0], lines[1])
	}
	if filtered[1] != lines[3] {
		t.Errorf("filtered[1] = %q, want %q", filtered[1], lines[3])
	}
}

func TestWriteLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")

	lines := []string{"alpha", "bravo", "charlie"}

	if err := WriteLines(path, lines); err != nil {
		t.Fatalf("WriteLines returned error: %v", err)
	}

	// Read the file back and verify.
	got, err := ReadLines(path)
	if err != nil {
		t.Fatalf("ReadLines returned error: %v", err)
	}

	if len(got) != len(lines) {
		t.Fatalf("got %d lines, want %d", len(got), len(lines))
	}

	for i := range lines {
		if got[i] != lines[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], lines[i])
		}
	}
}

func TestReadLinesNotFound(t *testing.T) {
	_, err := ReadLines("/tmp/nonexistent_file_that_does_not_exist.txt")
	if err == nil {
		t.Error("expected an error for a nonexistent file, got nil")
	}
}
