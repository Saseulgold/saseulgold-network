package model

import (
	"fmt"
	"hello/pkg/core/storage"
	"hello/pkg/core/structure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockDeserialization(t *testing.T) {
	// Test block JSON string
	blockJson := `{"height":1,"s_timestamp":1734319679804410,"previous_blockhash":"","blockhash":"06295ac47113fa512b78ff5782b8a9afaa80ab7ce907b458eb15ce4b9f213c6b9b5ca7dad4562f","difficulty":0,"reward_address":"","vout":"","nonce":"","transactions":{"06295ac47102c0c62f7fd925aca23697532dd27ddb9e30cefa3f15f198a29b906df3bb49434664":{"transaction":{"type":"Genesis","timestamp":1734319679800000},"signature":"3dca866aab17ab9ec55597ce7efe526857a4d8cc3fb23f90014983ee9bbd9c5a76499a858b19b88c84acae4c37fcf31a1c4f786a6863ad2fdeea9098f28d9b07","public_key":"391e87c9ceedb34ecd7f74d4536a33851ce54dbb0c2dfbf1a529816f8ed78afd"}},"universal_updates":{"02a35620a542dd255bc1d258ae935bdd4a05b479001b8a4ca630d214b1dbd21700":{"old":null,"new":true}},"local_updates":{}}`
	om, err := structure.ParseOrderedMap(blockJson)
	assert.NoError(t, err)
	fmt.Println("om: ", om.Ser())

	// Parse block
	block, err := storage.ParseBlock([]byte(blockJson))
	assert.NoError(t, err)

	// Validate basic fields
	assert.Equal(t, 1, block.Height, "Block height validation")
	assert.Equal(t, int64(1734319679804410), block.Timestamp_s, "Timestamp validation")
	assert.Equal(t, "", block.PreviousBlockhash, "Previous block hash validation")
	assert.Equal(t, 0, block.Difficulty, "Difficulty validation")
	assert.Equal(t, "", block.RewardAddress, "Reward address validation")
	assert.Equal(t, "", block.Vout, "Vout validation")
	assert.Equal(t, "", block.Nonce, "Nonce validation")

	// Validate transactions
	expectedTxID := "06295ac47102c0c62f7fd925aca23697532dd27ddb9e30cefa3f15f198a29b906df3bb49434664"
	assert.Equal(t, 1, len(*block.Transactions), "Number of transactions")
	// assert.Contains((*block.Transactions), expectedTxID, "Transaction ID exists")

	tx := (*block.Transactions)[expectedTxID]
	assert.Equal(t, "Genesis", tx.GetType(), "Transaction type")
	assert.Equal(t, int64(1734319679800000), tx.GetTimestamp(), "Transaction timestamp")
	assert.Equal(t, "3dca866aab17ab9ec55597ce7efe526857a4d8cc3fb23f90014983ee9bbd9c5a76499a858b19b88c84acae4c37fcf31a1c4f786a6863ad2fdeea9098f28d9b07", tx.GetSignature(), "Transaction signature")
	assert.Equal(t, "391e87c9ceedb34ecd7f74d4536a33851ce54dbb0c2dfbf1a529816f8ed78afd", tx.GetXpub(), "Transaction public key")

	// Validate universal updates
	expectedUpdateKey := "02a35620a542dd255bc1d258ae935bdd4a05b479001b8a4ca630d214b1dbd21700"
	assert.Equal(t, 1, len(*block.UniversalUpdates), "Number of universal updates")
	assert.Contains(t, *block.UniversalUpdates, expectedUpdateKey, "Universal update key exists")

	update := (*block.UniversalUpdates)[expectedUpdateKey]
	assert.Nil(t, update.Old, "Universal update old value")
	assert.Equal(t, true, update.New, "Universal update new value")

	// Validate local updates
	assert.Equal(t, 0, len(*block.LocalUpdates), "Number of local updates")
}
