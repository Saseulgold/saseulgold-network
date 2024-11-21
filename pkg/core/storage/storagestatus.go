package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	F "hello/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// StorageTask represents a single storage operation
type StorageTask struct {
	FilePath string `json:"file_path"`
	Seek     int64  `json:"seek"`
	Data     []byte `json:"data"`
}

type StatusFile struct {
	CachedUniversalIndexes map[string]StorageIndexCursor
	CachedLocalIndexes     map[string]StorageIndexCursor
	Tasks                  []StorageTask
}

var statusFileInstance *StatusFile
var statusFileonce sync.Once

func GetStatusFileInstance() *StatusFile {
	statusFileonce.Do(func() {
		statusFileInstance = &StatusFile{
			CachedUniversalIndexes: make(map[string]StorageIndexCursor),
			CachedLocalIndexes:     make(map[string]StorageIndexCursor),
			Tasks:                  make([]StorageTask, 0),
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
		sf.UniversalBundle("0000"),
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

		// 임시 파일이 비어있는 경우 Commit() 호출을 건너뜁니다
		tmpData, err := os.ReadFile(sf.TempFile())
		if err != nil || len(tmpData) == 0 {
			sf.CachedUniversalIndexes = ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)
			sf.CachedLocalIndexes = ReadStatusStorageIndex(sf.LocalBundleIndex(), true)
			return nil
		}

		if err := sf.Commit(); err != nil {
			return fmt.Errorf("Failed to commit: %v", err)
		}

		sf.CachedUniversalIndexes = ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)
		sf.CachedLocalIndexes = ReadStatusStorageIndex(sf.LocalBundleIndex(), true)
	}
	return nil
}

func (sf *StatusFile) Flush() {
	sf.CachedUniversalIndexes = make(map[string]StorageIndexCursor)
	sf.CachedLocalIndexes = make(map[string]StorageIndexCursor)
}

func (sf *StatusFile) StatusBundle() string {
	return filepath.Join(DataRootDir(), "status_bundle")
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
	res := F.Bin2Hex(F.DecBin(int(F.Hex2UInt64(fileID)+1), C.DATA_ID_BYTES))
	if len(res) != C.DATA_ID_BYTES*2 {
		DebugPanic(fmt.Sprintf("NextFileID length error: %d", len(res)))
	}
	return res
}

