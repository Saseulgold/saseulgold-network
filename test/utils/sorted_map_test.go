package utils_test

import (
	"hello/pkg/core/model"
	"hello/pkg/core/structure"
	"hello/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedValueK(t *testing.T) {
	// comment: Create test data with SignedTransactions
	testMap := make(map[string]*model.SignedTransaction)

	// comment: Create dummy OrderedMap data for each transaction
	data1 := structure.NewOrderedMap()
	data1.Set("name", "Charlie")
	data1.Set("amount", 300)

	data2 := structure.NewOrderedMap()
	data2.Set("name", "Alice")
	data2.Set("amount", 100)

	data3 := structure.NewOrderedMap()
	data3.Set("name", "Bob")
	data3.Set("amount", 200)

	// comment: Create sample transactions with different hashes and data
	tx1 := &model.SignedTransaction{
		Data: data1,
	}
	tx2 := &model.SignedTransaction{
		Data: data2,
	}
	tx3 := &model.SignedTransaction{
		Data: data3,
	}

	testMap["c"] = tx1
	testMap["a"] = tx2
	testMap["b"] = tx3

	// comment: Get sorted transactions
	sorted := util.SortedValueK(testMap)

	// comment: Verify length is preserved
	assert.Equal(t, 3, len(sorted))

	// comment: Verify data content is preserved
	amount1, _ := sorted[0].Data.Get("amount")
	assert.Equal(t, 100, amount1)
	name1, _ := sorted[0].Data.Get("name")
	assert.Equal(t, "Alice", name1)

	amount2, _ := sorted[1].Data.Get("amount")
	assert.Equal(t, 200, amount2)
	name2, _ := sorted[1].Data.Get("name")
	assert.Equal(t, "Bob", name2)

	amount3, _ := sorted[2].Data.Get("amount")
	assert.Equal(t, 300, amount3)
	name3, _ := sorted[2].Data.Get("name")
	assert.Equal(t, "Charlie", name3)
}
