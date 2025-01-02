package model

import (
	"fmt"
	"hello/pkg/core/model"
	S "hello/pkg/core/structure"
	. "hello/pkg/core/vm"
	"hello/pkg/util"
	"testing"

	"hello/pkg/core/storage"

	"github.com/stretchr/testify/assert"
)

func TestSerializationDeserialization(t *testing.T) {

	// Comment: Initialize test environment
	machine := GetMachineInstance()
	machine.GetInterpreter().Reset(true)

	// Comment: Create test transaction data
	txData := createTestTransaction(t)
	tx, err := model.NewSignedTransaction(txData)
	assert.NoError(t, err)
	txHash := tx.GetTxHash()

	// Comment: Set up test updates with various data types
	updates := S.NewOrderedMap()
	updates.Set("string_value", "테스트 문자열")
	updates.Set("int_value", 12345)
	updates.Set("float_value", 123.456)
	updates.Set("boolean_value", true)
	updates.Set("array_value", []string{"a", "b", "c"})
	nestedMap := S.NewOrderedMap()
	nestedMap.Set("nested_key", "nested_value")
	updates.Set("map_value", nestedMap)

	// Comment: Create and initialize test block
	originalBlock := &model.Block{
		Height:            100,
		PreviousBlockhash: "previous_hash_test",
		Timestamp_s:       util.Utime(),
		Vout:              "vout_test",
		Nonce:             "nonce_test",
		RewardAddress:     "reward_address_test",
		Difficulty:        4,
		Transactions:      &map[string]*model.SignedTransaction{txHash: &tx},
		LocalUpdates:      &map[string]model.Update{txHash: {Key: txHash, Old: nil, New: nil}},
	}
	originalBlock.Init()

	txSer, err := (*originalBlock.Transactions)[txHash].Ser()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("block txs: ", txSer)
	s := originalBlock.Ser("full")
	fmt.Println("serialized: ", s)
	// Comment: Test Block Serialization/Deserialization
	t.Run("Test Block Serialization/Deserialization", func(t *testing.T) {
		serialized := originalBlock.Ser("full")
		fmt.Println(serialized)
		deserializedBlock, err := storage.ParseBlock([]byte(serialized))
		fmt.Println(err)
		fmt.Println(fmt.Sprintf("deserializedBlock.Transactions: %v", deserializedBlock.Transactions))

		assert.NoError(t, err)

		// Comment: Verify block fields
		assert.Equal(t, originalBlock.Height, deserializedBlock.Height, "블록 높이 검증")
		assert.Equal(t, originalBlock.PreviousBlockhash, deserializedBlock.PreviousBlockhash, "이전 블록 해시 검증")
		assert.Equal(t, originalBlock.Timestamp_s, deserializedBlock.Timestamp_s, "타임스탬프 검증")
		assert.Equal(t, originalBlock.Difficulty, deserializedBlock.Difficulty, "난이도 검증")
	})

}

// Comment: Test Transaction Serialization/Deserialization

func createTestTransaction(t *testing.T) *S.OrderedMap {
	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	txData.Set("from", "60c3a6cd858c90574bcdc35b2da5dbc7225275f50efd")
	txData.Set("amount", "1000")
	txData.Set("timestamp", util.Utime())

	data.Set("transaction", txData)
	data.Set("public_key", "test_public_key")
	data.Set("signature", "test_signature")

	return data
}
