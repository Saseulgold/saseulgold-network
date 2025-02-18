package main

import (
	C "hello/pkg/core/config"
	"hello/pkg/core/storage"
	"testing"
)

func TestChainStorage_Block(t *testing.T) {
	// Test Settings
	C.CORE_TEST_MODE = true

	cs := &storage.ChainStorage{}
	last := cs.LastIdx()
	t.Logf("last idx: %v", last)

	// test cases
	tests := []struct {
		name    string
		height  int // Read index to block height
		wantErr bool
	}{
		{
			name:    "genesis block",
			height:  3,
			wantErr: false,
		},
		{
			name:    "genesis block",
			height:  4,
			wantErr: false,
		},
		{
			name:    "genesis block",
			height:  5,
			wantErr: false,
		},
		{
			name:    "genesis block",
			height:  5,
			wantErr: false,
		},
	}

	// Prepare real index data for each test case
	for i := range tests {
		needles, err := cs.Index(tests[i].height)
		if err != nil {
			t.Fatalf("Failed to read index: %v", err)
		}
		t.Logf("needles: %v", needles)

		block, err := cs.Block(tests[i].height)
		if err != nil {
			t.Fatalf("Failed to read block: %v", err)
		}
		block.Init()

		if err != nil {
			t.Fatalf("Failed to read block: %v", err)
		}
		t.Logf("height: %v", block.Height)
		t.Logf("timestamp: %v", block.Timestamp_s)
		t.Logf("previous blockhash: %v", block.PreviousBlockhash)
		t.Logf("universal updates: %v", block.UniversalUpdates)
		t.Logf("local updates: %v", block.LocalUpdates)
		t.Logf("transactions: %v", block.Transactions)
		// t.Logf("blockhash: %v", block.BlockHash())
		t.Logf("difficulty: %v", block.Difficulty)
		t.Logf("reward address: %v", block.RewardAddress)
		t.Logf("vout: %v", block.Vout)
		t.Logf("nonce: %v", block.Nonce)
		idx := cs.ReadIdx(tests[i].height)
		t.Logf("Test case '%s' idx: %v", tests[i].name, idx)

		index, err := cs.ReadIndex(idx)
		if err != nil {
			t.Fatalf("Failed to read index: %v", err)
		}
		t.Logf("Test case '%s' index: %v", tests[i].name, index)

		if tests[i].wantErr {
			if tests[i].height != int(block.Height) {
				t.Errorf("Expected height: %d, got: %d", tests[i].height, block.Height)
			}
			if err == nil {
				t.Errorf("Expected error but got none")
			}
			return
		}

		if err != nil {
			t.Errorf("Unexpected error occurred: %v", err)
			return
		}

		// Validate block data
		if block == nil {
			t.Error("Block is nil")
			return
		}

		// Validate specific block fields if needed
		// e.g. block.Height, block.Hash etc
	}

}
