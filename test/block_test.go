package main

import (
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
	t.Log(actualJSON)

	if actualJSON != expectedJSON {
		t.Errorf("Ser() = %v; want %v", actualJSON, expectedJSON)
	}

	expectedHash := "0625f96a8fd358fbb536f1ee332a8470fadc5c049af6835552a83a566a4d811fcb1f459cf50187"
	actualHash := tx.GetTxHash()

	timehex := F.HexTime(int64(1730603025159000))
	t.Log("timehex: ", timehex)

	if timehex != "0625f96a8fd358" {
		t.Errorf("GetTxHash() = %v; want %v", timehex, "0625f96a8fd358")
	}

	t.Log("actualHash: ", actualHash)
	if actualHash != expectedHash {
		t.Errorf("GetTxHash() = %v; want %v", actualHash, expectedHash)
	}
}
