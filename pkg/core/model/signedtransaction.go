package model

import (
	"fmt"
	S "hello/pkg/core/structure"
	"hello/pkg/crypto"
	"hello/pkg/util"
	F "hello/pkg/util"
)

type SignedTransaction struct {
	Data      *S.OrderedMap `json:"data"`
	Xpub      string        `json:"xpub"`
	Signature string        `json:"signature"`
}

func NewSignedTransaction(data *S.OrderedMap) (SignedTransaction, error) {
	txData, ok := data.Get("transaction")

	if !ok || txData == nil {
		return SignedTransaction{}, fmt.Errorf("the signed transaction must contain the transaction parameter")
	}

	if txStr, ok := txData.(string); ok {
		fmt.Printf("txStr: %s\n", txStr)
		txMap, err := S.ParseOrderedMap(txStr)
		if err != nil {
			return SignedTransaction{}, err
		}
		data.Set("transaction", txMap)
	}

	tx := SignedTransaction{Data: data}

	publicKey, ok := data.Get("public_key")

	if !ok || publicKey == nil {
		return tx, fmt.Errorf("the signed transaction must contain the public_key parameter")
	}
	tx.Xpub = publicKey.(string)

	signature, ok := data.Get("signature")
	if !ok || signature == nil {
		return tx, fmt.Errorf("the signed transaction must contain the signature parameter")
	}

	tx.Signature = signature.(string)

	return tx, nil
}

func (tx *SignedTransaction) Ser() (string, error) {
	if tx.Xpub == "" {
		return "", fmt.Errorf("The signed transaction must contain the public_key parameter")
	}

	if tx.Signature == "" {
		return "", fmt.Errorf("The signed transaction must contain the signature parameter")
	}

	transaction, ok := tx.Data.Get("transaction")
	if !ok || transaction == nil {
		return "", fmt.Errorf("the signed transaction must contain the transaction parameter")
	}

	copy := S.NewOrderedMap()
	for _, key := range tx.Data.Keys() {
		if val, ok := tx.Data.Get(key); ok {
			copy.Set(key, val)
		}
	}

	transactionStr := transaction.(*S.OrderedMap).Ser()
	copy.Set("transaction", transactionStr)

	return copy.Ser(), nil
}

func (tx *SignedTransaction) GetSize() (int, error) {
	ser, err := tx.Ser()
	if err != nil {
		return 0, err
	}
	return len(ser), nil
}

func (tx *SignedTransaction) GetTimestamp() int64 {
	transaction, ok := tx.Data.Get("transaction")

	if !ok || transaction == nil {
		return 0
	}

	timestamp, ok := transaction.(*S.OrderedMap).Get("timestamp")

	if !ok {
		return 0
	}

	timestampInt64, ok := timestamp.(int64)
	if !ok {
		if timestampInt, ok := timestamp.(int); ok {
			timestampInt64 = int64(timestampInt)
		} else {
			return 0
		}
	}

	return timestampInt64
}

func (tx *SignedTransaction) GetCID() string {
	transaction, ok := tx.Data.Get("transaction")

	if !ok || transaction == nil {
		return ""
	}

	cid, ok := transaction.(*S.OrderedMap).Get("cid")
	if !ok {
		return F.RootSpaceId()
	}
	return cid.(string)
}

func (tx *SignedTransaction) GetType() string {
	transaction, ok := tx.Data.Get("transaction")
	if !ok || transaction == nil {
		return ""
	}

	txType, ok := transaction.(*S.OrderedMap).Get("type")
	if !ok {
		return ""
	}
	return txType.(string)
}

func (tx *SignedTransaction) GetTxData() *SignedData {
	transaction, ok := tx.Data.Get("transaction")
	if !ok || transaction == nil {
		return nil
	}

	return NewSignedDataFromTransaction(tx)
}

func (tx *SignedTransaction) GetTxHash() string {
	transaction, ok := tx.Data.Get("transaction")
	if !ok || transaction == nil {
		return ""
	}

	timestamp, ok := transaction.(*S.OrderedMap).Get("timestamp")
	if !ok {
		return ""
	}

	ser := transaction.(*S.OrderedMap).Ser()

	timestampInt64, ok := timestamp.(int64)
	if !ok {
		if timestampInt, ok := timestamp.(int); ok {
			timestampInt64 = int64(timestampInt)
		} else {
			return ""
		}
	}

	return util.TimeHash(util.Hash(ser), timestampInt64)
}

func (tx *SignedTransaction) Validate() error {
	if tx.Xpub == "" {
		return fmt.Errorf("the signed transaction must contain the xpub parameter")
	}

	transaction, ok := tx.Data.Get("transaction")
	if !ok || transaction == nil {
		return fmt.Errorf("the signed transaction must contain the transaction parameter")
	}

	txType, ok := transaction.(*S.OrderedMap).Get("type")
	if !ok || txType == nil {
		return fmt.Errorf("the signed transaction must contain the transaction.type parameter")
	}

	if _, ok := txType.(string); !ok {
		return fmt.Errorf("Parameter transaction.type must be of string type")
	}

	timestamp, ok := transaction.(*S.OrderedMap).Get("timestamp")
	if !ok || timestamp == nil {
		return fmt.Errorf("the signed transaction must contain the transaction.timestamp parameter")
	}

	if _, ok := timestamp.(int64); !ok {
		if _, ok := timestamp.(int); !ok {
			return fmt.Errorf("parameter transaction.timestamp must be of integer type")
		}
	}

	if tx.Signature == "" {
		return fmt.Errorf("the signed transaction must contain the signature parameter")
	}

	// Verify signature
	hash := tx.GetTxHash()

	if !crypto.SignatureValidity(hash, tx.Xpub, tx.Signature) {
		return fmt.Errorf("invalid signature: %s (transaction hash: %s)", tx.Signature, hash, tx.Xpub)
	}

	return nil
}

func (tx *SignedTransaction) GetXpub() string {
	return tx.Xpub
}

func (tx *SignedTransaction) GetSignature() string {
	return tx.Signature
}
