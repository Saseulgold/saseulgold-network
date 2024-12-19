package main

import (
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"

	// . "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	F "hello/pkg/util"
	"testing"
)

func TestFileIdBin(t *testing.T) {
	f := FileIdBin("1000")
	// BinToFileId 테스트
	fileId := BinToFileId(f)
	DebugLog(fmt.Sprintf("BinToFileId 결과: %s", fileId))

	if fileId != "1000" {
		t.Errorf("BinToFileId 변환 실패: 예상값 '1000', 실제값 '%s'", fileId)
	}
	DebugLog(fmt.Sprintf("FileIdBin: %v", F.Bin2Hex(f))) // %v로 실제 바이트 값을 출력

	if len(f) != C.DATA_ID_BYTES && f[0] != 0 {
		t.Errorf("FileIdBin length error: %d", len(f))
	}

	f = FileIdBin("0100")
	DebugLog(fmt.Sprintf("FileIdBin: %v", F.Bin2Hex(f)))

	if len(f) != C.DATA_ID_BYTES {
		t.Errorf("FileIdBin length error: %d", len(f))
	}
	sf := GetStatusFileInstance()
	// StatusFile 인스턴스 상태 확인
	DebugLog(fmt.Sprintf("StatusFile 인스턴스: %+v", sf))

	NextFileId := sf.NextFileID("0000")
	DebugLog(fmt.Sprintf("NextFileId: %v", NextFileId))

	h0 := F.Hex2UInt64("0000")
	DebugLog(fmt.Sprintf("h0: %v", h0))

	h1 := F.Hex2UInt64("0010")
	DebugLog(fmt.Sprintf("h1: %v", h1))

	h2 := F.Hex2UInt64("0000") + 1
	DebugLog(fmt.Sprintf("h2: %v", h2))
}

func TestIndexRawAndParse(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "index_test"

	// 테스트 데이터 준비
	// asfsdafasdfd""
	expectedKey := "a12a66e790ac21dd95aa37bb6aaa22e8ac19559ed88297f127f51772513cb3910000000000000000000000000000000000000000000000000000000000000000"
	expectedFileID := "0000"
	expectedSeek := int64(0)
	expectedLength := int64(3)

	// indexRaw 함수 테스트
	indexData := IndexRaw(expectedKey, expectedFileID, expectedSeek, expectedLength)
	DebugLog(fmt.Sprintf("Generated index data: %v (length: %d bytes)", indexData, len(indexData)))
	DebugLog(fmt.Sprintf("expectedKey 길이: %d", len(expectedKey)))

	// ParseIndexRaw 함수로 데이터 파싱
	key, fileID, seek, length, err := ParseIndexRaw(indexData)
	if err != nil {
		t.Errorf("parseIndexRaw 에러: %v", err)
	}

	DebugLog(fmt.Sprintf("파싱된 key 길이: %d", len(key)))
	DebugLog(fmt.Sprintf("expectedKey: %s", expectedKey))
	DebugLog(fmt.Sprintf("파싱된 key: %s", key))

	// 결과 검증
	if key != expectedKey {
		t.Errorf("Key Error key = %v; want %v", key, expectedKey)
	}

	if fileID != expectedFileID {
		t.Errorf("fileID = %v; want %v", fileID, expectedFileID)
	}

	if seek != expectedSeek {
		t.Errorf("seek = %v; want %v", seek, expectedSeek)
	}

	if length != expectedLength {
		t.Errorf("length = %v; want %v", length, expectedLength)
	}

	// 디버그 로그 출력
	DebugLog(fmt.Sprintf("\n=== 인덱스 파싱 결과 ==="))
	DebugLog(fmt.Sprintf("Key: %s", key))
	DebugLog(fmt.Sprintf("FileID: %s", fileID))
	DebugLog(fmt.Sprintf("Seek: %d", seek))
	DebugLog(fmt.Sprintf("Length: %d", length))

	/**
	genesisUpdate := Update{
		Key: F.StatusHash(C.ZERO_ADDRESS, F.RootSpace(), "thisisgenesis", C.ZERO_ADDRESS), // 길이를 108로 맞추기 위해 해시 추가
		Old: nil,
		New: "1",
	}
	**/

	// sf := GetStatusFileInstance()
	// sf.WriteUniversal(map[string]Update{genesisUpdate.Key: genesisUpdate})
}
