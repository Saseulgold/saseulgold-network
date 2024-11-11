package main

import (
	_ "fmt"
	. "hello/pkg/core/model"
	S "hello/pkg/core/structure"
	"testing"
)

func TestUpdate_GetHash0(t *testing.T) {
	oldValue := "99999999999999996783125000"
	newValue := "99999999999999993566250000"
	update := Update{
		Key: "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770da53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4",
		Old: oldValue,
		New: newValue,
	}

	expectedHash := "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770da53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c49b7a4cfee50af6cfe8d35060e2ed250039a31ad30d18a02f8b5f7934cd2004f6"
	actualHash := update.GetHash()

	if actualHash != expectedHash {
		t.Errorf("GetHash() = %v; want %v", actualHash, expectedHash)
	}
}

func TestBlock_WithMultipleUpdates(t *testing.T) {
	// 첫 번째 Send 트랜잭션 생성
	tx1Data := S.NewOrderedMap()
	tx1Data.Set("type", "Send")
	tx1Data.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	tx1Data.Set("amount", 3142500000)
	tx1Data.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	tx1Data.Set("timestamp", int64(1731062859308000))
	tx1 := NewSignedTransaction(tx1Data)

	// 두 번째 Send 트랜잭션 생성
	tx2Data := S.NewOrderedMap()
	tx2Data.Set("type", "Send")
	tx2Data.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	tx2Data.Set("amount", 3142500000)
	tx2Data.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	tx2Data.Set("timestamp", int64(1731062859742000))
	tx2 := NewSignedTransaction(tx2Data)

	// 트랜잭션 해시 검증
	expectedTx1Hash := "0626647ad01be020dd16f9dc2cbdeeb44bf7046f0fcef1ceb64b7f5ff4d6868c8ed6350e12e2ad"
	actualTx1Hash := tx1.GetTxHash()
	if actualTx1Hash != expectedTx1Hash {
		t.Errorf("Transaction 1 Hash = %v; want %v", actualTx1Hash, expectedTx1Hash)
	}

	expectedTx2Hash := "0626647ad6bb308377ef1921bbaba020038cb9a14a75aa3c0bf79c4e2e0b13f90ec769d94a4e4d"
	actualTx2Hash := tx2.GetTxHash()
	if actualTx2Hash != expectedTx2Hash {
		t.Errorf("Transaction 2 Hash = %v; want %v", actualTx2Hash, expectedTx2Hash)
	}

	// 블록 생성
	previousBlockhash := "0626647acb68c0fa085be6ebfbafdc3b3afbcde8bc0bff1ba1f9b8f49a16faded2edbee8c0abb7"
	block := NewBlock(5, previousBlockhash)
	block.SetTimestamp(1731062860000000)

	// Universal Updates 추가
	// Universal Updates 추가
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

	// Local Updates 추가
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

	// 트랜잭션 추가
	block.AppendTransaction(tx1)
	block.AppendTransaction(tx2)

	// UHash 리스트 검증
	expectedUHashes := []string{
		"12194c0ef66a96758afcf4e7ddd3a0b851bba110c7dd2ffff358cbabd725b3fc000000000000000000000000000000000000000000002676e636c7fe20bcafc28139f910ad1aaa58ad755f98af6bffa0e00f3f803625",
		"290eed314ce4d91c387028c290936b5b261e06f05d871bad42dfdf7436e89e9c00000000000000000000000000000000000000000000d9e1873da47c578c3fcfdb3fe8a1d63742444a4ff084c5efbcaef2f3ec0b2b2f",
		"724d2935080d38850e49b74927eb0351146c9ee955731f4ef53f24366c5eb9b100000000000000000000000000000000000000000000b0959c88ebb43c49949e999f1918f76f5d3cd70fd44d5bcb336c77d9d997362c",
		"87abdca0d3d3be9f71516090a362e5e79546f3183b1793789902c2e5176f62ae00000000000000000000000000000000000000000000cd08f201c5a1102e707c971626050b21e828b977451193c61671a25fa3294faa",
		"b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf37051c0ed209c0909f567fab50dda1dd3268b3f4e9a5f5c29eaf64de268a0884",
		"b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770da53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4b139da1fb6537bdca92b9533dc005d517f9e99a955283887f5496d2f3a246701",
		"c5ca2cb405daf22453b559420907bb12d7fb34519ac55d81f47829054374512fa53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4c74020b3885e369b9c0019c8da8c3342b432c52a6ba98df94227df15d471e812",
		"c8c603ff91a3c59d637c7bda83e732dea6ec74e1001b35600f0ba7831dbfe32900000000000000000000000000000000000000000000343ea3e5441ea4f6506426204bee366c89cbff5a86cd35618d56d10e2e141cae",
		"fbab6eb9aa47eeb4d14b9473201b5aedbe0c03ba583be29f5840452ad2f172420000000000000000000000000000000000000000000034d56b9daa4df8e4c3dfaaeb71d716602a448745bd0c868e6991156994364d23",
	}

	actualUHashes := block.UHashs()
	// 길이 검증
	if len(actualUHashes) != len(expectedUHashes) {
		t.Errorf("UHash length = %v; want %v", len(actualUHashes), len(expectedUHashes))
		return
	}

	// 각 해시값 검증
	for i, expectedHash := range expectedUHashes {
		if actualUHashes[i] != expectedHash {
			t.Errorf("UHash[%d] = %v; want %v", i, actualUHashes[i], expectedHash)
		}
	}

	// 검증
	expectedBlockRoot := "7e8d8bc16377cb1157fa0cfc001ea2958ad17df7119b528fca809b6360a2c9df"
	actualBlockRoot := block.BlockRoot()
	if actualBlockRoot != expectedBlockRoot {
		t.Errorf("BlockRoot() = %v; want %v", actualBlockRoot, expectedBlockRoot)
	}

	expectedTxRoot := "61c1ced40de12595e08f5858aa7e38bbb57db57456bd4fbc4b7bdf0f298c515a"
	actualTxRoot := block.TransactionRoot()
	if actualTxRoot != expectedTxRoot {
		t.Errorf("TransactionRoot() = %v; want %v", actualTxRoot, expectedTxRoot)
	}

	expectedUpdateRoot := "044c6a958a5a94594d1488f7256d264b9032fb8795202e8eab8aadb4fa91e541"
	actualUpdateRoot := block.UpdateRoot()
	if actualUpdateRoot != expectedUpdateRoot {
		t.Errorf("UpdateRoot() = %v; want %v", actualUpdateRoot, expectedUpdateRoot)
	}

	// Print block data
	t.Logf("\n=== Block Data Verification Results ===")
	t.Logf("Height: %d", block.Height)
	t.Logf("Previous Block Hash: %s", block.PreviousBlockhash)
	t.Logf("Timestamp: %d", block.Timestamp_s)
	t.Logf("Block Root: %s", block.BlockRoot())
	t.Logf("Transaction Root: %s", block.TransactionRoot())
	t.Logf("Update Root: %s", block.UpdateRoot())

	t.Logf("\n=== Transaction List ===")
	for txHash, tx := range block.Transactions {
		t.Logf("Transaction Hash: %s", txHash)
		if val, _ := tx.Data.Get("type"); val != nil {
			t.Logf("- Type: %v", val)
		}
		if val, _ := tx.Data.Get("from"); val != nil {
			t.Logf("- From Address: %v", val)
		}
		if val, _ := tx.Data.Get("to"); val != nil {
			t.Logf("- To Address: %v", val)
		}
		if val, _ := tx.Data.Get("amount"); val != nil {
			t.Logf("- Amount: %v", val)
		}
		if val, _ := tx.Data.Get("timestamp"); val != nil {
			t.Logf("- Timestamp: %v\n", val)
		}
	}

	t.Logf("\n=== Universal Updates List ===")
	for _, update := range block.UniversalUpdates {
		t.Logf("Key: %s", update.Key)
		t.Logf("Old Value: %v", update.Old)
		t.Logf("New Value: %v\n", update.New)
	}

	t.Logf("\n=== Local Updates List ===")
	for _, update := range block.LocalUpdates {
		t.Logf("Key: %s", update.Key)
		t.Logf("Old Value: %v", update.Old)
		t.Logf("New Value: %v\n", update.New)
	}
}
