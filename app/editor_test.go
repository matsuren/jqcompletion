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

func TestEditorCopyFile(t *testing.T) {
	jsonPath, err := getTempJsonPath()
	defer os.Remove(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	err = copyFile("../data/sample.json", jsonPath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEditorSaveJsonToFile(t *testing.T) {
	jsonPath, err := getTempJsonPath()
	defer os.Remove(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := LoadJsonFile("../data/sample.json")
	if err != nil {
		t.Fatal(err)
	}
	err = SaveJsonToFile(jsonPath, jsonData)
	if err != nil {
		t.Fatal(err)
	}
}
