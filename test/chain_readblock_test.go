package main

import (
	C "hello/pkg/core/config"
	"hello/pkg/core/storage"
	"testing"
)

func TestChainStorage_Block(t *testing.T) {
	// 테스트 설정
	C.IS_TEST = true

	cs := &storage.ChainStorage{}
	testDir := "main_chain"
	testIndexDir := "main_chain"

	last := cs.LastIdx("main_chain")
	t.Logf("last idx: %v", last)

	// 테스트 케이스들
	tests := []struct {
		name    string
		height  int // 블록 높이로 인덱스를 읽어옴
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

	// 각 테스트 케이스에 대해 실제 인덱스 데이터 준비
	for i := range tests {
		needles, err := cs.Index(testDir, tests[i].height)
		if err != nil {
			t.Fatalf("Failed to read index: %v", err)
		}
		t.Logf("needles: %v", needles)

		block, err := cs.Block(testDir, needles)
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
		idx := cs.ReadIdx(testDir, tests[i].height)
		t.Logf("Test case '%s' idx: %v", tests[i].name, idx)

		index, err := cs.ReadIndex(testIndexDir, idx)
		if err != nil {
			t.Fatalf("Failed to read index: %v", err)
		}
		t.Logf("Test case '%s' index: %v", tests[i].name, index)

		if tests[i].wantErr {
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
