package main

import (
	. "hello/pkg/util"
	F "hello/pkg/util"
	"os"
	"testing"

	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	S "hello/pkg/core/storage"
	. "hello/pkg/core/structure"

	"fmt"
	"strings"
)

func TestNewStatusFile(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "storagestatus_test0"

	t.Log("Starting NewStatusFile test")
	sf := S.GetStatusFileInstance()

	if sf.CachedUniversalIndexes == nil {
		t.Error("CachedUniversalIndexes was not initialized")
	}
	t.Log("Verified CachedUniversalIndexes initialization")

	if sf.CachedLocalIndexes == nil {
		t.Error("CachedLocalIndexes was not initialized")
	}
	t.Log("Verified CachedLocalIndexes initialization")

	if sf.Tasks == nil {
		t.Error("Tasks was not initialized")
	}
	t.Log("Verified Tasks initialization")
}

func TestStatusFile_Touch(t *testing.T) {
	t.Log("Starting Touch test")
	C.CORE_TEST_MODE = true

	sf := S.GetStatusFileInstance()
	t.Log("Created new status file")

	err := sf.Touch()
	if err != nil {
		t.Errorf("Error occurred during Touch(): %v", err)
	}
	t.Log("Successfully executed Touch operation")

	t.Log("Verifying created files...")
	expectedFiles := []string{
		sf.TempFile(),
		sf.InfoFile(),
		sf.LocalFile(),
		sf.LocalBundle(),
		sf.UniversalBundle("0000"),
		sf.LocalBundleIndex(),
		sf.UniversalBundleIndex(),
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected file was not created: %s", file)
		}
		t.Logf("Verified file exists: %s", file)
	}
}

func TestStatusFile_Reset(t *testing.T) {
	t.Log("Starting Reset test")
	C.CORE_TEST_MODE = true

	sf := S.GetStatusFileInstance()
	t.Log("Created new status file")

	t.Log("Initializing with Touch operation")
	err := sf.Touch()
	if err != nil {
		t.Errorf("Error occurred during Touch(): %v", err)
	}

	t.Log("Executing Reset operation")
	err = sf.Reset()
	if err != nil {
		t.Errorf("Error occurred during Reset(): %v", err)
	}

	t.Log("Verifying StatusBundle recreation")
	if _, err := os.Stat(sf.StatusBundle()); os.IsNotExist(err) {
		t.Error("StatusBundle was not recreated after Reset")
	}
	t.Log("Reset test completed successfully")
}

func TestStatusFile_Cache(t *testing.T) {
	t.Log("Starting Cache test")
	C.CORE_TEST_MODE = true

	sf := S.GetStatusFileInstance()
	t.Log("Created new status file")

	err := sf.Cache()
	if err != nil {
		t.Errorf("Error occurred during Cache(): %v", err)
	}
	t.Log("Cache operation completed successfully")
}

func TestMainUpdate(t *testing.T) {
	AddressFromInt64 := func(n int64) string {
		// 정수를 16진수 문자열로 변환
		hex := fmt.Sprintf("%x", n)

		// 44자리가 되도록 앞에 0으로 패딩
		padded := strings.Repeat("0", 44-len(hex)) + hex

		return padded
	}

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

	// 50개의 Update 객체 생성 및 추가
	for i := 0; i < 8; i++ {
		update := Update{
			Key: StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", AddressFromInt64(int64(i))),
			Old: nil,
			New: fmt.Sprintf("%d000000000000000000000000", 100+i),
		}
		block.AppendUniversalUpdate(update)
		t.Logf("Generated Update Key[%d]: %s", i, update.Key)
	}

	sf.Write(&block)
	cursors := S.ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)

	for key, cursor := range cursors {
		prefix := key[:64]
		suffix := key[64:]

		addr := suffix[:44]
		part, err := S.ReadPart(sf.UniversalBundle("0000"), cursor.Seek, int(cursor.Length))
		if err != nil {
			t.Errorf("Error occurred during ReadPart(): %v", err)
		}
		if false {
			t.Logf("prefix: %s, address: %s, value: %s, length: %d", prefix, addr, part, cursor.Length)
		}
		// Generate expected value

	}
	// sf.Update(&block)
}

func aaTestStorageStatusOverwrite(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "storage_status_overwrite_test"
	sf := S.GetStatusFileInstance()

	// First block
	block1 := NewBlock(0, "")
	block1.Init()

	// First update - base length
	update1 := Update{
		Key: StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", C.ZERO_ADDRESS),
		Old: nil,
		New: "1000000000000000000000000",
	}
	block1.AppendUniversalUpdate(update1)
	sf.Write(&block1)

	// Check first data
	cursors1 := S.ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)
	for key, cursor := range cursors1 {
		part, err := S.ReadPart(sf.UniversalBundle("0000"), cursor.Seek, int(cursor.Length))
		if err != nil {
			t.Errorf("First read error: %v", err)
		}
		t.Logf("First write - Key: %s, Value: %s, Length: %d", key, part, cursor.Length)
	}

	// Second block - longer data
	block2 := NewBlock(1, "")
	block2.Init()

	update2 := Update{
		Key: StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", C.ZERO_ADDRESS),
		Old: "1000000000000000000000000",
		New: "100000000000000000000000000000",
	}
	block2.AppendUniversalUpdate(update2)
	sf.Write(&block2)

	// Check second data
	cursors2 := S.ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)
	for key, cursor := range cursors2 {
		part, err := S.ReadPart(sf.UniversalBundle("0000"), cursor.Seek, int(cursor.Length))
		if err != nil {
			t.Errorf("Second read error: %v", err)
		}
		t.Logf("Second write - Key: %s, Value: %s, Length: %d", key, part, cursor.Length)
	}

	// Third block - shorter data
	block3 := NewBlock(2, "")
	block3.Init()

	update3 := Update{
		Key: StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", C.ZERO_ADDRESS),
		Old: "100000000000000000000000000000",
		New: "1000",
	}
	block3.AppendUniversalUpdate(update3)
	sf.Write(&block3)

	// Check third data
	cursors3 := S.ReadStatusStorageIndex(sf.UniversalBundleIndex(), true)
	for key, cursor := range cursors3 {
		part, err := S.ReadPart(sf.UniversalBundle("0000"), cursor.Seek, int(cursor.Length))
		if err != nil {
			t.Errorf("Third read error: %v", err)
		}
		t.Logf("Third write - Key: %s, Value: %s, Length: %d", key, part, cursor.Length)
	}
}
