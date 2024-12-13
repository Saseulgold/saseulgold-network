package main

import (
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	S "hello/pkg/core/structure"
	. "hello/pkg/crypto"
	"testing"
)

func aaTestSignedTransaction_WithRealData(t *testing.T) {
	C.CORE_TEST_MODE = true

	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "867b8991f4f2eb94398d4647f5ddc57b30f8cb36acdf")
	txData.Set("amount", "10000000000")
	txData.Set("from", "4a8fd2ebb308370a689c3ef47cb83ba182683def3d4f")
	txData.Set("timestamp", int64(1732359603011000))

	data.Set("transaction", txData)
	data.Set("public_key", "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2")
	data.Set("signature", "b1adf108db92fe4e062bd7e79c4d48b8202267f2968740a59fbdaf66f6f826768d7fb740cec4075dbe71abf4143aed9fa901e8b90b357ded06a5647ecb74da0c")

	tx, err := NewSignedTransaction(data)

	if err != nil {
		t.Errorf("NewSignedTransaction(): %s", err)
	}

	hash := tx.GetTxHash()
	DebugLog("hash: " + hash)
	if hash != "06279266c2bdb852cc9ec8e60fcbbef15442fa88bb398728b9e7797fdfa0cd8878e5373dbc7089" {
		t.Errorf("tx.GetTxHash() failed\nExpected: %s\nActual: %s", "06279266c2bdb852cc9ec8e60fcbbef15442fa88bb398728b9e7797fdfa0cd8878e5373dbc7089", hash)
	}

	privateKey := "dd97b057aa5d0fcc01acd23bdde9243dc22ec93110440c36800623b70c1c78c3"
	signature := Signature(hash, privateKey)
	DebugLog("signature: " + signature)

	if signature != "b1adf108db92fe4e062bd7e79c4d48b8202267f2968740a59fbdaf66f6f826768d7fb740cec4075dbe71abf4143aed9fa901e8b90b357ded06a5647ecb74da0c" {
		t.Errorf("Signature() failed\nExpected: %s\nActual: %s", "b1adf108db92fe4e062bd7e79c4d48b8202267f2968740a59fbdaf66f6f826768d7fb740cec4075dbe71abf4143aed9fa901e8b90b357ded06a5647ecb74da0c", signature)
	}

	if err != nil {
		t.Errorf("tx: %s", err)
	}

	// 검증
	errMsg := tx.Validate()
	if errMsg != nil {
		t.Errorf("tx.Validate(): %s", errMsg)
	}

	// 서명 검증
	if tx.Signature == "" {
		t.Error("서명이 없습니다")
	}

	if tx.Xpub == "" {
		t.Error("공개키가 없습니다")
	}
}

func TestSignedTransaction_WithRealData2(t *testing.T) {
	C.CORE_TEST_MODE = true

	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "e9696a5f58d43772e87a37f21bae4b2eee2f2a750d1c")
	txData.Set("amount", "1000000000000000000")
	txData.Set("from", "43da67e738d53473cf4fd307f0acb534b72c62a806c4")
	txData.Set("timestamp", int64(1733211394654000))

	data.Set("transaction", txData)
	data.Set("public_key", "8860ecfe5711c9096f43411ce1ebefcb292200fbca73aa14fbf187a52cc29898")
	data.Set("signature", "1493bd19ea174751810b3fece0f23fa24c9e6d884118624e863b8fc3892f5604dba420407e2a833440d58a2e3c3e1048109f4f14f5f9fd3c9788a86e0bc5f400")

	tx, err := NewSignedTransaction(data)
	if err != nil {
		t.Errorf("NewSignedTransaction(): %s", err)
	}

	// 검증
	errMsg := tx.Validate()
	if errMsg != nil {
		t.Errorf("tx.Validate(): %s", errMsg)
	}

	// 서명 검증
	if tx.Signature == "" {
		t.Error("서명이 없습니다")
	}

	if tx.Xpub == "" {
		t.Error("공개키가 없습니다")
	}
}
