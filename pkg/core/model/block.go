package model

import (
	"encoding/json"
	"fmt"
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

func (block Block) BaseObj() map[string]interface{} {
	return map[string]interface{}{
		"height":             block.Height,
		"s_timestamp":        block.Timestamp_s,
		"previous_blockhash": block.PreviousBlockhash,
		"blockhash":          block.BlockHash(),
		"difficulty":         block.Difficulty,
		"reward_address":     block.RewardAddress,
		"vout":               block.Vout,
		"nonce":              block.Nonce,
	}
}

func (block Block) FullObj() map[string]interface{} {
	obj := block.BaseObj()
	obj["transactions"] = block.Transactions
	obj["universal_updates"] = block.UniversalUpdates
	obj["local_updates"] = block.LocalUpdates
	return obj
}

func (block Block) Ser(t string) string {
	if t == "full" {
		j, _ := json.Marshal(block.FullObj())
		return string(j)
	} else {
		j, _ := json.Marshal(block.BaseObj())
		return string(j)
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
