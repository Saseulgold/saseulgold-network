package main

import (
	. "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	S "hello/pkg/core/structure"
	"hello/pkg/util"
	"testing"
)

func TestMempoolStorage_AddTransaction(t *testing.T) {
	mp := GetMempoolInstance()
	mp.Clear() // Clear before test

	tests := []struct {
		name    string
		tx      SignedTransaction
		wantErr error
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mp.AddTransaction(tt.tx)
			if err != tt.wantErr {
				t.Errorf("AddTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createTestTransaction(txType string) SignedTransaction {
	tx1Data := S.NewOrderedMap()
	tx1Data.Set("type", txType)
	tx1Data.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	tx1Data.Set("amount", 3142500000)
	tx1Data.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	tx1Data.Set("timestamp", util.CurrentTime(nil))

	tx := NewSignedTransaction(tx1Data)
	return tx
}

func TestMempoolStorage_GetTransaction(t *testing.T) {
	mp := GetMempoolInstance()
	mp.Clear()

	// Add test transaction with specific hash
	tx := createTestTransaction("tx1")
	mp.AddTransaction(tx)

	// Query existing transaction
	if got := mp.GetTransaction(tx.GetTxHash()); got == nil {
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
	mp.AddTransaction(tx)

	// Remove transaction
	mp.RemoveTransaction(tx.GetTxHash())

	// Query removed transaction
	if got := mp.GetTransaction(tx.GetTxHash()); got != nil {
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
		t.Logf("Transaction hash: %v", tx.GetTxHash())
	}

	for _, tx := range txs {
		t.Logf("Adding transaction: %v", tx)
		mp.AddTransaction(tx)
	}

	// Query transactions (should be sorted by fee)
	got := mp.GetTransactions(2)
	if len(got) != 2 {
		t.Errorf("GetTransactions() returned %d transactions, want 2", len(got))
	}
}
