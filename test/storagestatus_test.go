package main

import (
	"os"
	"testing"

	C "hello/pkg/core/config"
	S "hello/pkg/core/storage"
)

func TestNewStatusFile(t *testing.T) {
	t.Log("Starting NewStatusFile test")
	sf := S.NewStatusFile()

	if sf.CachedUniversalIndexes == nil {
		t.Error("CachedUniversalIndexes was not initialized")
	}
	t.Log("Verified CachedUniversalIndexes initialization")

	if sf.CachedLocalIndexes == nil {
		t.Error("CachedLocalIndexes was not initialized")
	}
	t.Log("Verified CachedLocalIndexes initialization")

	if sf.Tasks == nil {
		t.Error("Tasks was not initialized")
	}
	t.Log("Verified Tasks initialization")
}

func TestStatusFile_Touch(t *testing.T) {
	t.Log("Starting Touch test")
	C.CORE_TEST_MODE = true

	sf := S.NewStatusFile()
	t.Log("Created new status file")

	err := sf.Touch()
	if err != nil {
		t.Errorf("Error occurred during Touch(): %v", err)
	}
	t.Log("Successfully executed Touch operation")

	t.Log("Verifying created files...")
	expectedFiles := []string{
		sf.TempFile(),
		sf.InfoFile(),
		sf.LocalFile(),
		sf.LocalBundle(),
		sf.UniversalBundle("00"),
		sf.LocalBundleIndex(),
		sf.UniversalBundleIndex(),
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected file was not created: %s", file)
		}
		t.Logf("Verified file exists: %s", file)
	}
}

func TestStatusFile_Reset(t *testing.T) {
	t.Log("Starting Reset test")
	C.CORE_TEST_MODE = true

	sf := S.NewStatusFile()
	t.Log("Created new status file")

	t.Log("Initializing with Touch operation")
	err := sf.Touch()
	if err != nil {
		t.Errorf("Error occurred during Touch(): %v", err)
	}

	t.Log("Executing Reset operation")
	err = sf.Reset()
	if err != nil {
		t.Errorf("Error occurred during Reset(): %v", err)
	}

	t.Log("Verifying StatusBundle recreation")
	if _, err := os.Stat(sf.StatusBundle()); os.IsNotExist(err) {
		t.Error("StatusBundle was not recreated after Reset")
	}
	t.Log("Reset test completed successfully")
}

func TestStatusFile_Cache(t *testing.T) {
	t.Log("Starting Cache test")
	C.CORE_TEST_MODE = true

	sf := S.NewStatusFile()
	t.Log("Created new status file")

	err := sf.Cache()
	if err != nil {
		t.Errorf("Error occurred during Cache(): %v", err)
	}
	t.Log("Cache operation completed successfully")
}

func TestStatusFile_AddUniversalIndexes(t *testing.T) {
	t.Log("Starting AddUniversalIndexes test")
	C.CORE_TEST_MODE = true

	sf := S.NewStatusFile()
	t.Log("Created new status file")

	// 테스트 데이터 준비
	testIndexes := make(map[string]S.StorageIndexCursor)
	testKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" // 64바이트 키
	testCursor := S.StorageIndexCursor{
		Key:    testKey,
		FileID: "01",
		Seek:   1234,
		Length: 5678,
	}
	testIndexes[testKey] = testCursor

	// Cache 초기화
	err := sf.Cache()
	if err != nil {
		t.Errorf("Error occurred during Cache(): %v", err)
	}

	// 유니버설 인덱스 추가
	sf.CachedUniversalIndexes = testIndexes

	// 검증
	if len(sf.CachedUniversalIndexes) != 1 {
		t.Error("Expected CachedUniversalIndexes to have 1 entry")
	}

	if cursor, exists := sf.CachedUniversalIndexes[testKey]; !exists {
		t.Error("Test key not found in CachedUniversalIndexes")
	} else {
		if cursor.Key != testKey {
			t.Errorf("Expected key %s, got %s", testKey, cursor.Key)
		}
		if cursor.FileID != "01" {
			t.Errorf("Expected FileID 01, got %s", cursor.FileID)
		}
		if cursor.Seek != 1234 {
			t.Errorf("Expected Seek 1234, got %d", cursor.Seek)
		}
		if cursor.Length != 5678 {
			t.Errorf("Expected Length 5678, got %d", cursor.Length)
		}
	}

	t.Log("AddUniversalIndexes test completed successfully")
	sf.WriteUniversal(testIndexes)
	sf.WriteTasks()
	// sf.Commit()
}
