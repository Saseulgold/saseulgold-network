package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// StatusFile 구조체는 상태 파일 관리를 담당합니다
type StatusFile struct {
	CachedUniversalIndexes map[string][]interface{}
	CachedLocalIndexes     map[string][]interface{}
	Tasks                  [][]interface{}
}

// NewStatusFile creates a new StatusFile instance
func NewStatusFile() *StatusFile {
	return &StatusFile{
		CachedUniversalIndexes: make(map[string][]interface{}),
		CachedLocalIndexes:     make(map[string][]interface{}),
		Tasks:                  make([][]interface{}, 0),
	}
}

// Touch creates necessary directories and files
func (sf *StatusFile) Touch() error {
	if err := os.MkdirAll(config.StatusBundle(), 0755); err != nil {
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
		if err := util.AppendFile(file); err != nil {
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
	if err := os.RemoveAll(config.StatusBundle()); err != nil {
		return err
	}
	return sf.Touch()
}

// Cache loads indexes into memory
func (sf *StatusFile) Cache() error {
	if len(sf.CachedLocalIndexes) == 0 && len(sf.CachedUniversalIndexes) == 0 {
		logger.Log("Bundling: caching..")
		if err := sf.Touch(); err != nil {
			return err
		}
		if err := sf.Commit(); err != nil {
			return err
		}

		var err error
		sf.CachedUniversalIndexes, err = Protocol.ReadStatusIndex(sf.UniversalBundleIndex(), true)
		if err != nil {
			return err
		}

		sf.CachedLocalIndexes, err = Protocol.ReadStatusIndex(sf.LocalBundleIndex(), true)
		if err != nil {
			return err
		}
	}
	return nil
}

// Flush clears cached indexes
func (sf *StatusFile) Flush() {
	sf.CachedUniversalIndexes = make(map[string][]interface{})
	sf.CachedLocalIndexes = make(map[string][]interface{})
}

// 파일 경로 관련 메서드들
func (sf *StatusFile) InfoFile() string {
	return filepath.Join(config.StatusBundle(), "info")
}

func (sf *StatusFile) TempFile() string {
	return filepath.Join(config.StatusBundle(), "tmp")
}

func (sf *StatusFile) LocalFile() string {
	return filepath.Join(config.StatusBundle(), "locals")
}

func (sf *StatusFile) UniversalFile(fileID string) string {
	return filepath.Join(config.StatusBundle(), fmt.Sprintf("universals-%s", fileID))
}

// ... 기타 메서드들은 비슷한 패턴으로 구현
