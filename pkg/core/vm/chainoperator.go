package vm

import (
	"hello/pkg/core/model"
	"hello/pkg/core/storage"
	"math"
)

func OpGetBlock(i *Interpreter, vars interface{}) interface{} {
	target, objType := Unpack1Or2(vars)
	cs := storage.GetChainStorageInstance()

	var targetNum int

	switch v := target.(type) {
	case int:
		targetNum = int(v)
	case int64:
		targetNum = int(v)
	}

	if targetNum < 1 || targetNum > cs.GetLastHeight() {
		return nil
	}

	var block *model.Block
	if target != nil {
		block, _ = cs.GetBlock(targetNum)
	} else {
		block, _ = cs.GetLastBlock()
	}

	return block.Ser(objType.(string))
}

func OpListBlock(i *Interpreter, vars interface{}) interface{} {
	page, count, responseType := Unpack3(vars)
	blocks := make(map[string]interface{})
	cs := storage.GetChainStorageInstance()
	lastHeight := cs.GetLastHeight()

	var pageNum, countNum int

	switch v := page.(type) {
	case int:
		pageNum = v
	case int64:
		pageNum = int(v)
	}

	switch v := count.(type) {
	case int:
		countNum = v
	case int64:
		countNum = int(v)
	}

	max := int(math.Min(float64(pageNum*countNum), float64(lastHeight)))
	min := int(math.Max(float64(max-countNum)+1, 1))

	for i := min; i <= int(max); i++ {
		block, _ := cs.GetBlock(int(i))
		blocks[block.BlockHash()] = block.Ser(responseType.(string))
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
