package model

import (
	S "hello/pkg/core/structure"
	"hello/pkg/crypto"
	"hello/pkg/util"
)

type Ia interface{}

type SignedTransaction struct {
	Data      *S.OrderedMap `json:"data"` 
	Xpub      string        `json:"xpub"`
	Signature string        `json:"signature"`
}

func NewSignedTransaction(data *S.OrderedMap) SignedTransaction { 
	return SignedTransaction{Data: data}
}

func (tx SignedTransaction) Ser() string {
	return tx.Data.Ser()
}

func (tx SignedTransaction) GetSize() int {
	return len(tx.Ser())
}

func (tx SignedTransaction) GetTxHash() string {
	ts, _ := tx.Data.Get("timestamp")
	return util.TimeHash(util.Hash(tx.Ser()), ts.(int64))
}

func (tx *SignedTransaction) Sign(privateKey string) string {
	txHash := tx.GetTxHash()
	tx.Signature = crypto.Sign(txHash, privateKey)
	return tx.Signature
}
