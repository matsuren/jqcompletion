package main

import (
	"os"
	"reflect"
	"testing"
)

func TestEditFile(t *testing.T) {
	jsonPath, err := getTempJsonPath()
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	defer func() { _ = os.Remove(jsonPath) }()
	editFileExecCmd(jsonPath)
	// cmd.Run() will block
}

func TestEditorCopyFile(t *testing.T) {
	jsonPath, err := getTempJsonPath()
	defer func() { _ = os.Remove(jsonPath) }()
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
	defer func() { _ = os.Remove(jsonPath) }()
	if err != nil {
		t.Fatal(err)
	}
	jsonData, _, err := LoadJsonFile("../data/sample.json")
	if err != nil {
		t.Fatal(err)
	}
	err = SaveJsonToFile(jsonPath, jsonData)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadJsonFile_JSONL(t *testing.T) {
	data, isSlurp, err := LoadJsonFile("../data/sample.jsonl")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !isSlurp {
		t.Error("Expected isSlurp=true for JSONL input")
	}
	arr, ok := data.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", data)
	}
	if len(arr) != 3 {
		t.Errorf("Expected 3 items, got %d", len(arr))
	}
	first, ok := arr[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map for first item, got %T", arr[0])
	}
	if first["name"] != "Alice" {
		t.Errorf("Expected name=Alice, got %v", first["name"])
	}
}

func TestLoadJsonFile_SingleJSON(t *testing.T) {
	data, isSlurp, err := LoadJsonFile("../data/sample.json")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if isSlurp {
		t.Error("Expected isSlurp=false for single JSON input")
	}
	m, ok := data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", data)
	}
	if m["tag"] != "Tag" {
		t.Errorf("Expected tag=Tag, got %v", m["tag"])
	}
}

func TestLoadJsonFile_SingleLineJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "single_line*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	content := `{"key": "value"}`
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	_ = tmpFile.Close()

	data, isSlurp, err := LoadJsonFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if isSlurp {
		t.Error("Expected isSlurp=false for single-line JSON")
	}
	expected := map[string]interface{}{"key": "value"}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Got %v, want %v", data, expected)
	}
}
