package main

import (
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	S "hello/pkg/core/structure"
	. "hello/pkg/core/vm"
	"testing"
)

func TestCommit(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "commit_test"

	sf := GetStatusFileInstance()
	ci := GetChainStorageInstance()

	sf.Touch()
	ci.Touch()

	machine := GetMachineInstance()

	block := createTestBlock1(t)
	hash := block.BlockHash()
	t.Logf("Block hash: %s", hash)

	err := machine.Commit(block)
	if err != nil {
		t.Errorf("Error committing block: %v", err)
	}
}

func createTestBlock1(t *testing.T) *Block {
	// Create first Send transaction
	data1 := S.NewOrderedMap()
	tx1Data := S.NewOrderedMap()
	tx1Data.Set("type", "Send")
	tx1Data.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	tx1Data.Set("amount", 3142500000)
	tx1Data.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	tx1Data.Set("timestamp", int64(1731062859308000))

	data1.Set("transaction", tx1Data)
	data1.Set("public_key", "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2")
	data1.Set("signature", "b1adf108db92fe4e062bd7e79c4d48b8202267f2968740a59fbdaf66f6f826768d7fb740cec4075dbe71abf4143aed9fa901e8b90b357ded06a5647ecb74da0c")

	tx1, err := NewSignedTransaction(data1)
	if err != nil {
		panic(err)
	}

	// Create second Send transaction
	data2 := S.NewOrderedMap()

	tx2Data := S.NewOrderedMap()
	tx2Data.Set("type", "Send")
	tx2Data.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	tx2Data.Set("amount", 3142500000)
	tx2Data.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	tx2Data.Set("timestamp", int64(1731062859742000))

	data2.Set("transaction", tx2Data)
	data2.Set("public_key", "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2")
	data2.Set("signature", "b1adf108db92fe4e062bd7e79c4d48b8202267f2968740a59fbdaf66f6f826768d7fb740cec4075dbe71abf4143aed9fa901e8b90b357ded06a5647ecb74da0c")

	tx2, err := NewSignedTransaction(data2)
	if err != nil {
		panic(err)
	}

	// Create block
	previousBlockhash := "0626647acb68c0fa085be6ebfbafdc3b3afbcde8bc0bff1ba1f9b8f49a16faded2edbee8c0abb7"
	block := NewBlock(1, previousBlockhash)
	// block.SetTimestamp(int(util.Utime()))
	block.SetTimestamp(1731062859742000)

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
