package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "hello/pkg/core/config"
	C "hello/pkg/core/config"
	F "hello/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type StatusFile struct {
	CachedUniversalIndexes map[string]StorageIndexCursor
	CachedLocalIndexes     map[string]StorageIndexCursor
	Tasks                  [][]interface{}
}

var statusFileInstance *StatusFile
var statusFileonce sync.Once

func GetStatusFileInstance() *StatusFile {
	statusFileonce.Do(func() {
		statusFileInstance = &StatusFile{
			CachedUniversalIndexes: make(map[string]StorageIndexCursor),
			CachedLocalIndexes:     make(map[string]StorageIndexCursor),
			Tasks:                  make([][]interface{}, 0),
		}
	})
	return statusFileInstance
}

// Touch creates necessary directories and files
func (sf *StatusFile) Touch() error {
	bundlePath := sf.StatusBundle()
	if err := os.MkdirAll(bundlePath, 0755); err != nil {
		return err
	}

	fmt.Println("sf.StatusBundle(): ", sf.StatusBundle())

	files := []string{
		sf.TempFile(),
		sf.InfoFile(),
		sf.LocalFile(),
		sf.LocalBundle(),
		sf.UniversalBundle("00"),
		sf.LocalBundleIndex(),
		sf.UniversalBundleIndex(),
	}

	for _, file := range files {
		fmt.Println("f: ", file)
		f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		f.Close()
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

		/** TODO
		if err := sf.Commit(); err != nil {
			return err
		}
		*/

		var err error

		sf.CachedUniversalIndexes = ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)

		if err != nil {
			return err
		}

		sf.CachedLocalIndexes = ReadStatusStorageIndex(sf.LocalBundleIndex(), true)
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

func (sf *StatusFile) DataRootDir() string {
	if C.CORE_TEST_MODE {
		return DATA_ROOT_TEST_DIR
	}
	return DATA_ROOT_DIR
}

func (sf *StatusFile) StatusBundle() string {
	return filepath.Join(sf.DataRootDir(), "statusbundle")
}

func (sf *StatusFile) LocalBundle() string {
	return filepath.Join(sf.StatusBundle(), "lbundle")
}

func (sf *StatusFile) LocalBundleIndex() string {
	return filepath.Join(sf.StatusBundle(), "lbundle_index")
}

func (sf *StatusFile) UniversalBundleIndex() string {
	return filepath.Join(sf.StatusBundle(), "ubundle_index")
}

func (sf *StatusFile) UniversalBundle(fileID string) string {
	return filepath.Join(sf.StatusBundle(), fmt.Sprintf("ubundle-%s", fileID))
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

func (sf *StatusFile) indexRaw(key string, fileID string, seek uint64, length uint64) string {
	keyStr := key
	fileIdStr := fileID
	seekStr := fmt.Sprintf("%d", seek)
	lengthStr := fmt.Sprintf("%d", length)

	result := keyStr + fileIdStr + seekStr + lengthStr
	return result
}

func (sf *StatusFile) WriteUniversal(universalUpdates map[string]StorageIndexCursor) error {
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

		data, err := json.Marshal(update.New)
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
			// Add data to new location
			fileID = latestFileID
			currSeek = seek
			seek += length

			if C.LEDGER_FILESIZE_LIMIT < currSeek+length {
				fileID = sf.NextFileID(fileID)
				currSeek = 0
				seek = length
			}

			if oldLength == 0 {
				// New data
				currIseek = iseek
				iseek += C.STATUS_HEAP_BYTES
			} else {
				// Update existing data
				currIseek = index.Iseek
			}
		} else {
			// Overwrite existing location
			fileID = index.FileID
			currSeek = index.Seek
			currIseek = index.Iseek
			length = oldLength
			// Pad data
			if uint64(len(data)) < length {
				data = append(data, bytes.Repeat([]byte{null}, int(length-uint64(len(data))))...)
			}
		}

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
			[]interface{}{sf.UniversalBundleIndex(), currIseek, indexData})
		fmt.Println("sf.a", currIseek)
	}
	return nil
}

func (sf *StatusFile) WriteTasks() error {
	// Serialize tasks to JSON
	tasksData, err := json.Marshal(sf.Tasks)
	if err != nil {
		return fmt.Errorf("Failed to serialize tasks: %v", err)
	}

	// Write to temporary file
	if err := ioutil.WriteFile(sf.TempFile(), tasksData, 0644); err != nil {
		return fmt.Errorf("Failed to write temporary file: %v", err)
	}

	// Reset tasks
	sf.Tasks = [][]interface{}{}
	return nil
}

func (sf *StatusFile) Commit() error {
	// Read tasks from temp file
	raw, err := ioutil.ReadFile(sf.TempFile())
	if err != nil {
		return fmt.Errorf("Failed to read temp file: %v", err)
	}

	var tasks [][]interface{}
	if err := json.Unmarshal(raw, &tasks); err != nil {
		return fmt.Errorf("Failed to unmarshal tasks: %v", err)
	}

	// Process each task
	for _, item := range tasks {
		file := item[0].(string)
		seek := int64(item[1].(float64))

		var data []byte
		switch v := item[2].(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return fmt.Errorf("unexpected data type: %T", item[2])
		}

		if file == sf.InfoFile() {
			// Overwrite info file
			if err := ioutil.WriteFile(file, data, 0644); err != nil {
				return fmt.Errorf("Failed to write info file: %v", err)
			}
		} else {
			// Write data at specific position
			f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("Failed to open file: %v", err)
			}
			defer f.Close()

			if _, err := f.WriteAt(data, seek); err != nil {
				return fmt.Errorf("Failed to write data: %v", err)
			}

			if err := f.Sync(); err != nil {
				return fmt.Errorf("Failed to sync file: %v", err)
			}
		}
	}

	// Clear temp file
	return ioutil.WriteFile(sf.TempFile(), []byte{}, 0644)
}

func (sf *StatusFile) GetUniversalIndexes(keys []string) map[string]StorageIndexCursor {

	return nil
}

func (sf *StatusFile) GetLocalIndexes(keys []string) map[string]StorageIndexCursor {
	return nil
}

func (sf *StatusFile) GetLocalStatus(key string) interface{} {
	return 1
}

func (sf *StatusFile) GetUniversalStatus(key string) interface{} {
	return 1
}
