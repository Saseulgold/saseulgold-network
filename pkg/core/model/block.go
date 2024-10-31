package model

import (
	"encoding/json"
	F "hello/pkg/util"
)

type TransactionMap = map[string]SignedTransaction
type UpdateMap = map[string]SignedTransaction

type BlockHeader struct {
	height      int64  `json:"height"`
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
	Updates           UpdateMap
	LocalUpdates      UpdateMap
	PreviousBlockhash string `json:"previous_blockhash"`
	Timestamp_s       int64  `json:"timestamp_s"`
}

func (block Block) BlockHeader() string {
	obj := BlockHeader{height: block.Height, Timestamp_s: block.Timestamp_s, BlockRoot: block.BlockRoot()}
	return F.Hash(obj.Ser())
}

func (block Block) BlockRoot() string {

	s := F.Concat(block.TransactionRoot(), block.UpdateRoot())
	return F.Hash(s)
}

func (block Block) THashs() []string {
	txs := F.SortedKeys(block.Transactions)
	F.Map(txs, func(a SignedTransaction) string {
		return a.GetTxHash()
	})
}

func (block Block) UHashs() []string {
	res := make([]string, len(block.Updates)+len(block.LocalUpdates))
	uhashs := F.SortedKeys(block.Updates)
	lhashs := F.SortedKeys(block.LocalUpdates)
	hashs := append(uhashs[:], lhashs[:])

	for _, h := range hashs {
		u := uhashs[i]
		res = append(res, u.Hash())
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
