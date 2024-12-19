package model

import (
    "fmt"
    S "hello/pkg/core/structure"
    "hello/pkg/crypto"
    "hello/pkg/util"
)

// EncodedTransaction represents a transaction with encoding and cryptographic attributes.
type EncodedTransaction struct {
    Data      *S.OrderedMap `json:"data"`
    Xpub      string        `json:"xpub"`
    Signature string        `json:"signature"`
}

// BaseObj returns the base data object of the transaction.
func (tx *EncodedTransaction) BaseObj() *S.OrderedMap {
    return tx.Data
}

// NewEncodedTransaction initializes an EncodedTransaction from given data.
func NewEncodedTransaction(data *S.OrderedMap) (EncodedTransaction, error) {
    txData, hasTransaction := data.Get("transaction")

    if !hasTransaction || txData == nil {
        return EncodedTransaction{}, fmt.Errorf("the transaction must contain the transaction data")
    }

    tx := EncodedTransaction{Data: data}

    if publicKey, ok := data.Get("public_key"); ok && publicKey != nil {
        tx.Xpub = publicKey.(string)
    } else {
        return tx, fmt.Errorf("the transaction must contain the public_key parameter")
    }

    if signature, ok := data.Get("signature"); ok && signature != nil {
        tx.Signature = signature.(string)
    } else {
        return tx, fmt.Errorf("the transaction must contain the signature parameter")
    }

    return tx, nil
}

// Serialize converts the transaction into a string representation.
func (tx *EncodedTransaction) Serialize() (string, error) {
    if tx.Xpub == "" || tx.Signature == "" {
        return "", fmt.Errorf("the transaction must contain both xpub and signature parameters")
    }

    transaction, ok := tx.Data.Get("transaction")
    if !ok || transaction == nil {
        return "", fmt.Errorf("the transaction must contain the transaction data")
    }

    copy := S.NewOrderedMap()
    for _, key := range tx.Data.Keys() {
        if val, ok := tx.Data.Get(key); ok {
            copy.Set(key, val)
        }
    }

    return copy.Ser(), nil
}

// GetSize returns the size of the serialized transaction data.
func (tx *EncodedTransaction) GetSize() (int, error) {
    ser, err := tx.Serialize()
    if err != nil {
        return 0, err
    }
    return len(ser), nil
}

// GetTimestamp retrieves the timestamp from the transaction data.
func (tx *EncodedTransaction) GetTimestamp() int64 {
    return getInt64ValueFromTransaction(tx.Data, "timestamp")
}

// GetCID retrieves the CID from the transaction data, falling back to RootSpaceId if not found.
func (tx *EncodedTransaction) GetCID() string {
    if cid, found := getStringValueFromTransaction(tx.Data, "cid"); found {
        return cid
    }
    return util.RootSpaceId()
}

// GetType retrieves the type from the transaction data.
func (tx *EncodedTransaction) GetType() string {
    if txType, found := getStringValueFromTransaction(tx.Data, "type"); found {
        return txType
    }
    return ""
}

// GetTxHash generates and returns the hash of the transaction based on its serialized data and timestamp.
func (tx *EncodedTransaction) GetTxHash() string {
    transaction, ok := tx.Data.Get("transaction")
    if !ok || transaction == nil {
        return ""
    }

    timestamp, ok := transaction.(*S.OrderedMap).Get("timestamp")
    if !ok {
        return ""
    }

    ser := transaction.(*S.OrderedMap).Ser()
    timestampInt64 := getInt64Value(timestamp)
   
    return util.TimeHash(util.Hash(ser), timestampInt64)
}

// Validate checks the validity of the transaction and its signature.
func (tx EncodedTransaction) Validate() error {
    if err := tx.basicValidation(); err != nil {
        return err
    }

    // Verify the signature
    if !crypto.SignatureValidity(tx.GetTxHash(), tx.Xpub, tx.Signature) {
        return fmt.Errorf("invalid signature: %s", tx.Signature)
    }

    return nil
}

func (tx EncodedTransaction) basicValidation() error {
    if tx.Xpub == "" {
        return fmt.Errorf("the transaction must contain the xpub parameter")
    }

    transaction, hasTransaction := tx.Data.Get("transaction")
    if !hasTransaction || transaction == nil {
        return fmt.Errorf("the transaction must include the transaction data")
    }

    if txType, ok := transaction.(*S.OrderedMap).Get("type"); !ok || txType == nil {
        return fmt.Errorf("the transaction must contain the transaction.type parameter")
    }

    timestamp, ok := transaction.(*S.OrderedMap).Get("timestamp")
    if !ok || timestamp == nil {
        return fmt.Errorf("the transaction must contain the transaction.timestamp parameter")
    }

    timestampTypeOk := okInt64OrInt(timestamp)
    if !timestampTypeOk {
        return fmt.Errorf("parameter transaction.timestamp must be of integer type")
    }

    if tx.Signature == "" {
        return fmt.Errorf("the transaction must contain the signature parameter")
    }

    return nil
}

func getInt64ValueFromTransaction(data *S.OrderedMap, key string) int64 {
    transaction, ok := data.Get("transaction")
    if !ok || transaction == nil {
        return 0
    }

    if val, found := transaction.(*S.OrderedMap).Get(key); found {
        return getInt64Value(val)
    }

    return 0
}

func getStringValueFromTransaction(data *S.OrderedMap, key string) (string, bool) {
    transaction, ok := data.Get("transaction")
    if !ok || transaction == nil {
        return "", false
    }

    if val, found := transaction.(*S.OrderedMap).Get(key); found {
        return val.(string), true
    }

    return "", false
}

func getInt64Value(val interface{}) int64 {
    if v, ok := val.(int64); ok {
        return v
    }
    if v, ok := val.(int); ok {
        return int64(v)
    }
    return 0
}

func okInt64OrInt(val interface{}) bool {
    _, isInt64 := val.(int64)
    _, isInt := val.(int)
    return isInt64 || isInt
}

// Sign signs the transaction with the given private key and updates the signature.
func (tx *EncodedTransaction) Sign(privateKey, publicKey string) string {
    hash := tx.GetTxHash()
    signature := crypto.Signature(hash, privateKey)
    tx.Signature = signature
    tx.Xpub = publicKey
    return signature
}

// ParseSerialized converts a serialized string into an EncodedTransaction object.
func ParseSerialized(serialized string) (EncodedTransaction, error) {
    data, err := S.ParseOrderedMap(serialized)
    if err != nil {
        return EncodedTransaction{}, fmt.Errorf("failed to parse serialized data: %v", err)
    }

    if txStr, ok := data.Get("transaction"); ok && txStr != nil {
        txMap, parseErr := S.ParseOrderedMap(txStr.(string))
        if parseErr != nil {
            return EncodedTransaction{}, fmt.Errorf("failed to parse transaction data: %v", parseErr)
        }
        data.Set("transaction", txMap)
    } else {
        return EncodedTransaction{}, fmt.Errorf("missing transaction in serialized data")
    }

    return NewEncodedTransaction(data)
}

// ParseSerializedEscaped creates an EncodedTransaction from a serialized string with escaped characters.
func ParseSerializedEscaped(serialized string) (EncodedTransaction, error) {
    data, err := S.ParseOrderedMap(serialized)
    if err != nil {
        return EncodedTransaction{}, fmt.Errorf("failed to parse serialized data: %v", err)
    }

    if txMap, ok := data.Get("transaction"); ok && txMap != nil {
        data.Set("transaction", txMap)
    } else {
        return EncodedTransaction{}, fmt.Errorf("missing transaction in serialized data")
    }

    return NewEncodedTransaction(data)
}