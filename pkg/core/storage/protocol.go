package storage

import (
	"encoding/hex"
	"fmt"
	C "hello/pkg/core/config"
	F "hello/pkg/util"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func StatusKey(raw string) string {
	if len(raw) < C.STATUS_KEY_BYTES {
		return ""
	}
	return hex.EncodeToString([]byte(raw)[:C.STATUS_KEY_BYTES])
}

type StorageIndexCursor struct {
	Key    string
	FileID string
	Seek   uint64
	Length uint64
	Iseek  uint64
	Value  F.Ia
	Old    F.Ia
	New    F.Ia
}

func StorageKey(raw string) string {
	return raw[:C.STATUS_KEY_BYTES]
}

func NewStorageCursor(raw string) StorageIndexCursor {
	key := StorageKey(raw)
	offset := C.STATUS_KEY_BYTES

	fileID := raw[offset : offset+C.DATA_ID_BYTES]
	offset += C.DATA_ID_BYTES

	seekBytes := raw[offset : offset+C.SEEK_BYTES]
	seek := uint64(seekBytes[0]) | uint64(seekBytes[1])<<8 | uint64(seekBytes[2])<<16 | uint64(seekBytes[3])<<24
	offset += C.SEEK_BYTES

	length := F.Hex2UInt64(raw[offset : offset+C.LENGTH_BYTES])

	return StorageIndexCursor{
		Key:    key,
		FileID: fileID,
		Seek:   seek,
		Length: length,
	}
}

func JoinRootPath(path string) string {
	return filepath.Join(C.QUANTUM_ROOT_DIR, path)
}

func JoinDataRootPath(path string) string {
	return filepath.Join(C.QUANTUM_ROOT_DIR, C.DATA_ROOT_DIR, path)
}

func ReadStatusStorageIndex(indexFile string, bundling bool) map[string]StorageIndexCursor {
	indexes := make(map[string]StorageIndexCursor)
	data, _ := ioutil.ReadFile(indexFile)

	for idx := 0; idx < len(data); idx += C.STATUS_HEAP_BYTES {
		end := idx + C.STATUS_HEAP_BYTES
		if end > len(data) {
			end = len(data)
		}
		raw := data[idx:end]

		if len(raw) == C.STATUS_HEAP_BYTES {
			key := StatusKey(string(raw))
			index := NewStorageCursor(string(raw))

			if bundling {
				iseek := idx * C.STATUS_HEAP_BYTES
				index.Iseek = uint64(iseek)
			}

			indexes[key] = index
		}
	}

	return indexes
}

func KeyBin(key string, keyBytes int) []byte {
	bin, err := hex.DecodeString(key)
	if err != nil {
		return nil
	}

	if len(bin) > keyBytes {
		return bin[:keyBytes]
	}

	result := make([]byte, keyBytes)
	copy(result, bin)
	return result
}

func FileIdBin(fileId string) []byte {
	bin, err := hex.DecodeString(fileId)
	if err != nil {
		return nil
	}

	if len(bin) > C.DATA_ID_BYTES {
		return bin[:C.DATA_ID_BYTES]
	}

	return bin
}

func SplitKey(key string) (string, string) {
	if len(key) < C.STATUS_PREFIX_SIZE {
		return "", ""
	}

	prefix := key[:C.STATUS_PREFIX_SIZE]
	suffix := key[C.STATUS_PREFIX_SIZE:]

	return prefix, suffix
}

func AppendFile(filename string, str string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(str); err != nil {
		return err
	}
	return nil
}

func ListFiles(dirname string, recursive bool) []string {
	info, err := os.Stat(dirname)
	if err != nil || !info.IsDir() {
		return []string{}
	}

	items := []string{}
	contents, err := filepath.Glob(filepath.Join(dirname, "*"))
	if err != nil {
		return []string{}
	}

	for _, item := range contents {
		info, err := os.Stat(item)
		if err != nil {
			continue
		}

		if info.IsDir() && recursive {
			items = append(items, ListFiles(item, recursive)...)
		} else if !info.IsDir() {
			items = append(items, item)
		}
	}

	return items
}

// GrepFiles returns a list of files matching the given prefix
func GrepFiles(dirname string, prefix string) []string {
	files := ListFiles(dirname, true)
	matches := []string{}

	for _, file := range files {
		if strings.HasPrefix(file, prefix) {
			matches = append(matches, file)
		}
	}

	return matches
}

func ChainKey(raw []byte) string {
	return hex.EncodeToString(raw[:C.CHAIN_KEY_BYTES])
}

func ChainIndex(raw []byte) []interface{} {
	offset := C.TIME_HASH_BYTES

	// height 추출
	height := F.BinDec(raw[offset : offset+C.CHAIN_HEIGHT_BYTES])
	offset += C.CHAIN_HEIGHT_BYTES

	// fileID 추출
	fileID := F.Bin2Hex(raw[offset : offset+C.DATA_ID_BYTES])
	offset += C.DATA_ID_BYTES

	// seek 추출
	seek := F.BinDec(raw[offset : offset+C.SEEK_BYTES])
	offset += C.SEEK_BYTES

	// length 추출
	length := F.BinDec(raw[offset : offset+C.LENGTH_BYTES])
	return []interface{}{height, fileID, seek, length}

}

func ReadChainIndexes(directory string) map[string][]interface{} {
	indexes := make(map[string][]interface{})

	indexFile := filepath.Join(directory, "index")
	header, _ := ReadPart(indexFile, 0, C.CHAIN_HEADER_BYTES)
	count := F.BinDec(header)
	length := count * C.CHAIN_HEAP_BYTES

	data, _ := ReadPart(indexFile, int64(C.CHAIN_HEADER_BYTES), length)

	for i := 0; i < len(data); i += C.CHAIN_HEAP_BYTES {
		raw := data[i : i+C.CHAIN_HEAP_BYTES]
		if len(raw) == C.CHAIN_HEAP_BYTES {
			key := ChainKey(raw)
			index := ChainIndex(raw)
			indexes[key] = index
		}
	}

	return indexes
}

func DataID(previousID string) string {
	if previousID != "" {
		// Convert hex string to int
		previous, err := strconv.ParseInt(previousID, 16, 64)
		if err == nil {
			// Increment and convert back to hex string with padding
			return fmt.Sprintf("%04x", previous+1)
		}
	}

	// Return "0000" if no previous ID or error
	return fmt.Sprintf("%04x", 0)
}

func WriteFile(filename string, offset int64, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteAt(data, offset)
	return err
}

func ReadPart(filename string, offset int64, length int) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.Seek(offset, 0); err != nil {
		return nil, err
	}

	buf := make([]byte, length)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf[:n], nil
}
