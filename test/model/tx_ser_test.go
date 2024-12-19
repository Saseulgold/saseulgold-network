package model

import (
	"fmt"
	. "hello/pkg/core/model"
	S "hello/pkg/core/structure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignedTransactionSerialization(t *testing.T) {
	// Comment: Test setup - create a sample transaction
	txData := S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "867b8991f4f2eb94398d4647f5ddc57b30f8cb36acdf")
	txData.Set("amount", "100000000000")
	txData.Set("from", "4a8fd2ebb308370a689c3ef47cb83ba182683def3d4f")
	txData.Set("timestamp", int64(1734143956554000))

	privateKey := "dd97b057aa5d0fcc01acd23bdde9243dc22ec93110440c36800623b70c1c78c3"
	publicKey := "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2"

	// Comment: Create signed transaction
	signedTx, err := FromRawData(txData, privateKey, publicKey)
	assert.NoError(t, err)

	// Comment: Test serialization
	serialized, err := signedTx.Ser()
	assert.NoError(t, err)
	assert.NotEmpty(t, serialized)

	// Comment: Test deserialization
	fmt.Println("serialized: ", serialized)
	deserializedTx, err := ParseEscaped(serialized)
	assert.NoError(t, err)

	// Comment: Verify deserialized data matches original
	assert.Equal(t, signedTx.Xpub, deserializedTx.Xpub)
	assert.Equal(t, signedTx.Signature, deserializedTx.Signature)

	// Comment: Verify transaction data
	originalTx := txData.Ser()
	deserializedTxData, _ := deserializedTx.Data.Get("transaction")
	deserializedTxStr := deserializedTxData.(*S.OrderedMap).Ser()

	fmt.Println(originalTx)
	fmt.Println(deserializedTxStr)
}

func aTestSignedTransactionSerializationErrors(t *testing.T) {
	// Comment: Test missing required fields
	t.Run("Missing public key", func(t *testing.T) {
		tx := SignedTransaction{
			Data: S.NewOrderedMap(),
		}
		_, err := tx.Ser()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "public_key parameter")
	})

	t.Run("Missing signature", func(t *testing.T) {
		tx := SignedTransaction{
			Data: S.NewOrderedMap(),
			Xpub: "test_pub_key",
		}
		_, err := tx.Ser()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature parameter")
	})

	t.Run("Missing transaction", func(t *testing.T) {
		tx := SignedTransaction{
			Data:      S.NewOrderedMap(),
			Xpub:      "test_pub_key",
			Signature: "test_signature",
		}
		_, err := tx.Ser()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction parameter")
	})
}

func TestSignedTransactionParseEscaped(t *testing.T) {
	// Set up test data
	txData := S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "867b8991f4f2eb94398d4647f5ddc57b30f8cb36acdf")
	txData.Set("amount", "100000000000")
	txData.Set("from", "4a8fd2ebb308370a689c3ef47cb83ba182683def3d4f")
	txData.Set("timestamp", int64(1734143956554000))

	data := S.NewOrderedMap()
	data.Set("transaction", txData)
	data.Set("public_key", "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2")
	data.Set("signature", "test_signature")

	// Serialize
	serialized := data.Ser()

	// Test ParseEscaped
	t.Run("Basic parsing test", func(t *testing.T) {
		deserializedTx, err := ParseEscaped(serialized)
		assert.NoError(t, err)
		assert.NotNil(t, deserializedTx)

		// Verify basic fields
		assert.Equal(t, "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2", deserializedTx.Xpub)
		assert.Equal(t, "test_signature", deserializedTx.Signature)

		// Verify transaction data
		txMap, ok := deserializedTx.Data.Get("transaction")
		assert.True(t, ok)
		assert.NotNil(t, txMap)

		transaction := txMap.(*S.OrderedMap)
		txType, _ := transaction.Get("type")
		assert.Equal(t, "Send", txType)

		amount, _ := transaction.Get("amount")
		assert.Equal(t, "100000000000", amount)

		timestamp, _ := transaction.Get("timestamp")
		assert.Equal(t, int64(1734143956554000), timestamp)
	})

	t.Run("Invalid format data test", func(t *testing.T) {
		_, err := ParseEscaped("invalid json data")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse serialized data")
	})

	t.Run("Missing transaction field test", func(t *testing.T) {
		invalidData := S.NewOrderedMap()
		invalidData.Set("public_key", "test_pub_key")
		invalidData.Set("signature", "test_signature")

		serialized := invalidData.Ser()
		_, err := ParseEscaped(serialized)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing transaction in serialized data")
	})
}
