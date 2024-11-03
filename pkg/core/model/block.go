package model

import (
	"encoding/json"
	F "hello/pkg/util"
)

type TransactionMap = map[string]SignedTransaction
type UpdateMap = map[string]Update

type BlockHeader struct {
	Height      int64  `json:"height"`
	Timestamp_s int64  `json:"timestamp_s"`
	BlockRoot   string `json:"block_root"`
}

func (bh BlockHeader) Ser() string {
	j, _ := json.Marshal(bh)
	return string(j)
}

type Block struct {
	Height            int64 `json:"height"`
	Transactions      TransactionMap
	UniversalUpdates  UpdateMap
	LocalUpdates      UpdateMap
	PreviousBlockhash string `json:"previous_blockhash"`
	Timestamp_s       int64  `json:"timestamp_s"`
	Vout              string `json:"vout"`
	Nonce             string `json:"nonce"`
	RewardAddress     string `json:"reward_address"`
}

func CreateBlock(
	height int64,
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
		Height: height, Transactions: transactions, UniversalUpdates: universalUpdates,
		LocalUpdates: localUpdates, PreviousBlockhash: previousBlockhash,
		Timestamp_s: timestamp_s, Vout: vout, Nonce: nonce, RewardAddress: rewardAddress,
	}
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
	txs := F.SortedValueK(block.Transactions)
	return F.Map(txs, func(tx SignedTransaction) string {
		return tx.GetTxHash()
	})
}

func (block Block) UHashs() []string {
	res := make([]string, len(block.UniversalUpdates)+len(block.LocalUpdates))
	uhashs := F.SortedValueK(block.UniversalUpdates)
	lhashs := F.SortedValueK(block.LocalUpdates)
	hashs := append(uhashs, lhashs...)

	for _, h := range hashs {
		res = append(res, h.GetHash())
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
