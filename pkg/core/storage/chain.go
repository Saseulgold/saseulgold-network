package storage

import (
	"bytes"
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	"hello/pkg/core/structure"
	F "hello/pkg/util"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
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
		height, _ = strconv.Atoi(string(data))
	}
	return height
}

func (c *ChainStorage) LastBlock() (*Block, error) {
	return c.GetBlock(LastHeight())
}

func (c *ChainStorage) SetLastHeight(height int) error {
	DebugLog(fmt.Sprintf("Set last height: %d", height))
	heightStr := fmt.Sprintf("%d", height)
	return os.WriteFile(ChainInfo(), []byte(heightStr), 0644)
}

func ConfirmedHeight() int {
	return LastHeight() - 5
}

func ParseBlock(data []byte) (*Block, error) {
	om, err := structure.ParseOrderedMap(string(data))
	if err != nil {
		return nil, err
	}

	var block Block
	blockHeight, _ := om.Get("height")

	block.Height = int(blockHeight.(int64))

	blockTimestamp, _ := om.Get("s_timestamp")
	block.Timestamp_s = blockTimestamp.(int64)

	blockPreviousBlockhash, _ := om.Get("previous_blockhash")
	block.PreviousBlockhash = blockPreviousBlockhash.(string)

	blockDifficulty, _ := om.Get("difficulty")
	if blockDifficulty == nil {
		block.Difficulty = 0
	} else {
		block.Difficulty = int(blockDifficulty.(int64))
	}


	if universalUpdatesRaw, exists := om.Get("universal_updates"); exists {
		updates := make(map[string]Update)
		block.UniversalUpdates = &updates

		universalUpdates, ok := universalUpdatesRaw.(*structure.OrderedMap)
		if ok {
			for _, key := range universalUpdates.Keys() {
				value, _ := universalUpdates.Get(key)
				valueMap := value.(*structure.OrderedMap)
				old, _ := valueMap.Get("old")
				new, _ := valueMap.Get("new")
				update := Update{
					Key: key,
					Old: old,
					New: new,
				}
				updates[key] = update
			}
		}
	}

	if localUpdatesRaw, exists := om.Get("local_updates"); exists {
		updates := make(map[string]Update)
		block.LocalUpdates = &updates
		localUpdates := localUpdatesRaw.(*structure.OrderedMap)

		for _, key := range localUpdates.Keys() {
			value, _ := localUpdates.Get(key)
			valueMap := value.(*structure.OrderedMap)
			old, _ := valueMap.Get("old")
			new, _ := valueMap.Get("new")
			update := Update{
				Key: key,
				Old: old,
				New: new,
			}
			updates[key] = update
		}
	}

	if transactionsRaw, exists := om.Get("transactions"); exists {
		transactions := make(map[string]*SignedTransaction)

		transactionsRaw, ok := transactionsRaw.(*structure.OrderedMap)
		if ok {
			for _, key := range transactionsRaw.Keys() {
				value, _ := transactionsRaw.Get(key)
				data := value.(*structure.OrderedMap)
				tx, err := NewSignedTransaction(data)
				if err != nil {
					return nil, err
				}
				transactions[key] = &tx
			}
		}
		block.Transactions = &transactions
	}

	block.Init()
	return &block, nil
}

