package storage

import (
	"bytes"
	"encoding/json"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	F "hello/pkg/util"
	"io"
	"math"
	"os"
	"path/filepath"
)

type ChainStorage struct{}

func (c *ChainStorage) Block(directory string, needle interface{}) (*Block, error) {
	data, err := c.ReadData(directory, needle.([]interface{}))
	if err != nil {
		return nil, err
	}

	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}

	return &block, nil
}

func (c *ChainStorage) IndexFile(directory string) string {
	path := C.DATA_ROOT_DIR
	if C.IS_TEST {
		path = C.DATA_ROOT_TEST_DIR
	}
	println("Chain index file path:", filepath.Join(path, directory, "index"))
	return filepath.Join(path, directory, "index")
}

func (c *ChainStorage) DataFile(directory, fileID string) string {
	if C.IS_TEST {
		return filepath.Join(C.DATA_ROOT_TEST_DIR, directory, fileID)
	}
	return filepath.Join(C.DATA_ROOT_DIR, directory, fileID)
}

func (c *ChainStorage) Touch(directory string) error {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	_, err := os.OpenFile(c.IndexFile(directory), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
}

func (c *ChainStorage) ResetData(directory string) error {
	if err := os.RemoveAll(directory); err != nil {
		return err
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	_, err := os.OpenFile(c.IndexFile(directory), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
}

func (c *ChainStorage) Index(directory string, needle interface{}) ([]interface{}, error) {
	if height, ok := needle.(int); ok {
		idx := c.ReadIdx(directory, height)
		return c.ReadIndex(directory, idx)
	}
	return c.SearchIndex(directory, needle.(string))
}

func (c *ChainStorage) ReadIdx(directory string, height int) int {
	idx := 0
	lastIdx := c.LastIdx(directory)
	println("마지막 인덱스:", lastIdx)
	lastIndex, _ := c.ReadIndex(directory, lastIdx)
	lastHeight := 0
	if len(lastIndex) > 0 {
		lastHeight = lastIndex[0].(int)
	}
	gap := lastHeight - height

	if gap >= 0 {
		idx = lastIdx - gap
	}

	return idx
}

func (c *ChainStorage) ReadIndex(directory string, idx int) ([]interface{}, error) {
	if idx > 0 {
		iseek := C.CHAIN_HEADER_BYTES + (idx-1)*C.CHAIN_HEAP_BYTES
		raw, err := ReadPart(c.IndexFile(directory), int64(iseek), C.CHAIN_HEAP_BYTES)

		if err != nil {
			return nil, err
		}

		if len(raw) == C.CHAIN_HEAP_BYTES {
			return ChainIndex(raw), nil
		}
	}
	return []interface{}{}, nil
}

func (c *ChainStorage) ReadData(directory string, index []interface{}) ([]byte, error) {
	fileID := index[1].(string)
	seek := index[2].(int)
	length := index[3].(int)

	return ReadPart(c.DataFile(directory, fileID), int64(seek), length)
}

func (c *ChainStorage) LastIdx(directory string) int {
	header, _ := ReadPart(c.IndexFile(directory), 0, C.CHAIN_HEADER_BYTES)
	return F.BinDec(header)
}

func (c *ChainStorage) LastIndex(directory string) ([]interface{}, error) {
	lastIdx := c.LastIdx(directory)
	return c.ReadIndex(directory, lastIdx)
}

func (c *ChainStorage) WriteData(directory string, height int, key string, data []byte) error {
	lastIdx := c.LastIdx(directory)
	lastIndex, _ := c.ReadIndex(directory, lastIdx)

	var lastHeight int
	if len(lastIndex) > 0 {
		lastHeight = lastIndex[0].(int)
	}

	if height != lastHeight+1 {
		return nil
	}

	var fileID string
	var lastSeek, lastLength int

	if len(lastIndex) > 0 {
		fileID = lastIndex[1].(string)
		lastSeek = lastIndex[2].(int)
		lastLength = lastIndex[3].(int)
	} else {
		fileID = DataID("")
	}

	seek := lastSeek + lastLength
	length := len(data)

	idx := lastIdx + 1
	iseek := C.CHAIN_HEADER_BYTES + lastIdx*C.CHAIN_HEAP_BYTES

	if err := AppendFile(c.DataFile(directory, fileID), ""); err != nil {
		return err
	}

	if C.LEDGER_FILESIZE_LIMIT < seek+length {
		fileID = DataID(fileID)
		seek = 0
	}

	headerData := c.headerRaw(idx)
	indexData := c.indexRaw(key, fileID, height, seek, length)

	if err := WriteFile(c.DataFile(directory, fileID), int64(seek), data); err != nil {
		return err
	}
	if err := WriteFile(c.IndexFile(directory), int64(iseek), indexData); err != nil {
		return err
	}
	return WriteFile(c.IndexFile(directory), 0, headerData)
}

func (c *ChainStorage) RemoveData(directory string, idx int) error {
	if idx <= 0 {
		idx = 1
	}
	headerData := c.headerRaw(idx - 1)
	return WriteFile(c.IndexFile(directory), 0, headerData)
}

func (c *ChainStorage) headerRaw(idx int) []byte {
	return F.DecBin(idx, C.CHAIN_HEADER_BYTES)
}

func (c *ChainStorage) indexRaw(key, fileID string, height, seek, length int) []byte {
	return append(append(append(append(
		KeyBin(key, C.CHAIN_KEY_BYTES),
		F.DecBin(height, C.CHAIN_HEIGHT_BYTES)...),
		FileIdBin(fileID)...),
		F.DecBin(seek, C.SEEK_BYTES)...),
		F.DecBin(length, C.LENGTH_BYTES)...)
}

func (c *ChainStorage) SearchIndex(directory string, hash string) ([]interface{}, error) {
	if len(hash) < C.HEX_TIME_SIZE {
		return []interface{}{}, nil
	}

	target := F.Hex2Bin(hash[:C.HEX_TIME_SIZE])
	header, err := ReadPart(c.IndexFile(directory), 0, C.CHAIN_HEADER_BYTES)
	if err != nil {
		return nil, err
	}
	count := F.BinDec(header)
	cycle := int(math.Log2(float64(count))) + 1

	min := 0
	max := count - 1

	a := min
	b := max + 1
	mid := 0

	bytes := len(target)
	f, err := os.Open(c.IndexFile(directory))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	for i := 0; i < cycle; i++ {
		mid = (a + b) / 2

		var l []byte
		if mid <= 0 {
			l = make([]byte, bytes)
		} else {
			if _, err := f.Seek(int64(C.CHAIN_HEADER_BYTES+(mid-1)*C.CHAIN_HEAP_BYTES), 0); err != nil {
				return nil, err
			}
			l = make([]byte, bytes)
			if _, err := f.Read(l); err != nil {
				return nil, err
			}
		}

		var r []byte
		if mid >= max {
			r = bytes_repeat(0xff, bytes)
		} else {
			if _, err := f.Seek(int64(C.CHAIN_HEADER_BYTES+mid*C.CHAIN_HEAP_BYTES), 0); err != nil {
				return nil, err
			}
			r = make([]byte, bytes)
			if _, err := f.Read(r); err != nil {
				return nil, err
			}
		}

		if mid == min || mid == max || (bytes_less(l, target) && bytes_less_equal(target, r)) {
			break
		} else if bytes_less_equal(target, l) {
			b = mid
		} else {
			a = mid
		}
	}

	if _, err := f.Seek(int64(C.CHAIN_HEADER_BYTES+mid*C.CHAIN_HEAP_BYTES), 0); err != nil {
		return nil, err
	}
	read := make([]byte, C.CHAIN_HEAP_BYTES)
	n, err := f.Read(read)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if n == C.CHAIN_HEAP_BYTES {
		return ChainIndex(read), nil
	}

	return []interface{}{}, nil
}

func bytes_repeat(b byte, count int) []byte {
	bytes := make([]byte, count)
	for i := 0; i < count; i++ {
		bytes[i] = b
	}
	return bytes
}

func bytes_less(a, b []byte) bool {
	return bytes.Compare(a, b) < 0
}

func bytes_less_equal(a, b []byte) bool {
	return bytes.Compare(a, b) <= 0
}
