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

		// Preparing test data
		testKey := StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "balance", C.ZERO_ADDRESS)
		update := Update{
			Key: testKey,
			Old: nil,
			New: "1000000000000000000000000",
		}

		// Perform status updates
		err := sf.Cache()
		if err != nil {
			t.Fatal(err)
		}

		updates := &map[string]Update{}
		(*updates)[update.Key] = update

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

		// Check the index
		cursor, exists := sf.CachedUniversalIndexes[testKey]
		if !exists {
			t.Error("Index not found")
		}

		// data verification
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

		// First Update
		update1 := Update{
			Key: testKey,
			Old: nil,
			New: "1000000000000000000000000",
		}

		updates := &map[string]Update{}
		(*updates)[update1.Key] = update1

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
			// Second update (larger data)
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
