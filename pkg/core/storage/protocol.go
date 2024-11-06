package storage

import (
	"encoding/hex"
	F "hello/pkg/util"
	"io/ioutil"
	"os"
)

const (
	STATUS_HEAP_BYTES  = 16
	DATA_ID_BYTES      = 2
	STATUS_PREFIX_SIZE = 4

	SEEK_BYTES   = 4
	LENGTH_BYTES = 4

	STATUS_KEY_BYTES = F.STATUS_HASH_BYTES
)

// var DATA_ROOT_DIR = os.Getenv("QUANTUM_DATA_DIR")
const DATA_ROOT_DIR = "/Users/louis/qcn/data"
const DATA_ROOT_TEST_DIR = "/Users/louis/qcn/data/test"

func StatusKey(raw string) string {
	if len(raw) < STATUS_KEY_BYTES {
		return ""
	}
	return hex.EncodeToString([]byte(raw)[:STATUS_KEY_BYTES])
}

type StorageIndexCursor struct {
	Key    string
	FileID string
	Seek   []byte
	Length []byte
	IsSeek int
	Value  F.Ia
}

func StorageKey(raw string) string {
	return raw[:STATUS_KEY_BYTES]
}

func NewStorageCursor(raw string) StorageIndexCursor {
	key := StorageKey(raw)
	offset := STATUS_KEY_BYTES

	fileID := raw[offset : offset+DATA_ID_BYTES]
	offset += DATA_ID_BYTES

	seek := F.Hex2Bin(raw[offset : offset+SEEK_BYTES])
	offset += SEEK_BYTES

	length := F.Hex2Bin(raw[offset : offset+LENGTH_BYTES])

	return StorageIndexCursor{
		Key:    key,
		FileID: fileID,
		Seek:   seek,
		Length: length,
	}
}

func ReadStorageIndex(indexFile string, bundling bool) map[string]StorageIndexCursor {
	indexes := make(map[string]StorageIndexCursor)
	data, _ := ioutil.ReadFile(indexFile)

	for idx := 0; idx < len(data); idx += STATUS_HEAP_BYTES {
		end := idx + STATUS_HEAP_BYTES
		if end > len(data) {
			end = len(data)
		}
		raw := data[idx:end]

		if len(raw) == STATUS_HEAP_BYTES {
			key := StatusKey(string(raw))
			index := NewStorageCursor(string(raw))

			if bundling {
				iseek := idx * STATUS_HEAP_BYTES
				index.IsSeek = iseek
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

	if len(bin) > DATA_ID_BYTES {
		return bin[:DATA_ID_BYTES]
	}

	return bin
}

func SplitKey(key string) (string, string) {
	if len(key) < STATUS_PREFIX_SIZE {
		return "", ""
	}

	prefix := key[:STATUS_PREFIX_SIZE]
	suffix := key[STATUS_PREFIX_SIZE:]

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
