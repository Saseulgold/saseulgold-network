package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "hello/pkg/core/config"
	C "hello/pkg/core/config"
	F "hello/pkg/util"
	"os"
	"path/filepath"
	"strings"
)

// StatusFile 구조체는 상태 파일 관리를 담당합니다
type StatusFile struct {
	CachedUniversalIndexes map[string]StorageIndexCursor
	CachedLocalIndexes     map[string]StorageIndexCursor
	Tasks                  [][]interface{}
}

// NewStatusFile creates a new StatusFile instance
func NewStatusFile() *StatusFile {
	return &StatusFile{
		CachedUniversalIndexes: make(map[string]StorageIndexCursor),
		CachedLocalIndexes:     make(map[string]StorageIndexCursor),
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

func (sf *StatusFile) Flush() {
	sf.CachedUniversalIndexes = make(map[string]StorageIndexCursor)
	sf.CachedLocalIndexes = make(map[string]StorageIndexCursor)
}

func (sf *StatusFile) RootDir() string {
	if C.CORE_TEST_MODE {
		return DATA_ROOT_TEST_DIR
	}
	return DATA_ROOT_DIR
}

func (sf *StatusFile) StatusBundle() string {
	return filepath.Join(DATA_ROOT_DIR, "statusbundle")
}

func (sf *StatusFile) LocalBundle() string {
	return filepath.Join(DATA_ROOT_DIR, "lbundle")
}

func (sf *StatusFile) LocalBundleIndex() string {
	return filepath.Join(DATA_ROOT_DIR, "lbundle_index")
}

func (sf *StatusFile) UniversalBundleIndex() string {
	return filepath.Join(DATA_ROOT_DIR, "ubundle_index")
}

func (sf *StatusFile) UniversalBundle(fileID string) string {
	return filepath.Join(DATA_ROOT_DIR, fmt.Sprintf("ubundle-%s", fileID))
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

func (sf *StatusFile) NextFileID(fileID string) string {
	return string(F.Hex2Int64(fileID) + 1)
}

func (sf *StatusFile) maxFileId(prefix string) string {
	files, err := filepath.Glob(filepath.Join(sf.StatusBundle(), prefix+"*"))
	if err != nil {
		return "00"
	}

	if len(files) > 0 {
		fileIds := make([]string, len(files))
		for i, file := range files {
			fileIds[i] = strings.TrimPrefix(filepath.Base(file), prefix)
		}

		maxId := fileIds[0]
		for _, id := range fileIds[1:] {
			if id > maxId {
				maxId = id
			}
		}
		return maxId
	}

	return fmt.Sprintf("%02x", 0)
}

func (sf *StatusFile) WriteUniversal(universalUpdates map[string]map[string]interface{}) error {
	null := byte(0)
	latestFileID := sf.maxFileId("ubundle-")
	latestFile := sf.UniversalBundle(latestFileID)
	indexFile := sf.UniversalBundleIndex()

	if err := AppendFile(latestFile, ""); err != nil {
		return err
	}

	latestFileInfo, err := os.Stat(latestFile)
	if err != nil {
		return err
	}
	seek := uint64(latestFileInfo.Size())

	indexFileInfo, err := os.Stat(indexFile)
	if err != nil {
		return err
	}

	iseek := uint64(indexFileInfo.Size())

	for key, update := range universalUpdates {
		key = F.FillHash(key)
		index, exists := sf.CachedUniversalIndexes[key]

		data, err := json.Marshal(update["new"])
		if err != nil {
			return err
		}

		length := uint64(len(data))
		var oldLength uint64
		if exists {
			oldLength = index.Length
		}

		var (
			fileID    string
			currSeek  uint64
			currIseek uint64
		)

		if oldLength < length {
			// 새로운 위치에 데이터 추가
			fileID = latestFileID
			currSeek = seek
			seek += length

			if C.LEDGER_FILESIZE_LIMIT < currSeek+length {
				fileID = sf.NextFileID(fileID)
				currSeek = 0
				seek = length
			}

			if oldLength == 0 {
				// 새로운 데이터
				currIseek = iseek
				iseek += C.STATUS_HEAP_BYTES
			} else {
				// 기존 데이터 업데이트
				currIseek = index.Iseek
			}
		} else {
			// 기존 위치에 덮어쓰기
			fileID = index.FileID
			currSeek = index.Seek
			currIseek = index.Iseek
			length = oldLength
			// 데이터 패딩
			if uint64(len(data)) < length {
				data = append(data, bytes.Repeat([]byte{null}, int(length-uint64(len(data))))...)
			}
		}

		/**
		// 인덱스 업데이트
		newIndex := StorageIndexCursor{
			FileID: fileID,
			Seek:   currSeek,
			Length: length,
			Iseek:  currIseek,
		}

		indexData := sf.indexRaw(key, fileID, currSeek, length)

		sf.CachedUniversalIndexes[key] = newIndex
		sf.Tasks = append(sf.Tasks,
			[]interface{}{sf.UniversalBundle(fileID), currSeek, data},
			[]interface{}{sf.UniversalBundleIndex(fileID), currIseek, indexData})
		**/
		fmt.Println("sf.a", currIseek)
	}
	return nil
}
