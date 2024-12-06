package main

import (
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	_ "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	S "hello/pkg/core/structure"
	"testing"
)

func createTestBlock2(t *testing.T) *Block {
	// Create first Send transaction
	tx1Data := S.NewOrderedMap()
	tx1Data.Set("type", "Send")
	tx1Data.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	tx1Data.Set("amount", 3142500000)
	tx1Data.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	tx1Data.Set("timestamp", int64(1731062859308000))
	tx1, err := NewSignedTransaction(tx1Data)
	if err != nil {
		t.Fatalf("Failed to create tx1: %v", err)
	}

	// Create second Send transaction
	tx2Data := S.NewOrderedMap()
	tx2Data.Set("type", "Send")
	tx2Data.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	tx2Data.Set("amount", 3142500000)
	tx2Data.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	tx2Data.Set("timestamp", int64(1731062859742000))
	tx2, err := NewSignedTransaction(tx2Data)
	if err != nil {
		t.Fatalf("Failed to create tx2: %v", err)
	}

	// Create block
	previousBlockhash := "0626647acb68c0fa085be6ebfbafdc3b3afbcde8bc0bff1ba1f9b8f49a16faded2edbee8c0abb7"
	block := NewBlock(1, previousBlockhash)
	block.SetTimestamp(1731062860000000)

	// Add Universal Updates
	block.AppendUniversalUpdate(Update{
		Key: "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770da53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4",
		Old: "99999999999999993566250000",
		New: "99999999999999990349375000",
	})

	block.AppendUniversalUpdate(Update{
		Key: "c8c603ff91a3c59d637c7bda83e732dea6ec74e1001b35600f0ba7831dbfe32900000000000000000000000000000000000000000000",
		Old: "148750000",
		New: "223125000",
	})

	block.AppendUniversalUpdate(Update{
		Key: "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf",
		Old: "6285000000",
		New: "9427500000",
	})

	block.AppendUniversalUpdate(Update{
		Key: "c5ca2cb405daf22453b559420907bb12d7fb34519ac55d81f47829054374512fa53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4",
		Old: nil,
		New: "100000000000000000000000000",
	})

	block.AppendUniversalUpdate(Update{
		Key: "87abdca0d3d3be9f71516090a362e5e79546f3183b1793789902c2e5176f62ae00000000000000000000000000000000000000000000",
		Old: "1864",
		New: "1832",
	})

	block.AppendUniversalUpdate(Update{
		Key: "fbab6eb9aa47eeb4d14b9473201b5aedbe0c03ba583be29f5840452ad2f1724200000000000000000000000000000000000000000000",
		Old: nil,
		New: "0023c5de767f70e88626023c5de767f70e88626023c5de767f70e88626023c5d",
	})

	// Add Local Updates
	block.AppendLocalUpdate(Update{
		Key: "724d2935080d38850e49b74927eb0351146c9ee955731f4ef53f24366c5eb9b100000000000000000000000000000000000000000000",
		Old: 5,
		New: 7,
	})

	block.AppendLocalUpdate(Update{
		Key: "12194c0ef66a96758afcf4e7ddd3a0b851bba110c7dd2ffff358cbabd725b3fc00000000000000000000000000000000000000000000",
		Old: nil,
		New: 1,
	})

	block.AppendLocalUpdate(Update{
		Key: "290eed314ce4d91c387028c290936b5b261e06f05d871bad42dfdf7436e89e9c00000000000000000000000000000000000000000000",
		Old: nil,
		New: "0",
	})

	// Add transactions
	block.AppendTransaction(tx1)
	block.AppendTransaction(tx2)

	return &block
}

func TestChainStorageWriteAndRead(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "chain_test"

	chain := &ChainStorage{}

	// Initialize
	err := chain.Touch()
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	block := createTestBlock2(t)
	// Write data test
	err = chain.Write(block)
	if err != nil {
		t.Errorf("Failed to write data at height %d: %v", block.Height, err)
	}

	// Read data test
	// Query index by height
	indices, err := chain.Index(block.Height)
	if err != nil {
		t.Errorf("Failed to query index at height %d: %v", block.Height, err)
	}

	// Read data
	data, err := chain.ReadData(indices)
	if err != nil {
		t.Errorf("Failed to read data at height %d: %v", block.Height, err)
	}
	parsedBlock, err := ParseBlock(data)
	if err != nil {
		t.Errorf("Failed to parse block: %v", err)
	}

	DebugLog(fmt.Sprintf("parsedBlock: %v", parsedBlock))

	// Verify block values
	if block.Height != parsedBlock.Height {
		t.Errorf("Block height mismatch. Expected: %d, Got: %d", block.Height, parsedBlock.Height)
	}

	// Last index test
	lastIdx := chain.LastIdx()
	if lastIdx != 1 {
		t.Errorf("Last index mismatch. Expected: %d, Got: %d",
			1, lastIdx)
	}

}
