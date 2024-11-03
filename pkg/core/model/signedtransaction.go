package model

import (
	"encoding/json"
	S "hello/pkg/core/structure"
	"hello/pkg/crypto"
	"hello/pkg/util"
)

type Ia interface{}

type SignedTransaction struct {
	Data      *S.OrderedMap `json:"data"` // OrderedMap으로 변경
	Xpub      string        `json:"xpub"`
	Signature string        `json:"signature"`
}

func NewSignedTransaction(data *S.OrderedMap) SignedTransaction { // OrderedMap을 인자로 받음
	return SignedTransaction{Data: data}
}

func (tx SignedTransaction) Ser() string {
	j, _ := json.Marshal(tx.Data) // OrderedMap을 JSON으로 직렬화
	return string(j)
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
