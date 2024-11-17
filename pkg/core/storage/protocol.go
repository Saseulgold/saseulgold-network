package storage

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	F "hello/pkg/util"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func StatusKey(raw []byte) string {
	res := F.Bin2Hex(raw[:C.STATUS_KEY_BYTES])
	DebugLog(fmt.Sprintf("StatusKey: %s", res))
	return res
}

type StorageIndexCursor struct {
	Key    string
	FileID string
	Seek   int64
	Length int64
	Iseek  int64
	Value  F.Ia
	Old    F.Ia
	New    F.Ia
}

func StorageKey(raw string) string {
	return hex.EncodeToString([]byte(raw[:C.STATUS_KEY_BYTES]))
}

func NewStorageCursor(key string, fileID string, seek int64, length int64) StorageIndexCursor {
	return StorageIndexCursor{
		Key:    key,
		FileID: fileID,
		Seek:   seek,
		Length: length,
	}
}

/*
*

	func NewStorageCursorRaw(raw string) StorageIndexCursor {
		key := StorageKey(raw)
		offset := C.STATUS_KEY_BYTES

		fileID := raw[offset : offset+C.DATA_ID_BYTES]
		offset += C.DATA_ID_BYTES

		seekBytes := raw[offset : offset+C.SEEK_BYTES]
		seek := int64(seekBytes[0]) | int64(seekBytes[1])<<8 | int64(seekBytes[2])<<16 | int64(seekBytes[3])<<24
		offset += C.SEEK_BYTES

		length := F.Hex2Int64(raw[offset : offset+C.LENGTH_BYTES])

		return StorageIndexCursor{
			Key:    key,
			FileID: fileID,
			Seek:   seek,
			Length: length,
		}
	}

	*
*/
func NewStorageCursorRaw(raw string) StorageIndexCursor {
	key, fileID, seek, length, err := ParseIndexRaw([]byte(raw))
	if err != nil {
		DebugPanic(fmt.Sprintf("NewStorageCursorRaw error: %v", err))
	}

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
	data, _ := os.ReadFile(indexFile)

	for idx := 0; idx < len(data); idx += C.STATUS_HEAP_BYTES {
		end := idx + C.STATUS_HEAP_BYTES
		if end > len(data) {
			end = len(data)
		}
		raw := data[idx:end]

		if len(raw) == C.STATUS_HEAP_BYTES {
			key := StatusKey(raw)
			DebugLog(fmt.Sprintf("raw : %s", raw))
			DebugLog(fmt.Sprintf("key : %s", key))
			index := NewStorageCursorRaw(string(raw))

			if bundling {
				iseek := idx * C.STATUS_HEAP_BYTES
				index.Iseek = int64(iseek)
			}

			indexes[key] = index
		}
	}

	return indexes
}

func KeyBin(key string, keyBytes int) []byte {
	if key == "" {
		result := make([]byte, keyBytes)
		return result
	}

	bin, err := hex.DecodeString(key)
	if err != nil {
		DebugLog(fmt.Sprintf("KeyBin Error: %v", err))
		result := make([]byte, keyBytes)
		return result
	}

	if len(bin) > keyBytes {
		return bin[:keyBytes]
	}

	result := make([]byte, keyBytes)
	copy(result, bin)
	return result
}

func FileIdBin(fileId string) []byte {
	bin := F.Hex2Bin(fileId)

	if len(bin) > C.DATA_ID_BYTES {
		return bin[:C.DATA_ID_BYTES]
	}

	if len(bin) != C.DATA_ID_BYTES {
		DebugPanic(fmt.Sprintf("FileIdBin length error: %d", len(bin)))
	}

	return bin
}

// BinToFileId는 바이트 배열을 파일 ID 문자열로 변환합니다
func BinToFileId(bin []byte) string {
	if len(bin) > C.DATA_ID_BYTES {
		bin = bin[:C.DATA_ID_BYTES]
	}

	// 16진수 문자열로 변환하고 왼쪽에 0을 채워 4자리로 맞춤
	return fmt.Sprintf("%04s", F.Bin2Hex(bin))
}

func SplitKey(key string) (string, string) {
	if len(key) < C.STATUS_PREFIX_SIZE {
		return "", ""
	}

	prefix := key[:C.STATUS_PREFIX_SIZE]
	suffix := key[C.STATUS_PREFIX_SIZE:]

	return prefix, suffix
}

func AppendFileBytes(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
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

func CopyFile(from string, to string) error {
	source, err := os.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(to)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func IndexRaw(key string, fileID string, seek int64, length int64) []byte {
	DebugLog(fmt.Sprintf("indexRaw - Key: %s, FileID: %s, Seek: %d, Length: %d", key, fileID, seek, length))
	var result []byte
	result = append(result, KeyBin(key, C.STATUS_KEY_BYTES)...)
	result = append(result, FileIdBin(fileID)...)

	seekBytes := make([]byte, C.SEEK_BYTES)
	binary.LittleEndian.PutUint32(seekBytes, uint32(seek))
	result = append(result, seekBytes...)

	lengthBytes := make([]byte, C.LENGTH_BYTES)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(length))
	result = append(result, lengthBytes...)

	if len(result) != C.STATUS_HEAP_BYTES {
		DebugPanic(fmt.Sprintf("IndexRaw length error: %d", len(result)))
	}

	return result
}

func ParseIndexRaw(data []byte) (key string, fileID string, seek int64, length int64, err error) {
	if len(data) != C.STATUS_HEAP_BYTES {
		return "", "", 0, 0, fmt.Errorf("invalid index data length: %d, expected: %d", len(data), C.STATUS_HEAP_BYTES)
	}

	keyBytes := data[:C.STATUS_KEY_BYTES]
	key = F.Bin2Hex(keyBytes)

	offset := C.STATUS_KEY_BYTES
	fileIDBytes := data[offset : offset+C.DATA_ID_BYTES]
	fileID = BinToFileId(fileIDBytes)

	offset += C.DATA_ID_BYTES
	seekBytes := data[offset : offset+C.SEEK_BYTES]
	seek = int64(binary.LittleEndian.Uint32(seekBytes))

	offset += C.SEEK_BYTES
	lengthBytes := data[offset : offset+C.LENGTH_BYTES]
	length = int64(binary.LittleEndian.Uint32(lengthBytes))

	return key, fileID, seek, length, nil
}
