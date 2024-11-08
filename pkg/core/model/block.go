package model

import (
	"encoding/json"
	"fmt"
	F "hello/pkg/util"
)

type TransactionMap = map[string]SignedTransaction
type UpdateMap = map[string]Update

type BlockHeader struct {
	Height      int64  `json:"height"`
	Timestamp_s int64  `json:"s_timestamp"`
	BlockRoot   string `json:"block_root"`
}

func (bh BlockHeader) Ser() string {
	j, _ := json.Marshal(bh)
	fmt.Println("bh Ser:", string(j))
	return string(j)
}

type Block struct {
	Height            int64 `json:"height"`
	Transactions      *TransactionMap
	UniversalUpdates  *UpdateMap
	LocalUpdates      *UpdateMap
	PreviousBlockhash string `json:"previous_blockhash"`
	Timestamp_s       int64  `json:"s_timestamp"`
	Vout              string `json:"vout"`
	Nonce             string `json:"nonce"`
	RewardAddress     string `json:"reward_address"`
}

func NewBlock(height int64, previous_blockhash string) Block {
	tm := make(TransactionMap, 8)
	uu := make(UpdateMap, 16)
	lu := make(UpdateMap, 4)

	return Block{Height: height, PreviousBlockhash: previous_blockhash, Transactions: &tm, UniversalUpdates: &uu, LocalUpdates: &lu}
}

func CreateBlock(
	height int64,
	transactions *TransactionMap,
	universalUpdates *UpdateMap,
	localUpdates *UpdateMap,
	previousBlockhash string,
	timestamp_s int64,
	vout string,
	nonce string,
	rewardAddress string,
) Block {
	return Block{
		Height: height, Transactions: transactions, UniversalUpdates: universalUpdates,
		LocalUpdates: localUpdates, PreviousBlockhash: previousBlockhash,
		Timestamp_s: timestamp_s, Vout: vout, Nonce: nonce, RewardAddress: rewardAddress,
	}
}

func (block *Block) SetTimestamp(timestamp int64) {
	block.Timestamp_s = timestamp
}

func (block *Block) AppendTransaction(tx SignedTransaction) bool {
	txHash := tx.GetTxHash()
	(*block.Transactions)[txHash] = tx
	return true
}

func (block *Block) AppendLocalUpdate(update Update) bool {
	updateHash := update.GetHash()
	(*block.LocalUpdates)[updateHash] = update
	return true
}

func (block *Block) AppendUniversalUpdate(update Update) bool {
	updateHash := update.GetHash()
	(*block.UniversalUpdates)[updateHash] = update
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
	return F.Map(txs, func(tx SignedTransaction) string {
		return tx.GetTxHash()
	})
}

func (block Block) UHashs() []string {
	hashs := make(map[string]Update)
	for k, v := range *block.UniversalUpdates {
		hashs[k] = v
	}
	for k, v := range *block.LocalUpdates {
		hashs[k] = v
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
	return F.TimeHash(s, block.Timestamp_s)
}
