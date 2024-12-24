package model

import (
	"encoding/json"
	"fmt"
	"hello/pkg/core/structure"
	F "hello/pkg/util"
)

type TransactionMap = *map[string]*SignedTransaction
type UpdateMap = *map[string]Update

type BlockHeader struct {
	Height      int    `json:"height"`
	Timestamp_s int64  `json:"s_timestamp"`
	BlockRoot   string `json:"block_root"`
}

func (bh BlockHeader) Ser() string {
	j, _ := json.Marshal(bh)
	fmt.Println(string(j))
	return string(j)
}

type Block struct {
	Height            int            `json:"height"`
	Transactions      TransactionMap `json:"transactions"`
	UniversalUpdates  UpdateMap
	LocalUpdates      UpdateMap
	PreviousBlockhash string `json:"previous_blockhash"`
	Timestamp_s       int64  `json:"s_timestamp"`
	Vout              string `json:"vout"`
	Nonce             string `json:"nonce"`
	RewardAddress     string `json:"reward_address"`
	Difficulty        int    `json:"difficulty"`
}

func NewBlock(height int, previous_blockhash string) Block {
	return Block{
		Height:            height,
		PreviousBlockhash: previous_blockhash,
		Transactions:      &map[string]*SignedTransaction{},
		UniversalUpdates:  &map[string]Update{},
		LocalUpdates:      &map[string]Update{},
	}
}

func CreateBlock(
	height int,
	transactions TransactionMap,
	universalUpdates UpdateMap,
	localUpdates UpdateMap,
	previousBlockhash string,
	timestamp_s int64,
	vout string,
	nonce string,
	rewardAddress string,
) Block {
	return Block{
		Height:            height,
		Transactions:      transactions,
		UniversalUpdates:  universalUpdates,
		LocalUpdates:      localUpdates,
		PreviousBlockhash: previousBlockhash,
		Timestamp_s:       timestamp_s,
		Vout:              vout,
		Nonce:             nonce,
		RewardAddress:     rewardAddress,
	}
}

func (block *Block) SetTimestamp(timestamp int64) {
	block.Timestamp_s = timestamp
}

func (block *Block) AppendTransaction(tx SignedTransaction) error {
	txHash := tx.GetTxHash()

	(*block.Transactions)[txHash] = &tx
	return nil
}

func (block *Block) AppendLocalUpdate(update Update) bool {
	key := F.FillHash(update.Key)
	(*block.LocalUpdates)[key] = update
	return true
}

func (block *Block) AppendUniversalUpdate(update Update) bool {
	key := F.FillHash(update.Key)
	(*block.UniversalUpdates)[key] = update
	return true
}

func (block Block) BlockHeader() string {
	obj := BlockHeader{Height: block.Height, Timestamp_s: block.Timestamp_s, BlockRoot: block.BlockRoot()}
	return F.Hash(obj.Ser())
}

func (block Block) BlockRoot() string {
	s := F.Concat(block.TransactionRoot(), block.UpdateRoot())
	return F.Hash(s)
}

func (block Block) THashs() []string {
	txs := F.SortedValueK(*block.Transactions)

	return F.Map(txs, func(tx *SignedTransaction) string {
		return tx.GetTxHash()
	})
}

func (block Block) UHashs() []string {
	hashs := make(map[string]Update)
	for _, update := range *block.UniversalUpdates {
		hashs[update.GetHash()] = update
	}
	for _, update := range *block.LocalUpdates {
		hashs[update.GetHash()] = update
	}

	res := make([]string, 0)

	for _, v := range F.SortedValueK(hashs) {
		res = append(res, v.GetHash())
	}

	return res
}

func (block Block) TransactionRoot() string {
	return F.MerkleRoot(block.THashs())
}

func (block Block) UpdateRoot() string {
	return F.MerkleRoot(block.UHashs())
}

func (block Block) BlockHash() string {
	s := F.Concat(block.PreviousBlockhash, block.BlockHeader())
	return F.TimeHash(s, int64(block.Timestamp_s))
}

func (block Block) BaseObj() *structure.OrderedMap {
	om := structure.NewOrderedMap()

	// 순서가 보장되도록 순차적으로 추가
	om.Set("height", block.Height)
	om.Set("s_timestamp", block.Timestamp_s)
	om.Set("previous_blockhash", block.PreviousBlockhash)
	om.Set("blockhash", block.BlockHash())
	om.Set("difficulty", block.Difficulty)
	om.Set("reward_address", block.RewardAddress)
	om.Set("vout", block.Vout)
	om.Set("nonce", block.Nonce)

	return om
}

func convertUpdates(updates map[string]Update) *structure.OrderedMap {
	// comment: Convert updates to simplified format with only old and new values
	result := structure.NewOrderedMap()
	for key, update := range updates {
		up := structure.NewOrderedMap()
		up.Set("old", update.Old)
		up.Set("new", update.New)
		result.Set(key, up)
	}
	return result
}

func (block Block) FullObj() *structure.OrderedMap {
	obj := block.BaseObj()
	if block.Transactions != nil {
		txOrderedMap := structure.NewOrderedMap()
		txHashes := F.SortedValueK(*block.Transactions)

		for _, tx := range txHashes {
			txOrderedMap.Set(tx.GetTxHash(), tx.BaseObj())
		}

		obj.Set("transactions", txOrderedMap)
	} else {
		obj.Set("transactions", structure.NewOrderedMap())
	}

	if block.UniversalUpdates != nil {
		// comment: Convert UniversalUpdates to simplified format
		convertedUpdates := convertUpdates(*block.UniversalUpdates)
		obj.Set("universal_updates", convertedUpdates)
	} else {
		obj.Set("universal_updates", structure.NewOrderedMap())
	}

	if block.LocalUpdates != nil {
		// comment: Convert LocalUpdates to simplified format
		convertedUpdates := convertUpdates(*block.LocalUpdates)
		obj.Set("local_updates", convertedUpdates)
	} else {
		obj.Set("local_updates", structure.NewOrderedMap())
	}

	return obj
}

func (block Block) Ser(t string) string {
	if t == "full" {
		return block.FullObj().Ser()
	} else {
		return block.BaseObj().Ser()
	}
}
func (block *Block) Init() {
	if block.Transactions == nil {
		tm := new(map[string]*SignedTransaction)
		block.Transactions = tm
	}
	if block.UniversalUpdates == nil {
		uu := new(map[string]Update)
		block.UniversalUpdates = uu
	}
	if block.LocalUpdates == nil {
		m := new(map[string]Update)
		block.LocalUpdates = m
	}
}

func (block *Block) GetTimestamp() int64 {
	return block.Timestamp_s
}

func (block *Block) GetTransactionCount() int {
	return len(*block.Transactions)
}

func (block *Block) GetTransactions() map[string]*SignedTransaction {
	return *block.Transactions
}
