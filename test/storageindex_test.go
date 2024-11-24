package main

import (
	"encoding/json"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	. "hello/pkg/util"
	F "hello/pkg/util"
	"testing"
)

func TestStorageIndex(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "storage_index_test"

	t.Run("Test index creation and reading", func(t *testing.T) {
		sf := GetStatusFileInstance()
		sf.Commit()

		// 테스트 데이터 준비
		testKey := StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", C.ZERO_ADDRESS)
		update := Update{
			Key: testKey,
			Old: nil,
			New: "1000000000000000000000000",
		}

		// 상태 업데이트 수행
		err := sf.Cache()
		if err != nil {
			t.Fatal(err)
		}

		updates := make(UpdateMap)
		updates[update.Key] = update

		err = sf.WriteUniversal(updates)
		if err != nil {
			t.Fatal(err)
		}

		err = sf.WriteTasks()
		if err != nil {
			t.Fatal(err)
		}

		err = sf.Commit()
		if err != nil {
			t.Fatal(err)
		}

		// 인덱스 확인
		cursor, exists := sf.CachedUniversalIndexes[testKey]
		if !exists {
			t.Error("Index not found")
		}

		// 데이터 검증
		data, err := ReadPart(sf.UniversalBundle(cursor.FileID), cursor.Seek, int(cursor.Length))
		if err != nil {
			t.Fatal(err)
		}

		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			t.Fatal(err)
		}
		if value != "1000000000000000000000000" {
			t.Errorf("Data mismatch: got %v, want %v", value, "1000000000000000000000000")
		}
	})

	t.Run("Test data size change scenario", func(t *testing.T) {
		sf := GetStatusFileInstance()
		sf.Reset()

		testKey := StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", C.ZERO_ADDRESS)

		// 첫 번째 업데이트
		update1 := Update{
			Key: testKey,
			Old: nil,
			New: "1000000000000000000000000",
		}

		updates := make(UpdateMap)
		updates[update1.Key] = update1

		err := sf.Cache()
		if err != nil {
			t.Fatal(err)
		}

		err = sf.WriteUniversal(updates)
		if err != nil {
			t.Fatal(err)
		}

		err = sf.WriteTasks()
		if err != nil {
			t.Fatal(err)
		}

		err = sf.Commit()
		if err != nil {
			t.Fatal(err)
		}

		// initialCursor := sf.CachedUniversalIndexes[testKey]
		// initialLength := initialCursor.Length

		/**
			// 두 번째 업데이트 (더 큰 데이터)
			update2 := Update{
				Key: testKey,
				Old: "1000000000000000000000000",
				New: "999999999999999999999999999999999",
			}

			updates[update2.Key] = update2

			err = sf.WriteUniversal(updates)
			if err != nil {
				t.Fatal(err)
			}

			err = sf.WriteTasks()
			if err != nil {
				t.Fatal(err)
			}

			err = sf.Commit()
			if err != nil {
				t.Fatal(err)
			}

			newCursor := sf.CachedUniversalIndexes[testKey]
			if newCursor.Length <= initialLength {
				t.Error("New data length should be larger")
			}
		**/
	})
}
