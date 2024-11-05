package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// StatusFile 구조체는 상태 파일 관리를 담당합니다
type StatusFile struct {
	CachedUniversalIndexes map[string]StorageIndex
	CachedLocalIndexes     map[string]StorageIndex
	Tasks                  [][]interface{}
}

// NewStatusFile creates a new StatusFile instance
func NewStatusFile() *StatusFile {
	return &StatusFile{
		CachedUniversalIndexes: make(map[string]StorageIndex),
		CachedLocalIndexes:     make(map[string]StorageIndex),
		Tasks:                  make([][]interface{}, 0),
	}
}

// Touch creates necessary directories and files
func (sf *StatusFile) Touch() error {
	if err := os.MkdirAll(sf.StatusBundle(), 0755); err != nil {
		return err
	}

	files := []string{
		sf.TempFile(),
		sf.InfoFile(),
		sf.LocalFile(),
		sf.LocalBundle(),
		sf.LocalBundleIndex(),
		sf.UniversalBundleIndex(),
	}

	for _, file := range files {
		if err := AppendFile(file, ""); err != nil {
			return err
		}
	}
	return nil
}

// Init initializes the status file
func (sf *StatusFile) Init() error {
	return sf.Touch()
}

// Reset deletes and recreates the status bundle
func (sf *StatusFile) Reset() error {
	if err := os.RemoveAll(sf.StatusBundle()); err != nil {
		return err
	}
	return sf.Touch()
}

// Cache loads indexes into memory
func (sf *StatusFile) Cache() error {
	if len(sf.CachedLocalIndexes) == 0 && len(sf.CachedUniversalIndexes) == 0 {
		if err := sf.Touch(); err != nil {
			return err
		}

		/**
		if err := sf.Commit(); err != nil {
			return err
		}
		**/

		var err error

		sf.CachedUniversalIndexes = ReadStorageIndex(sf.UniversalBundleIndex(), true)

		if err != nil {
			return err
		}

		sf.CachedLocalIndexes = ReadStorageIndex(sf.LocalBundleIndex(), true)
		if err != nil {
			return err
		}
	}
	return nil
}

// Flush clears cached indexes
func (sf *StatusFile) Flush() {
	sf.CachedUniversalIndexes = make(map[string]StorageIndex)
	sf.CachedLocalIndexes = make(map[string]StorageIndex)
}

func (sf *StatusFile) StatusBundle() string {
	return filepath.Join(DATA_ROOT_DIR, "statusbundle")
}

func (sf *StatusFile) LocalBundle() string {
	return filepath.Join(DATA_ROOT_DIR, "localbundle")
}

func (sf *StatusFile) LocalBundleIndex() string {
	return filepath.Join(DATA_ROOT_DIR, "localbundleindex")
}

func (sf *StatusFile) UniversalBundleIndex() string {
	return filepath.Join(DATA_ROOT_DIR, "universalsbundleindex")
}

// 파일 경로 관련 메서드들
func (sf *StatusFile) InfoFile() string {
	return filepath.Join(sf.StatusBundle(), "info")
}

func (sf *StatusFile) TempFile() string {
	return filepath.Join(sf.StatusBundle(), "tmp")
}

func (sf *StatusFile) LocalFile() string {
	return filepath.Join(sf.StatusBundle(), "locals")
}

func (sf *StatusFile) UniversalFile(fileID string) string {
	return filepath.Join(sf.StatusBundle(), fmt.Sprintf("universals-%s", fileID))
}

// Commit 메서드 추가

/**
func (sf *StatusFile) Commit() error {
	if err := WriteStatusIndex(sf.UniversalBundleIndex(), sf.CachedUniversalIndexes); err != nil {
		return err
	}
	if err := WriteStatusIndex(sf.LocalBundleIndex(), sf.CachedLocalIndexes); err != nil {
		return err
	}
	return nil
}
	**/

// ... 기타 메서드들은 비슷한 패턴으로 구현
