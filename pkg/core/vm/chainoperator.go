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
	page, count := Unpack2(vars)
	blocks := make(map[string]interface{})
	cs := storage.GetChainStorageInstance()

	pageNum := int(page.(int))
	countNum := int(count.(int))

	max := pageNum * countNum
	min := int(math.Max(float64(max-countNum), 1))

	for i := min; i <= int(max); i++ {
		block, _ := cs.GetBlock(int(i))
		blocks[block.BlockHash()] = block.Ser("full")
	}

	return blocks
}

func OpBlockCount(i *Interpreter, vars interface{}) interface{} {
	cs := storage.GetChainStorageInstance()
	return cs.GetLastHeight()
}

func OpListTransaction(i *Interpreter, vars interface{}) interface{} {
	count := Unpack1(vars)
	countNum := int(count.(int))

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
			payload, _ := tx.Ser()
			transactions[hash] = payload

			if len(transactions) >= countNum {
				OperatorLog("OpListTransaction", "input:", vars, "result:", transactions)
				return transactions
			}
		}
	}

	OperatorLog("OpListTransaction", "input:", vars, "result:", transactions)
	return transactions
}