func (c *ChainStorage) Block(needle interface{}) (*Block, error) {
	if height, ok := needle.(int); ok {

		idx := c.ReadIdx(height)
		index, err := c.ReadIndex(idx)
		if err != nil {
			return nil, err
		}

		data, err := c.ReadData(index)
		if err != nil {
			return nil, err
		}

		return ParseBlock(data)
	}

	data, err := c.ReadData(needle.(ChainIndexCursor))
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

func (c *ChainStorage) Index(needle interface{}) (ChainIndexCursor, error) {
	if height, ok := needle.(int); ok {
		idx := c.ReadIdx(height)
		index, err := c.ReadIndex(idx)
		DebugLog(fmt.Sprintf("index: %v", index))
		if err != nil {
			return ChainIndexCursor{}, err
		}
		return index, nil
	}
	return c.SearchIndex(needle.(string), "")
}

func (c *ChainStorage) ReadIdx(height int) int {
	idx := 0
	lastIdx := c.LastIdx()

	lastIndex, _ := c.ReadIndex(lastIdx)
	lastHeight := 0
	if lastIndex.Height > 0 {
		lastHeight = lastIndex.Height
	}
	gap := lastHeight - height

	if gap >= 0 {
		idx = lastIdx - gap
	}

	return idx
}

func (c *ChainStorage) ReadIndex(idx int) (ChainIndexCursor, error) {
	if idx > 0 {
		iseek := C.CHAIN_HEADER_BYTES + (idx-1)*C.CHAIN_HEAP_BYTES
		DebugLog(fmt.Sprintf("index seek: %d", iseek))
		raw, err := ReadPart(c.IndexFile(), int64(iseek), C.CHAIN_HEAP_BYTES)

		if err != nil {
			return ChainIndexCursor{}, err
		}

		if len(raw) == C.CHAIN_HEAP_BYTES {
			index := ChainIndex(raw)
			return index, nil
		}
	}
	return ChainIndexCursor{}, nil
}

func (c *ChainStorage) ReadData(index ChainIndexCursor) ([]byte, error) {
	fileID := index.FileID
	seek := index.Seek
	length := index.Length

	return ReadPart(c.DataFile(fileID), seek, int(length))
}

func (c *ChainStorage) LastIdx() int {
	header, _ := ReadPart(c.IndexFile(), 0, C.CHAIN_HEADER_BYTES)
	return F.BinDec(header)
}

func (c *ChainStorage) LastIndex() (ChainIndexCursor, error) {
	lastIdx := c.LastIdx()
	return c.ReadIndex(lastIdx)
}

func (c *ChainStorage) Write(block *Block) error {
	lastHeight := LastHeight()
	height := lastHeight + 1

	DebugLog(fmt.Sprintf("Write block: %v", block.Ser("full")))

	if err := c.WriteData(height, block.BlockHash(), []byte(block.Ser("full"))); err != nil {
		return err
	}
	return c.SetLastHeight(height)
}

func (c *ChainStorage) WriteData(height int, key string, data []byte) error {
	lastIdx := c.LastIdx()
	lastIndex, _ := c.ReadIndex(lastIdx)

	var lastHeight int
	if lastIndex.Height > 0 {
		lastHeight = lastIndex.Height
	}

	if height != lastHeight+1 {
		return nil
	}

	var fileID string
	var lastSeek, lastLength int64

	if lastIndex.FileID != "" {
		fileID = lastIndex.FileID
		lastSeek = lastIndex.Seek
		lastLength = lastIndex.Length
	} else {
		fileID = DataID("")
	}

	seek := lastSeek + lastLength
	length := len(data)

	idx := lastIdx + 1
	iseek := C.CHAIN_HEADER_BYTES + lastIdx*C.CHAIN_HEAP_BYTES
	if err := AppendFile(c.DataFile(fileID), ""); err != nil {
		return err
	}

	if C.LEDGER_FILESIZE_LIMIT < seek+int64(length) {
		fileID = DataID(fileID)
		seek = 0
	}

	headerData := c.headerRaw(idx)
	indexData := c.indexRaw(key, fileID, height, int(seek), length)

	if err := WriteFile(c.DataFile(fileID), int64(seek), data); err != nil {
		return err
	}

	if err := WriteFile(c.IndexFile(), int64(iseek), indexData); err != nil {
		return err
	}

	DebugLog(fmt.Sprintf("write block index - key: %s, fileId: %s, height: %d, seek: %d, length: %d, iseek: %d, indexData length: %d", key, fileID, height, seek, length, iseek, len(indexData)))
	DebugLog(fmt.Sprintf("headerData: %v", headerData))
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

func (c *ChainStorage) SearchIndex(directory string, hash string) (ChainIndexCursor, error) {
	if len(hash) < C.HEX_TIME_SIZE {
		return ChainIndexCursor{}, nil
	}

	target := F.Hex2Bin(hash[:C.HEX_TIME_SIZE])
	header, err := ReadPart(c.IndexFile(), 0, C.CHAIN_HEADER_BYTES)
	if err != nil {
		return ChainIndexCursor{}, err
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
		return ChainIndexCursor{}, err
	}
	defer f.Close()

	for i := 0; i < cycle; i++ {
		mid = (a + b) / 2

		var l []byte
		if mid <= 0 {
			l = make([]byte, bytes)
		} else {
			if _, err := f.Seek(int64(C.CHAIN_HEADER_BYTES+(mid-1)*C.CHAIN_HEAP_BYTES), 0); err != nil {
				return ChainIndexCursor{}, err
			}
			l = make([]byte, bytes)
			if _, err := f.Read(l); err != nil {
				return ChainIndexCursor{}, err
			}
		}

		var r []byte
		if mid >= max {
			r = bytes_repeat(0xff, bytes)
		} else {
			if _, err := f.Seek(int64(C.CHAIN_HEADER_BYTES+mid*C.CHAIN_HEAP_BYTES), 0); err != nil {
				return ChainIndexCursor{}, err
			}
			r = make([]byte, bytes)
			if _, err := f.Read(r); err != nil {
				return ChainIndexCursor{}, err
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
		return ChainIndexCursor{}, err
	}
	read := make([]byte, C.CHAIN_HEAP_BYTES)
	n, err := f.Read(read)
	if err != nil && err != io.EOF {
		return ChainIndexCursor{}, err
	}

	if n == C.CHAIN_HEAP_BYTES {
		return ChainIndex(read), nil
	}

	return ChainIndexCursor{}, nil
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

func (c *ChainStorage) GetLastHeight() int {
	return LastHeight()
}

func (c *ChainStorage) GetLastBlock() (*Block, error) {
	return c.GetBlock(LastHeight())
}
