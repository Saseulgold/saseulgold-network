package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	F "hello/pkg/util"
	"io"
	"math"
	"os"
	"path/filepath"
)

type ChainStorage struct{}

var chainStorageInstance *ChainStorage

func GetChainStorageInstance() *ChainStorage {
	if chainStorageInstance == nil {
		chainStorageInstance = &ChainStorage{}
	}
	return chainStorageInstance
}

func ChainInfo() string {
	return filepath.Join(DataRootDir(), "chain_info")
}

func MainChain() string {
	return filepath.Join(DataRootDir(), "main_chain")
}

func LastHeight() int {
	data, _ := os.ReadFile(ChainInfo())
	height := 0
	if len(data) > 0 {
		height = int(F.BinDec(data))
	}
	return height
}

func (c *ChainStorage) LastBlock() (*Block, error) {
	return c.GetBlock(LastHeight())
}

func ConfirmedHeight() int {
	return LastHeight() - 5
}

func ParseBlock(data []byte) (*Block, error) {
	var block Block
	var rawData map[string]interface{}

	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}

	if universalUpdatesRaw, exists := rawData["universal_updates"].(map[string]interface{}); exists {
		block.UniversalUpdates = make(UpdateMap)
		for key, value := range universalUpdatesRaw {
			updateBytes, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			var update Update
			if err := json.Unmarshal(updateBytes, &update); err != nil {
				return nil, err
			}
			block.UniversalUpdates[key] = update
		}
	}

	if localUpdatesRaw, exists := rawData["local_updates"].(map[string]interface{}); exists {
		block.LocalUpdates = make(UpdateMap)
		for key, value := range localUpdatesRaw {
			updateBytes, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			var update Update
			if err := json.Unmarshal(updateBytes, &update); err != nil {
				return nil, err
			}
			block.LocalUpdates[key] = update
		}
	}

	block.Init()

	return &block, nil
}

func (c *ChainStorage) Block(needle interface{}) (*Block, error) {
	data, err := c.ReadData(needle.([]interface{}))
	if err != nil {
		return nil, err
	}

	block, err := ParseBlock(data)

	if err != nil {
		return nil, err
	}

	return block, nil
}

func (c *ChainStorage) GetBlock(height int) (*Block, error) {
	return c.Block(height)
}

func (c *ChainStorage) IndexFile() string {
	return filepath.Join(MainChain(), "index")
}

func (c *ChainStorage) DataFile(fileID string) string {
	return filepath.Join(MainChain(), fileID)
}

func (c *ChainStorage) Touch() error {
	if err := os.MkdirAll(MainChain(), 0755); err != nil {
		return err
	}

	indexPath := c.IndexFile()

	if info, err := os.Stat(indexPath); err == nil && info.IsDir() {
		if err := os.RemoveAll(indexPath); err != nil {
			return err
		}
	}

	_, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
}

func (c *ChainStorage) ResetData(directory string) error {
	if err := os.RemoveAll(directory); err != nil {
		return err
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	_, err := os.OpenFile(c.IndexFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
}

func (c *ChainStorage) Index(needle interface{}) ([]interface{}, error) {
	if height, ok := needle.(int); ok {
		idx := c.ReadIdx(height)
		return c.ReadIndex(idx)
	}
	return c.SearchIndex(needle.(string), "")
}

func (c *ChainStorage) ReadIdx(height int) int {
	idx := 0
	lastIdx := c.LastIdx()
	println("Last index:", lastIdx)
	lastIndex, _ := c.ReadIndex(lastIdx)
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

func (c *ChainStorage) ReadIndex(idx int) ([]interface{}, error) {
	if idx > 0 {
		iseek := C.CHAIN_HEADER_BYTES + (idx-1)*C.CHAIN_HEAP_BYTES
		DebugLog(fmt.Sprintf("index seek: %d", iseek))
		raw, err := ReadPart(c.IndexFile(), int64(iseek), C.CHAIN_HEAP_BYTES)

		if err != nil {
			return nil, err
		}

		if len(raw) == C.CHAIN_HEAP_BYTES {
			return ChainIndex(raw), nil
		}
	}
	return []interface{}{}, nil
}

func (c *ChainStorage) ReadData(index []interface{}) ([]byte, error) {
	fileID := index[1].(string)
	seek := index[2].(int)
	length := index[3].(int)

	return ReadPart(c.DataFile(fileID), int64(seek), length)
}

func (c *ChainStorage) LastIdx() int {
	header, _ := ReadPart(c.IndexFile(), 0, C.CHAIN_HEADER_BYTES)
	return F.BinDec(header)
}

func (c *ChainStorage) LastIndex() ([]interface{}, error) {
	lastIdx := c.LastIdx()
	return c.ReadIndex(lastIdx)
}

func (c *ChainStorage) Write(block *Block) error {
	lastHeight := LastHeight()
	height := lastHeight + 1

	if err := c.WriteData(height, block.BlockHash(), []byte(block.Ser("full"))); err != nil {
		return err
	}

	return nil
}

func (c *ChainStorage) WriteData(height int, key string, data []byte) error {
	lastIdx := c.LastIdx()
	lastIndex, _ := c.ReadIndex(lastIdx)

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
	DebugLog(fmt.Sprintf("데이터 파일: %s, 파일ID: %s, 시크: %d, 길이: %d", c.DataFile(fileID), fileID, seek, length))
	if err := AppendFile(c.DataFile(fileID), ""); err != nil {
		return err
	}

	if C.LEDGER_FILESIZE_LIMIT < seek+length {
		fileID = DataID(fileID)
		seek = 0
	}

	headerData := c.headerRaw(idx)
	indexData := c.indexRaw(key, fileID, height, seek, length)
	DebugLog(fmt.Sprintf("파일ID: %s", fileID))

	if err := WriteFile(c.DataFile(fileID), int64(seek), data); err != nil {
		return err
	}

	if err := WriteFile(c.IndexFile(), int64(iseek), indexData); err != nil {
		return err
	}

	return WriteFile(c.IndexFile(), 0, headerData)
}

func (c *ChainStorage) RemoveData(directory string, idx int) error {
	if idx <= 0 {
		idx = 1
	}
	headerData := c.headerRaw(idx - 1)
	return WriteFile(c.IndexFile(), 0, headerData)
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
	header, err := ReadPart(c.IndexFile(), 0, C.CHAIN_HEADER_BYTES)
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
	f, err := os.Open(c.IndexFile())
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
