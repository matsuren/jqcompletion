package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type editorFinishedMsg struct {
	tempJsonPath string
	err          error
}

func editFileExecCmd(filepath string) *exec.Cmd {
	return editorFileExecCmd(filepath, false)
}

func editorFileExecCmd(filepath string, readOnly bool) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	var c *exec.Cmd
	if readOnly && (editor == "vim" || editor == "nvim") {
		c = exec.Command(editor, "-R", filepath)
	} else {
		c = exec.Command(editor, filepath)
	}
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	return c
}

func getTempJsonPath() (string, error) {
	f, err := os.CreateTemp("", fmt.Sprintf("tmp_jqc%s_*.json", time.Now().Format("2006-01-02T15:04:05")))
	defer func() { _ = f.Close() }()
	return f.Name(), err
}

func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer func() { _ = sourceFile.Close() }()

	// Create or truncate the destination file
	// The os.O_TRUNC flag ensures any existing file is truncated (overwritten)
	destFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return fmt.Errorf("open dest file: %w", err)
	}
	defer func() { _ = destFile.Close() }()

	// Copy the contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("copy files: %w", err)
	}

	return nil
}

func LoadJsonFile(jsonPath string) (interface{}, error) {
	logger.Debug("Loading:", "jsonPath", jsonPath)
	// Read the JSON file
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	// Parse the JSON
	logger.Debug("Parsing:", "jsonPath", jsonPath)
	var rawJsonData interface{}
	err = json.Unmarshal(jsonData, &rawJsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	logger.Debug("Done LoadJsonFile")
	return rawJsonData, nil
}

func SaveJsonToFile(jsonPath string, jsonData interface{}) error {
	// Convert jsonData to formatted JSON bytes
	resultBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}
	err = os.WriteFile(jsonPath, resultBytes, 0o644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
