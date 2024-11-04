package main

import (
	"encoding/json"
	_ "fmt"
	. "hello/pkg/core/model"
	S "hello/pkg/core/structure"
	F "hello/pkg/util"
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
	t.Log("actualHash: ", actualHash)

	if actualHash != expectedHash {
		t.Errorf("GetHash() = %v; want %v", actualHash, expectedHash)
	}
}

func TestSignedTransaction_Ser(t *testing.T) {
	data := S.NewOrderedMap()
	data.Set("type", "Send")
	data.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	data.Set("amount", 3142500000)
	data.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	data.Set("timestamp", int64(1730603025159000))

	tx := NewSignedTransaction(data)

	expectedJSON := `{"type":"Send","to":"50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf","amount":3142500000,"from":"a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4","timestamp":1730603025159000}`
	actualJSON := tx.Ser()

	if actualJSON != expectedJSON {
		t.Errorf("Ser() = %v; want %v", actualJSON, expectedJSON)
	}

	expectedHash := "0625f96a8fd358fbb536f1ee332a8470fadc5c049af6835552a83a566a4d811fcb1f459cf50187"
	actualHash := tx.GetTxHash()

	timehex := F.HexTime(int64(1730603025159000))

	if timehex != "0625f96a8fd358" {
		t.Errorf("GetTxHash() = %v; want %v", timehex, "0625f96a8fd358")
	}

	if actualHash != expectedHash {
		t.Errorf("GetTxHash() = %v; want %v", actualHash, expectedHash)
	}

	// txMap := make(map[string]SignedTransaction, 1)
	// updateMap := make(map[string]Update, 3)

	update0 := Update{
		Key: "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770da53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4",
		Old: "99999999999999993566250000",
		New: "99999999999999990349375000",
	}

	if update0.GetHash() != "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770da53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4b139da1fb6537bdca92b9533dc005d517f9e99a955283887f5496d2f3a246701" {
		t.Errorf("GetHash() = %v; want %v", update0.GetHash(), "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770da53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4b139da1fb6537bdca92b9533dc005d517f9e99a955283887f5496d2f3a246701")
	}

	update1 := Update{
		Key: "c8c603ff91a3c59d637c7bda83e732dea6ec74e1001b35600f0ba7831dbfe32900000000000000000000000000000000000000000000",
		Old: "148750000",
		New: "223125000",
	}

	euh1 := "c8c603ff91a3c59d637c7bda83e732dea6ec74e1001b35600f0ba7831dbfe32900000000000000000000000000000000000000000000343ea3e5441ea4f6506426204bee366c89cbff5a86cd35618d56d10e2e141cae"
	if update1.GetHash() != euh1 {
		t.Errorf("GetHash() = %v; want %v", update1.GetHash(), euh1)
	}

	update2 := Update{
		Key: "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf",
		Old: "6285000000",
		New: "9427500000",
	}

	if update2.SerUpdateLog() != `{"old":"6285000000","new":"9427500000"}` {
		t.Errorf("SerUpdateLog() = %v; want %v", update2.SerUpdateLog(), `{"old":"6285000000","new":"9427500000"}`)
	}

	euh := "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf37051c0ed209c0909f567fab50dda1dd3268b3f4e9a5f5c29eaf64de268a0884"
	if update2.GetHash() != euh {
		t.Errorf("GetHash() = %v; want %v", update2.GetHash(), euh)
	}

	update3 := Update{
		Key: "724d2935080d38850e49b74927eb0351146c9ee955731f4ef53f24366c5eb9b100000000000000000000000000000000000000000000",
		Old: 4,
		New: 5,
	}

	blockPreviousBlockhash := "0625f96a9ca880efae3b7b47dc7ba9410ff36176096e9dfd321ca5e565cffaa4e908fcabcca389"
	block4 := NewBlock(4, blockPreviousBlockhash)

	block4.AppendUniversalUpdate(update0)
	block4.AppendUniversalUpdate(update1)
	block4.AppendUniversalUpdate(update2)

	block4.AppendLocalUpdate(update3)
	uhashs := block4.UHashs()
	j, _ := json.Marshal(uhashs)
	t.Logf("UHashs: %s", string(j))

	ur := block4.UpdateRoot()
	expectedUpdateRoot := "68af6d7009201e21283a75b345739ccea7c821ce6a0bc4fab105c8038ba9dd09"

	if ur != expectedUpdateRoot {
		t.Errorf("GetUpdateRoot() = %v; want %v", ur, expectedUpdateRoot)
	}
}
