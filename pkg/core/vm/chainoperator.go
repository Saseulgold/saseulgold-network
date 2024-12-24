package vm

import (
	"hello/pkg/core/model"
	"hello/pkg/core/storage"
	"math"
)

func OpGetBlock(i *Interpreter, vars interface{}) interface{} {
	target, objType := Unpack1Or2(vars)
	cs := storage.GetChainStorageInstance()

	var block *model.Block
	if target != nil {
		block, _ = cs.GetBlock(int(target.(float64)))
	} else {
		block, _ = cs.GetLastBlock()
	}

	return block.Ser(objType.(string))
}

func OpListBlock(i *Interpreter, vars interface{}) interface{} {
	page, count, sort := Unpack3(vars)
	blocks := make(map[string]interface{})
	cs := storage.GetChainStorageInstance()

	pageNum := int(page.(int))
	countNum := int(count.(int))
	sortNum := int(sort.(int))

	if sortNum == 1 {
		max := pageNum * countNum
		min := int(math.Max(float64(max-countNum), 1))

		for i := min; i <= int(max); i++ {
			block, _ := cs.GetBlock(int(i))
			blocks[block.BlockHash()] = block.Ser("full")
		}
	} else {
		min := cs.GetLastHeight() - (pageNum * countNum)
		max := min + countNum

		for i := max; i > int(math.Max(float64(min), 0)); i-- {
			block, _ := cs.GetBlock(int(i))
			blocks[block.BlockHash()] = block.Ser("full")
		}
	}

	return blocks
}

func OpBlockCount(i *Interpreter, vars interface{}) interface{} {
	cs := storage.GetChainStorageInstance()
	return cs.GetLastHeight()
}

func OpListTransaction(i *Interpreter, vars interface{}) interface{} {
	count := 20
	if arr, ok := vars.([]interface{}); ok && len(arr) > 0 {
		if c, ok := arr[0].(float64); ok {
			count = int(c)
		}
	}

	transactions := make(map[string]interface{})
	cs := storage.GetChainStorageInstance()
	lastHeight := cs.GetLastHeight()

	for height := lastHeight; height > 0; height-- {
		block, err := cs.GetBlock(height)
		if err != nil {
			continue
		}

		txs := block.GetTransactions()
		for hash, tx := range txs {
			transactions[hash] = tx

			if len(transactions) >= count {
				OperatorLog("OpListTransaction", "input:", vars, "result:", transactions)
				return transactions
			}
		}
	}

	OperatorLog("OpListTransaction", "input:", vars, "result:", transactions)
	return transactions
}
