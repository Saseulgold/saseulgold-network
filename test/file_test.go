package main

import (
	"bytes"
	"encoding/hex"
	"hello/pkg/core/config"
	"hello/pkg/core/storage"
	"hello/pkg/util"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	owner := config.ZeroAddress()
	space := config.RootSpace()

	hash := util.StatusHash(owner, space, "balance", owner)
	t.Logf("hash: %s", hash)
}

// TestListFiles is a test function for StorageFileIndex's ListFiles method.
func TestListFiles(t *testing.T) {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(tempDir) // Clean up after the test

	// Create test files
	files := []struct {
		name string
	}{
		{"prefix_file1.txt"},
		{"prefix_file2.txt"},
		{"other_file.txt"},
	}

	for _, file := range files {
		err := ioutil.WriteFile(tempDir+"/"+file.name, []byte("test data"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file.name, err)
		}
	}

	// Initialize StorageFileIndex
	storageIndex := storage.StorageFileIndex{
		Prefix: "prefix",
		Dir:    tempDir,
	}

	// Call ListFiles method
	matchedFiles := storageIndex.ListFiles()

	// Expected result
	expectedFiles := []string{"prefix_file1.txt", "prefix_file2.txt"}

	// Verify the result
	if len(matchedFiles) != len(expectedFiles) {
		t.Errorf("Expected %d files, but got %d. Something went off the rails!", len(expectedFiles), len(matchedFiles))
	}

	// Validate the file list
	for _, expectedFile := range expectedFiles {
		found := false
		for _, matchedFile := range matchedFiles {
			if expectedFile == matchedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file '%s' was not found in the result list. Make sure it's included!", expectedFile)
		}
	}
}

func TestStatusKey(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    strings.Repeat("a", 128),
			expected: hex.EncodeToString([]byte(strings.Repeat("a", 64))),
		},
		{
			input:    "short",
			expected: "",
		},
	}

	for _, tc := range testCases {
		result := storage.StatusKey(tc.input)
		if result != tc.expected {
			t.Errorf("StatusKey(%s) = %s; want %s", tc.input, result, tc.expected)
		}
	}
}

func TestNewStatusIndex(t *testing.T) {
	// Create test data with proper length
	raw := strings.Repeat("a", 64) + // Status key bytes
		"bb" +
		"cccc" +
		"dddd"

	index := storage.NewStatusIndex(raw)

	if index.FileID != "bb" {
		t.Errorf("Expected FileID 'bb', got %s", index.FileID)
	}

	expectedSeek := util.Hex2Bin("cccc")
	if !bytes.Equal(index.Seek, expectedSeek) {
		t.Errorf("Expected Seek %v, got %v", expectedSeek, index.Seek)
	}

	expectedLength := util.Hex2Bin("dddd")
	if !bytes.Equal(index.Length, expectedLength) {
		t.Errorf("Expected Length %v, got %v", expectedLength, index.Length)
	}
}

func TestSplitKey(t *testing.T) {
	testCases := []struct {
		input          string
		expectedPrefix string
		expectedSuffix string
	}{
		{
			input:          "12345678",
			expectedPrefix: "1234",
			expectedSuffix: "5678",
		},
		{
			input:          "123",
			expectedPrefix: "",
			expectedSuffix: "",
		},
	}

	for _, tc := range testCases {
		result := storage.SplitKey(tc.input)
		if result[0] != tc.expectedPrefix || result[1] != tc.expectedSuffix {
			t.Errorf("SplitKey(%s) = [%s, %s]; want [%s, %s]",
				tc.input, result[0], result[1], tc.expectedPrefix, tc.expectedSuffix)
		}
	}
}
