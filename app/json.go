package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/sahilm/fuzzy"
)

type KeySearchEngine struct {
	keys []string
}

func (e KeySearchEngine) Query(query string) []string {
	return FuzzyFind(query, e.keys)
}

func FuzzyFind(query string, candidates []string) []string {
	if len(candidates) == 0 {
		return nil
	}
	if query == "" {
		return candidates
	}
	matches := fuzzy.Find(query, candidates)
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		result = append(result, match.Str)
	}
	return result
}

var (
	logLevel = new(slog.LevelVar)
	logger   = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
)

func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}

func GetKeys(jsonData interface{}) ([]string, error) {
	query, err := gojq.Parse("keys")
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	code, err := gojq.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("failed to compile query: %w", err)
	}

	iter := code.Run(jsonData)
	result, ok := iter.Next()
	if !ok {
		return []string{}, nil
	}

	if err, ok := result.(error); ok {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	keys, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	// Convert []interface{} to []string
	strKeys := make([]string, len(keys))
	for i, key := range keys {
		strKey, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("key at index %d is not a string: %v", i, key)
		}
		strKeys[i] = strKey
	}

	return strKeys, nil
}

func GetUnnestedKeys(jsonData interface{}) ([]string, error) {
	query, err := gojq.Parse(". | paths")
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	code, err := gojq.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("failed to compile query: %w", err)
	}

	iter := code.Run(jsonData)
	joinedPaths := make([]string, 0, 20)
	for {
		data, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := data.(error); ok {
			return nil, fmt.Errorf("got error: %v", err)
		}
		switch v := data.(type) {
		case []interface{}:
			joinedPaths = append(joinedPaths, JoinPath(v))
		default:
			return nil, fmt.Errorf("wrong data type: %T", v)
		}
	}

	// Remove duplicate
	slices.Sort(joinedPaths)
	return slices.Compact(joinedPaths), nil
}

func JoinPath(v []interface{}) string {
	joinedPath := ""
	for _, path := range v {
		switch path.(type) {
		case int:
			joinedPath += "[]"
		default:
			joinedPath += fmt.Sprintf(".%v", path)
		}
	}
	return joinedPath
}

// addTryOperators adds ? after each [] in key that produces an iteration error on jsonData.
// For example, .[].orders[].id becomes .[].orders[]?.id if .orders is missing in some elements.
func addTryOperators(key string, jsonData interface{}) string {
	var result strings.Builder
	i := 0
	for i < len(key) {
		if i+1 < len(key) && key[i] == '[' && key[i+1] == ']' {
			result.WriteString("[]")
			// Check if the prefix up to this [] produces an iteration error
			if hasIterationError(key[:i+2], jsonData) {
				result.WriteByte('?')
			}
			i += 2
		} else {
			result.WriteByte(key[i])
			i++
		}
	}
	return result.String()
}

// hasIterationError checks if running queryStr on jsonData produces any iteration error
// (e.g., "null has no elements" when a key is missing in some array elements).
func hasIterationError(queryStr string, jsonData interface{}) bool {
	query, err := gojq.Parse(queryStr)
	if err != nil {
		return false
	}
	code, err := gojq.Compile(query)
	if err != nil {
		return false
	}
	iter := code.Run(jsonData)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if _, ok := v.(error); ok {
			return true
		}
	}
	return false
}

func QueryJsonData(queryStr string, jsonData interface{}) (interface{}, error) {
	query, err := gojq.Parse(queryStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}
	code, err := gojq.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("failed to compile query: %w", err)
	}
	iter := code.Run(jsonData)
	var results []interface{}
	nonNilCount := 0
	hasError := false
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if _, ok := v.(error); ok {
			// Replace iteration errors with null (e.g., when a JSONL line lacks a key)
			hasError = true
			results = append(results, nil)
			continue
		}
		if v != nil {
			nonNilCount++
		}
		results = append(results, v)
	}

	if len(results) == 0 || (nonNilCount == 0 && !hasError) {
		return "", fmt.Errorf("empty results")
	}
	if nonNilCount == 0 && hasError {
		return "", fmt.Errorf("query produced only errors")
	}
	if len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}

func RobustQueryJsonData(queryStr string, jsonData interface{}) (string, interface{}) {
	queryCandidates := generateQueryCandidates(queryStr)
	for _, query := range queryCandidates {
		results, err := QueryJsonData(query, jsonData)
		logger.Debug(fmt.Sprintf("query: %v and got: %v, %v", query, results, err))
		if err == nil {
			return query, results
		}
	}
	return "", "ERROR: Cannot find"
}

func generateQueryCandidates(queryStr string) []string {
	if queryStr == "keys" {
		return []string{"keys"}
	}
	parts := strings.Split(strings.Trim(queryStr, "."), ".")
	candidates := make([]string, 0, len(parts)+1)
	// Generate queries from most specific to least specific
	// Example: .user.wrongkey -> .user -> .
	for i := 0; i <= len(parts); i++ {
		subParts := parts[:len(parts)-i]
		candidates = append(candidates, "."+strings.Join(subParts, "."))
	}
	logger.Debug("Show generated", "candidates", candidates)
	return candidates
}
