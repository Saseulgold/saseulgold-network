package main

import (
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	S "hello/pkg/core/structure"
	"testing"
	"time"
)

func TestSignedTransaction_GetSize(t *testing.T) {
	data := S.NewOrderedMap()
	data.Set("type", "test_tx")
	data.Set("timestamp", time.Now().Unix())

	tx := NewSignedTransaction(data)

	size := tx.GetSize()
	if size <= 0 {
		t.Errorf("트랜잭션 크기가 잘못되었습니다: got %v", size)
	}
}

func TestSignedTransaction_WithRealData(t *testing.T) {
	C.CORE_TEST_MODE = true
	// 실제 트랜잭션 데이터 생성

	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "867b8991f4f2eb94398d4647f5ddc57b30f8cb36acdf")
	txData.Set("amount", "10000000000")
	txData.Set("from", "4a8fd2ebb308370a689c3ef47cb83ba182683def3d4f")
	txData.Set("timestamp", int64(1732359603011000))
	data.Set("transaction", txData)

	tx := NewSignedTransaction(data)
	DebugLog("tx.Data.Ser(): " + tx.Data.Ser())

	tx.Xpub = "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2"
	tx.Signature = "b1adf108db92fe4e062bd7e79c4d48b8202267f2968740a59fbdaf66f6f826768d7fb740cec4075dbe71abf4143aed9fa901e8b90b357ded06a5647ecb74da0c"

	// 검증
	valid, errMsg := tx.Validate()
	if !valid {
		t.Errorf("실제 트랜잭션 데이터 검증 실패: %s", errMsg)
	}

	// 서명 검증
	if tx.Signature == "" {
		t.Error("서명이 없습니다")
	}

	if tx.Xpub == "" {
		t.Error("공개키가 없습니다")
	}
}
