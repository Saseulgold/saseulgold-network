package main

import (
	. "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	S "hello/pkg/core/structure"
	"testing"
)

func createTestTransaction(txType string) SignedTransaction {
	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", txType)
	txData.Set("to", "867b8991f4f2eb94398d4647f5ddc57b30f8cb36acdf")
	txData.Set("amount", "10000000000")
	txData.Set("from", "4a8fd2ebb308370a689c3ef47cb83ba182683def3d4f")
	txData.Set("timestamp", int64(1732359603011000))

	data.Set("transaction", txData)
	data.Set("public_key", "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2")
	data.Set("signature", "b1adf108db92fe4e062bd7e79c4d48b8202267f2968740a59fbdaf66f6f826768d7fb740cec4075dbe71abf4143aed9fa901e8b90b357ded06a5647ecb74da0c")

	tx, err := NewSignedTransaction(data)
	if err != nil {
		panic(err)
	}

	return tx
}

func TestMempoolStorage_GetTransaction(t *testing.T) {
	mp := GetMempoolInstance()
	mp.Clear()

	// Add test transaction with specific hash
	tx := createTestTransaction("tx1")
	mp.AddTransaction(&tx)

	// Query existing transaction
	txHash, err := tx.GetTxHash()
	if err != nil {
		panic(err)
	}
	if got := mp.GetTransaction(txHash); got == nil {
		t.Error("GetTransaction() = nil, want transaction")
	}

	// Query non-existent transaction
	if got := mp.GetTransaction("nonexistent"); got != nil {
		t.Error("GetTransaction() = transaction, want nil")
	}
}

func TestMempoolStorage_RemoveTransaction(t *testing.T) {
	mp := GetMempoolInstance()
	mp.Clear()

	// Add test transaction
	tx := createTestTransaction("tx1")
	mp.AddTransaction(&tx)

	// Remove transaction
	txHash, err := tx.GetTxHash()
	if err != nil {
		panic(err)
	}
	mp.RemoveTransaction(txHash)

	// Query removed transaction
	if got := mp.GetTransaction(txHash); got != nil {
		t.Error("RemoveTransaction() failed to remove transaction")
	}
}

func TestMempoolStorage_GetTransactions(t *testing.T) {
	mp := GetMempoolInstance()
	mp.Clear()

	// Add multiple transactions with different hashes
	txs := []SignedTransaction{
		createTestTransaction("tx1"),
		createTestTransaction("tx2"),
	}

	for _, tx := range txs {
		txHash, err := tx.GetTxHash()
		if err != nil {
			panic(err)
		}
		t.Logf("Transaction hash: %v", txHash)
	}

	for _, tx := range txs {
		t.Logf("Adding transaction: %v", tx)
		mp.AddTransaction(&tx)
	}

	// Query transactions (should be sorted by fee)
	got := mp.GetTransactions(2)
	if len(got) != 2 {
		t.Errorf("GetTransactions() returned %d transactions, want 2", len(got))
	}
}
