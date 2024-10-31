package model

import (
    F "hello/pkg/util"
)

type TransactionMap = map[string]SignedTransaction
type UpdateMap = map[string]SignedTransaction

type BlockHeader struct {
    height      int64   `json:"height"`
    Timestamp_s int64   `json:"timestamp_s"`
    BlockRoot   string  `json:"block_root"`
}


type Block struct {
	Height              int64               `json:"height"`
	Transactions        TransactionMap
    Updates             UpdateMap
    LocalUpdates        UpdateMap
    PreviousBlockhash   string              `json:"previous_blockhash"`
    Timestamp_s         int64               `json:"timestamp_s"`
}

func (block Block) BlockHeader() string {
    obj := BlockHeader{ height: block.Height, Timestamp_s: block.TimeHash, BlockRoot: block.BlockRoot() }
    return F.Hash()
}

func (block Block) BlockRoot() string {
    s := F.Concat(block.THashs(), block.UHashs())
    return F.Hash(s)
}

func (block Block) THashs() string[] {
    return F.SortedKeys(Transactions)
}

func (block Block) UHashs() string[] {
    res := []
    uhashs = F.SortedKeys(block.Updates)
    lhashs = F.SortedKeys(block.LocalUpdates)
    hashs = append(uhashs[:], lhashs[:])

    for(i := 0; i < len(hashs); i++) {
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