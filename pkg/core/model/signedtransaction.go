package model

import (
	"encoding/json"
	"hello/pkg/crypto"
	"hello/pkg/util"
)

type Ia interface{}

type AttributeMap = map[string]Ia
type CachedMap = map[string]Ia
type TransactinData = map[string]Ia

type SignedTransaction struct {
	Data      TransactinData `json:"data"`
	Xpub      string         `json:"xpub"`
	Signature string         `json:"signature"`
}

func NewSignedTransaction(data TransactinData) SignedTransaction {
	return SignedTransaction{Data: data}
}

func (tx SignedTransaction) Ser() string {
	j, _ := json.Marshal(tx.Data)
	return string(j)
}

func (tx SignedTransaction) GetSize() int {
	return len(tx.Ser())
}

func (tx SignedTransaction) GetTxHash() string {
	return util.TimeHash(util.Hash(tx.Ser()), tx.Data["timestamp"].(int64))
}

func (tx *SignedTransaction) Sign(privateKey string) string {
	// xpub := crypto.GetXpub(privateKey)
	// address :

	txHash := tx.GetTxHash()
	tx.Signature = crypto.Sign(txHash, privateKey)
	return tx.Signature
}
