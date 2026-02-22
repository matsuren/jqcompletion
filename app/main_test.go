package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadStdinToFile(t *testing.T) {
	t.Run("reads data and writes to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		dstPath := filepath.Join(tmpDir, "out.json")
		input := `{"key": "value"}`

		err := readStdinToFile(strings.NewReader(input), dstPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}
		if string(got) != input {
			t.Errorf("got %q, want %q", string(got), input)
		}
	})

	t.Run("writes empty input", func(t *testing.T) {
		tmpDir := t.TempDir()
		dstPath := filepath.Join(tmpDir, "out.json")

		err := readStdinToFile(strings.NewReader(""), dstPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}
		if string(got) != "" {
			t.Errorf("got %q, want empty", string(got))
		}
	})

	t.Run("returns error for invalid dest path", func(t *testing.T) {
		err := readStdinToFile(strings.NewReader("data"), "/nonexistent/dir/file.json")
		if err == nil {
			t.Fatal("expected error for invalid path, got nil")
		}
	})
}
