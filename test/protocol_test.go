package main

import (
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/storage"
	F "hello/pkg/util"
	"testing"
)

func TestFileIdBin(t *testing.T) {
	f := FileIdBin("1000")
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
