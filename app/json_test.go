package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestFuzzyFind(t *testing.T) {
	query := "banana"
	candidates := []string{"apple", "banana", "cherry"}
	expected := []string{"banana"}

	result := FuzzyFind(query, candidates)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FuzzyFind(%q, %v) = %v, want %v", query, candidates, result, expected)
	}
}

func TestGetKeys(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		"status": "active",
	}
	want := []string{"status", "user"}
	wantErr := false

	got, err := GetKeys(input)

	// Check error cases
	if (err != nil) != wantErr {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
		return
	}

	if wantErr {
		return
	}

	// Sort both slices for comparison
	sort.Strings(got)
	sort.Strings(want)

	// Use reflect.DeepEqual for comparison
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestGetUnnestedKeys(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		"status": "active",
	}
	want := []string{".status", ".user", ".user.age", ".user.name"}

	got, err := GetUnnestedKeys(input)
	// Check error cases
	if err != nil {
		t.Errorf("Got error: %v", err)
		return
	}

	// Sort both slices for comparison
	sort.Strings(got)
	sort.Strings(want)

	// Use reflect.DeepEqual for comparison
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestJoinPath(t *testing.T) {
	input := []interface{}{"status", "user", 0, "id"}
	want := ".status.user[].id"
	got := JoinPath(input)

	if got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestQueryJsonData(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		"status": "active",
	}
	got, err := QueryJsonData(".", input)
	if err != nil {
		t.Errorf("Got err: %v", err)
	}

	if !reflect.DeepEqual(input, got) {
		t.Errorf("Got %v, want %v", got, input)
	}
}

func TestRobustQueryJsonData(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		"status": "active",
	}

	want := struct {
		result map[string]interface{}
		query  string
	}{
		result: map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		query: ".user",
	}

	gotQuery, got := RobustQueryJsonData(".user.wrongkey", input)
	if !reflect.DeepEqual(got, want.result) {
		t.Errorf("Got %v, want %v", got, want)
	}

	if !reflect.DeepEqual(gotQuery, want.query) {
		t.Errorf("Got %v, want %v", gotQuery, want.query)
	}
}

func TestRobustQueryJsonData2(t *testing.T) {
	SetLogLevel(slog.LevelInfo)

	filePath := "../data/sample.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON
	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	inputQuery := ".results[].aaa"

	want := struct {
		result interface{}
		query  string
	}{
		result: []map[string]string{
			{
				"level":   "Info",
				"message": "Info message",
			},
			{
				"level":   "Debug",
				"message": "Debug message",
			},
		},
		query: ".results[]",
	}

	gotQuery, got := RobustQueryJsonData(inputQuery, data)
	if !reflect.DeepEqual(gotQuery, want.query) {
		t.Errorf("Got %v, want %v", gotQuery, want.query)
	}
	if !reflect.DeepEqual(jsonDataToStrings(got), jsonDataToStrings(want.result)) {
		t.Errorf("Got %v, want %v", got, want.result)
	}
}

func jsonDataToStrings(jsonData interface{}) string {
	resultBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "error"
	}
	return string(resultBytes)
}

func TestLoadJsonFile(t *testing.T) {
	t.Run("Valid JSON File", func(t *testing.T) {
		// Create a temporary JSON file
		tmpFile, err := os.CreateTemp("", "jqctest*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		jsonContent := `{"key": "value"}`
		if _, err := tmpFile.Write([]byte(jsonContent)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		// Call the function
		result, err := LoadJsonFile(tmpFile.Name())
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify the result
		expected := map[string]interface{}{"key": "value"}
		if result == nil {
			t.Errorf("Expected non-nil result")
		}

		if !reflect.DeepEqual(jsonDataToStrings(result), jsonDataToStrings(expected)) {
			t.Errorf("Got %v, want %v", result, expected)
		}
	})

	t.Run("Invalid JSON File", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "invalid.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte("invalid json")); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		_, err = LoadJsonFile(tmpFile.Name())
		if err == nil {
			t.Errorf("Expected an error for invalid JSON but got none")
		}
	})

	t.Run("Non-Existent File", func(t *testing.T) {
		_, err := LoadJsonFile("non_existent_file.json")
		if err == nil {
			t.Errorf("Expected an error for non-existent file but got none")
		}
	})
}
