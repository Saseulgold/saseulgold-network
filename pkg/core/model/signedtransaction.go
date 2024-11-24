package model

import (
	"fmt"
	S "hello/pkg/core/structure"
	"hello/pkg/crypto"
	"hello/pkg/util"
)

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
	transaction, ok := tx.Data.Get("transaction")
	if !ok || transaction == nil {
		return ""
	}

	txMap, ok := transaction.(*S.OrderedMap)
	if !ok {
		return ""
	}

	timestamp, ok := txMap.Get("timestamp")
	if !ok {
		return ""
	}

	return util.TimeHash(util.Hash(tx.Ser()), timestamp.(int64))
}

func (tx *SignedTransaction) GetTimestamp() int64 {
	transaction, ok := tx.Data.Get("transaction")
	if !ok || transaction == nil {
		return 0
	}

	txMap, ok := transaction.(*S.OrderedMap)
	if !ok {
		return 0
	}

	timestamp, ok := txMap.Get("timestamp")
	if !ok {
		return 0
	}

	return timestamp.(int64)
}

func (tx *SignedTransaction) Validate() error {
	if tx.Xpub == "" {
		return fmt.Errorf("The signed transaction must contain the xpub parameter")
	}

	transaction, ok := tx.Data.Get("transaction")
	if !ok || transaction == nil {
		return fmt.Errorf("The signed transaction must contain the transaction parameter")
	}

	txMap, ok := transaction.(*S.OrderedMap)
	if !ok {
		return fmt.Errorf("Transaction parameter must be an OrderedMap")
	}

	txType, ok := txMap.Get("type")
	if !ok || txType == nil {
		return fmt.Errorf("The signed transaction must contain the transaction.type parameter")
	}

	if _, ok := txType.(string); !ok {
		return fmt.Errorf("Parameter transaction.type must be of string type")
	}

	timestamp, ok := txMap.Get("timestamp")
	if !ok || timestamp == nil {
		return fmt.Errorf("The signed transaction must contain the transaction.timestamp parameter")
	}

	if _, ok := timestamp.(int64); !ok {
		return fmt.Errorf("Parameter transaction.timestamp must be of integer type")
	}

	if tx.Signature == "" {
		return fmt.Errorf("The signed transaction must contain the signature parameter")
	}

	// Verify signature
	if !crypto.SignatureValidity(tx.Data.Ser(), tx.Xpub, tx.Signature) {
		return fmt.Errorf("Invalid signature: %s (transaction hash: %s)", tx.Signature, tx.GetTxHash())
	}

	return nil
}