func (sf *StatusFile) maxFileId(prefix string) string {
	files, err := filepath.Glob(filepath.Join(sf.StatusBundle(), prefix+"*"))
	if err != nil {
		return "0000"
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

func (sf *StatusFile) WriteLocal(localUpdates UpdateMap) error {
	null := byte(0)
	fileID := "0000" // Default file ID for local bundle
	file := sf.LocalBundle()
	indexFile := sf.LocalBundleIndex()

	if err := AppendFile(file, ""); err != nil {
		return err
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	seek := fileInfo.Size()

	indexInfo, err := os.Stat(indexFile)
	if err != nil {
		return err
	}
	iseek := indexInfo.Size()

	for key, update := range localUpdates {
		key = F.FillHash(key)
		index, exists := sf.CachedLocalIndexes[key]
		data, err := json.Marshal(update.New)
		if err != nil {
			return err
		}
		length := int64(len(data))
		var storedLength int64
		if exists {
			storedLength = index.Length
		}

		var currSeek, currIseek int64

		if storedLength < length {
			// append new line
			currSeek = seek
			seek += length

			if storedLength == 0 {
				// new data
				currIseek = iseek
				iseek += C.STATUS_HEAP_BYTES
			} else {
				// existing data
				currIseek = index.Iseek
			}
		} else {
			// overwrite
			currSeek = index.Seek
			currIseek = index.Iseek
			length = storedLength
			padding := make([]byte, length-int64(len(data)))
			for i := range padding {
				padding[i] = null
			}
			data = append(data, padding...)
		}

		newIndex := NewStorageCursor(key, fileID, currSeek, length)
		newIndex.Iseek = currIseek
		indexData := IndexRaw(key, fileID, currSeek, length)

		sf.CachedLocalIndexes[key] = newIndex
		sf.Tasks = append(sf.Tasks,
			StorageTask{FilePath: sf.LocalBundle(), Seek: currSeek, Data: data},
			StorageTask{FilePath: sf.LocalBundleIndex(), Seek: currIseek, Data: indexData})
	}

	return nil
}

func (sf *StatusFile) WriteUniversal(blockUpdates UpdateMap) error {
	for key, index := range sf.CachedUniversalIndexes {
		DebugLog(fmt.Sprintf("Cached Universal Index - Key: %s, FileID: %s, Seek: %d, Length: %d, Iseek: %d",
			key, index.FileID, index.Seek, index.Length, index.Iseek))
	}

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
	seek := int64(latestFileInfo.Size())

	indexFileInfo, err := os.Stat(indexFile)
	if err != nil {
		return err
	}

	iseek := int64(indexFileInfo.Size())

	for key, update := range blockUpdates {
		key = F.FillHash(key)
		index, exists := sf.CachedUniversalIndexes[key]

		data, err := json.Marshal(update.New)
		DebugLog(fmt.Sprintf("WriteUniversal - Data: %s", data))
		DebugLog(fmt.Sprintf("Index: %v", index))

		if err != nil {
			return err
		}

		length := int64(len(data))
		var oldLength int64
		if exists {
			oldLength = index.Length
		}

		var (
			fileID    string
			currSeek  int64
			currIseek int64
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
				DebugLog(fmt.Sprintf("New data: Key=%s, FileID=%s, Seek=%d, Length=%d, exists=%t\n", key, fileID, currSeek, length, exists))
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
			if int64(len(data)) < length {
				data = append(data, bytes.Repeat([]byte{null}, int(length-int64(len(data))))...)
			}
		}

		newIndex := StorageIndexCursor{
			FileID: fileID,
			Seek:   currSeek,
			Length: length,
			Iseek:  currIseek,
		}

		indexData := IndexRaw(key, fileID, currSeek, length)
		DebugAssert(len(indexData) == C.STATUS_HEAP_BYTES, "invalid index data length: %d, expected: %d", len(indexData), C.STATUS_HEAP_BYTES)

		sf.CachedUniversalIndexes[key] = newIndex
		sf.Tasks = append(sf.Tasks,
			StorageTask{FilePath: sf.UniversalBundle(fileID), Seek: currSeek, Data: data},
			StorageTask{FilePath: sf.UniversalBundleIndex(), Seek: currIseek, Data: indexData})
		fmt.Println("sf.a", currIseek)
	}
	return nil
}

func (sf *StatusFile) WriteTasks() error {
	tasksData, err := json.Marshal(sf.Tasks)
	if err != nil {
		return fmt.Errorf("Failed to serialize tasks: %v", err)
	}

	if err := ioutil.WriteFile(sf.TempFile(), tasksData, 0644); err != nil {
		return fmt.Errorf("Failed to write temporary file: %v", err)
	}

	sf.Tasks = []StorageTask{}
	return nil
}

func (sf *StatusFile) Commit() error {
	// Read tasks from temp file
	raw, err := os.ReadFile(sf.TempFile())
	if err != nil {
		return fmt.Errorf("failed to read temp file: %v", err)
	}

	var tasks []StorageTask
	if err := json.Unmarshal(raw, &tasks); err != nil {
		return fmt.Errorf("failed to unmarshal tasks: %v", err)
	}

	// Process each task
	for _, task := range tasks {
		if task.FilePath == sf.InfoFile() {
			// Overwrite info file
			if err := os.WriteFile(task.FilePath, task.Data, 0644); err != nil {
				return fmt.Errorf("failed to write info file: %v", err)
			}
		} else {
			// Write data at specific position
			f, err := os.OpenFile(task.FilePath, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("failed to open file: %v", err)
			}
			defer f.Close()
			DebugLog(fmt.Sprintf("write data at %s, seek: %d, data length: %d", task.FilePath, task.Seek, len(task.Data)))

			if _, err := f.WriteAt(task.Data, task.Seek); err != nil {
				return fmt.Errorf("failed to write data: %v", err)
			}

			if err := f.Sync(); err != nil {
				return fmt.Errorf("failed to sync file: %v", err)
			}
		}
	}

	return os.WriteFile(sf.TempFile(), []byte{}, 0644)
}

func (sf *StatusFile) GetUniversalIndexes(keys []string) map[string]StorageIndexCursor {
	return GetStatusIndexInstance().UniversalIndexes(keys)
}

func (sf *StatusFile) GetLocalIndexes(keys []string) map[string]StorageIndexCursor {
	return GetStatusIndexInstance().LocalIndexes(keys)
}

func (sf *StatusFile) BundleHeight() int {
	data, err := ioutil.ReadFile(sf.InfoFile())
	if err != nil {
		return 0
	}
	height := 0
	if len(data) > 0 {
		height, _ = strconv.Atoi(string(data))
	}
	return height
}

func (sf *StatusFile) Write(block *Block) error {
	err := sf.Cache()
	if err != nil {
		return err
	}

	err = sf.WriteUniversal(block.UniversalUpdates)
	if err != nil {
		return err
	}
	err = sf.WriteLocal(block.LocalUpdates)
	if err != nil {
		return err
	}

	// make function for making height to byte
	heightData := []byte(strconv.FormatInt(int64(block.Height), 10))
	sf.addTask(sf.InfoFile(), 0, heightData)

	err = sf.WriteTasks()
	if err != nil {
		return err
	}

	return sf.Commit()
}

func (sf *StatusFile) addTask(filePath string, seek int64, data []byte) {
	sf.Tasks = append(sf.Tasks, StorageTask{
		FilePath: filePath,
		Seek:     seek,
		Data:     data,
	})
}

func (sf *StatusFile) UpdateUniversal(indexes map[string]StorageIndexCursor, universalUpdates UpdateMap) map[string]StorageIndexCursor {
	null := []byte{0}
	latestFileID := sf.maxFileId("ubundle-")
	latestFile := sf.UniversalBundle(latestFileID)

	fileInfo, _ := os.Stat(latestFile)
	seek := fileInfo.Size()

	for key, update := range universalUpdates {
		key = F.FillHash(key)
		index, exists := indexes[key]
		data, _ := json.Marshal(update.New)
		length := int64(len(data))
		var storedLength int64
		if exists {
			storedLength = int64(index.Length)
		}

		var fileID string
		if storedLength < length {
			// append new line
			fileID = latestFileID
			if C.LEDGER_FILESIZE_LIMIT < seek+length {
				fileID = sf.NextFileID(fileID)
				seek = 0
			}
			seek += length
		} else {
			// overwrite
			fileID = index.FileID
			seek = int64(index.Seek)
			length = storedLength
			padding := make([]byte, length-int64(len(data)))
			for i := range padding {
				padding[i] = null[0]
			}
			data = append(data, padding...)
		}

		indexes[key] = NewStorageCursor(key, fileID, seek, length)
		AppendFileBytes(sf.UniversalBundle(fileID), data)
	}

	return indexes
}

func (sf *StatusFile) UpdateLocal(indexes map[string]StorageIndexCursor, localUpdates UpdateMap) map[string]StorageIndexCursor {
	null := []byte{0}
	latestFile := sf.LocalBundle()

	fileInfo, _ := os.Stat(latestFile)
	seek := fileInfo.Size()

	for key, update := range localUpdates {
		key = F.FillHash(key)
		index, exists := indexes[key]
		data, _ := json.Marshal(update.New)
		length := int64(len(data))
		var storedLength int64
		if exists {
			storedLength = int64(index.Length)
		}

		var fileID string
		if storedLength < length {
			// append new line
			fileID = "0000"
			if C.LEDGER_FILESIZE_LIMIT < seek+length {
				fileID = sf.NextFileID(fileID)
				seek = 0
			}
			seek += length
		} else {
			// overwrite
			fileID = index.FileID
			seek = int64(index.Seek)
			length = storedLength
			padding := make([]byte, length-int64(len(data)))
			for i := range padding {
				padding[i] = null[0]
			}
			data = append(data, padding...)
		}

		indexes[key] = NewStorageCursor(key, fileID, seek, length)
		AppendFile(sf.LocalBundle(), string(data))
	}

	return indexes
}

func (sf *StatusFile) CopyBundles() error {
	// Copy local bundle
	localBundle := sf.LocalBundle()
	localFile := sf.LocalFile()

	if err := CopyFile(localBundle, localFile); err != nil {
		return err
	}

	// Copy universal bundles
	universalBundles, err := filepath.Glob(filepath.Join(sf.StatusBundle(), "ubundle-*"))
	if err != nil {
		return err
	}

	universalFiles, err := filepath.Glob(filepath.Join(sf.StatusBundle(), "universals-*"))
	if err != nil {
		return err
	}

	// Delete existing universal files
	for _, file := range universalFiles {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	// Copy universal bundles to new files
	for _, bundle := range universalBundles {
		from := bundle
		to := strings.Replace(bundle, "ubundle-", "universals-", 1)

		if err := CopyFile(from, to); err != nil {
			return err
		}
	}

	return nil
}

func (sf *StatusFile) Update(block *Block) error {
	localUpdates := block.LocalUpdates
	universalUpdates := block.UniversalUpdates

	// Handle local updates
	localKeys := make([]string, 0, len(localUpdates))
	for k := range localUpdates {
		localKeys = append(localKeys, k)
	}
	localIndexes := sf.GetLocalIndexes(localKeys)
	sf.UpdateLocal(localIndexes, localUpdates)

	// Handle universal updates
	universalKeys := make([]string, 0, len(universalUpdates))
	for k := range universalUpdates {
		universalKeys = append(universalKeys, k)
	}
	universalIndexes := sf.GetUniversalIndexes(universalKeys)
	sf.UpdateUniversal(universalIndexes, universalUpdates)

	return nil
}

func (sf *StatusFile) CountLocalStatus(prefix string) int {
	return GetStatusIndexInstance().CountLocalIndexes(prefix)
}

func (sf *StatusFile) CountUniversalStatus(prefix string) int {
	return GetStatusIndexInstance().CountUniversalIndexes(prefix)
}
