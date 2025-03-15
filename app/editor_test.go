package main

import (
	"os"
	"testing"
)

func TestEditFile(t *testing.T) {
	jsonPath, err := getTempJsonPath()
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	defer os.Remove(jsonPath)
	editFileExecCmd(jsonPath)
	// cmd.Run() will block
}
