package main

import (
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	S "hello/pkg/core/storage"
	"testing"
)

func TestBundle(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "storagestatus_test0"

	t.Log("Starting NewStatusFile test")
	sf := S.GetStatusFileInstance()
	t.Log("Created new status file")

	err := sf.Cache()

	if err != nil {
		t.Errorf("Error occurred during Cache(): %v", err)
	}

	t.Log("Successfully executed Cache operation")

	if sf.CachedUniversalIndexes == nil {
		t.Error("CachedUniversalIndexes가 초기화되지 않았습니다")
	}

	balancePrefix := "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770d"

	for key, cursor := range sf.CachedUniversalIndexes {
		prefix := key[:64]
		surfix := key[64:]
		if prefix == balancePrefix {
			value, err := S.ReadPart(sf.UniversalBundle("0000"), cursor.Seek, int(cursor.Length))

			if err != nil {
				t.Errorf("값 읽기 실패: %v", err)
				continue
			}
			DebugLog(prefix, surfix, value)
		}
	}
}
