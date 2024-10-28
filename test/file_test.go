package storage

import (
	"hello/pkg/core/storage"
	"io/ioutil"
	"os"
	"testing"
)

// TestListFiles is a test function for StorageFileIndex's ListFiles method.
func TestListFiles(t *testing.T) {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(tempDir) // Clean up after the test
	t.Logf("Temporary directory created at: %s", tempDir)

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
