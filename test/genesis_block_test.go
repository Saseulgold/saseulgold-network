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

	// Universal Updates 정의
	balanceUpdate := Update{
		Key: F.StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", C.ZERO_ADDRESS),
		Old: nil,
		New: "10000000000000000000000000",
	}

	genesisUpdate := Update{
		Key: F.StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "thisisgenesis", C.ZERO_ADDRESS),
		Old: nil,
		New: "1",
	}

	t.Logf("balanceUpdate Key: %s", balanceUpdate.Key)
	t.Logf("genesisUpdate Key: %s", genesisUpdate.Key)
	// Universal Updates 추가
	// block.AppendUniversalUpdate(balanceUpdate)
	block.AppendUniversalUpdate(genesisUpdate)

	if len(balanceUpdate.Key) != 108 {
		t.Errorf("balanceUpdate.Key length = %d; want %d", len(balanceUpdate.Key), 108)
	}

	if len(genesisUpdate.Key) != 108 {
		t.Errorf("genesisUpdate.Key length = %d; want %d", len(genesisUpdate.Key), 108)
	}

	// block.AppendTransaction(tx)

	err := sf.Cache()
	if err != nil {
		t.Errorf("Cache Error: %v", err)
	}
	// 1t.Logf("CachedUniversalIndexes: %+v", sf.CachedUniversalIndexes)

	sf.Write(&block)
	sf.Update(&block)

	// height := S.LastHeight()
	// t.Logf("Last Height: %d", height)
	universalIndexes := S.ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)
	for key := range universalIndexes {
		t.Logf("Universal Index Key: %s, Value: %s, FileID: %s", key, universalIndexes[key].Value, universalIndexes[key].FileID)
	}
}
