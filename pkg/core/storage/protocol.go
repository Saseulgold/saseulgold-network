package storage

import (
	"encoding/hex"
	"io/ioutil"

	F "hello/pkg/util"
)

const (
	STATUS_HEAP_BYTES  = 16
	DATA_ID_BYTES      = 2
	STATUS_PREFIX_SIZE = 4

	SEEK_BYTES   = 4
	LENGTH_BYTES = 4

	STATUS_KEY_BYTES = F.STATUS_HASH_BYTES
)

func StatusKey(raw string) string {
	if len(raw) < STATUS_KEY_BYTES {
		return ""
	}
	return hex.EncodeToString([]byte(raw)[:STATUS_KEY_BYTES])
}

type StatusIndex struct {
	FileID string
	Seek   []byte
	Length []byte
	IsSeek int
}

func NewStatusIndex(raw string) StatusIndex {
	offset := STATUS_KEY_BYTES

	fileID := raw[offset : offset+DATA_ID_BYTES]
	offset += DATA_ID_BYTES

	seek := F.Hex2Bin(raw[offset : offset+SEEK_BYTES])
	offset += SEEK_BYTES

	length := F.Hex2Bin(raw[offset : offset+LENGTH_BYTES])

	return StatusIndex{
		FileID: fileID,
		Seek:   seek,
		Length: length,
	}
}

func ReadStatusIndex(indexFile string, bundling bool) map[string]StatusIndex {
	indexes := make(map[string]StatusIndex)

	data, _ := ioutil.ReadFile(indexFile)

	for idx := 0; idx < len(data); idx += STATUS_HEAP_BYTES {
		end := idx + STATUS_HEAP_BYTES
		if end > len(data) {
			end = len(data)
		}
		raw := data[idx:end]

		if len(raw) == STATUS_HEAP_BYTES {
			key := StatusKey(string(raw))
			index := NewStatusIndex(string(raw))

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

func SplitKey(key string) [2]string {
	if len(key) < STATUS_PREFIX_SIZE {
		return [2]string{"", ""}
	}

	prefix := key[:STATUS_PREFIX_SIZE]
	suffix := key[STATUS_PREFIX_SIZE:]

	return [2]string{prefix, suffix}
}
