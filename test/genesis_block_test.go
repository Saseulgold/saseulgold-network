package main

import (
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	S "hello/pkg/core/storage"
	. "hello/pkg/core/structure"
	F "hello/pkg/util"
	"testing"
)

func TestGenesisBlock(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "genesis_test"
	// sfi := S.GetStatusIndexInstance()
	sf := S.GetStatusFileInstance()

	// Genesis 블록 생성
	block := NewBlock(0, "")
	block.SetTimestamp(1731062860000000)

	txData := NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	txData.Set("amount", 3142500000)
	txData.Set("from", "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4")
	txData.Set("timestamp", int64(1731062859742000))

	//tx := NewSignedTransaction(txData)
	block.Init()

	// Universal Updates 추가
	block.AppendUniversalUpdate(Update{
		Key: F.StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", C.ZERO_ADDRESS),
		Old: nil,
		New: "10000000000000000000000000",
	})

	block.AppendUniversalUpdate(Update{
		Key: F.StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "thisisgenesis", C.ZERO_ADDRESS),
		Old: nil,
		New: "1",
	})

	// block.AppendTransaction(tx)

	sf.Cache()
	t.Logf("CachedUniversalIndexes: %+v", sf.CachedUniversalIndexes)

	sf.Write(&block)
	// sf.Update(&block)

	// height := S.LastHeight()
	// t.Logf("Last Height: %d", height)
	/**
	universalIndexes := S.ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)
	for key := range universalIndexes {
		t.Logf("Universal Index Key: %s, Value: %s, FileID: %s", key, universalIndexes[key].Value, universalIndexes[key].FileID)
	}
	/**
	localIndexes := S.ReadStatusStorageIndex(sf.LocalBundleIndex(), true)

	if len(localIndexes) != 0 {
		t.Errorf("Local Indexes 개수 = %d; want %d", len(localIndexes), 0)
	}

	if len(universalIndexes) != 2 {
		t.Errorf("Universal Indexes 개수 = %d; want %d", len(universalIndexes), 2)
	}
	**/

	// if height != 0 {
	// 	t.Errorf("Last Height = %d; want %d", height, 0)
	// }
}
