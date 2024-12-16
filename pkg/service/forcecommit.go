package service

import (
	"fmt"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
	. "hello/pkg/core/vm"
	"hello/pkg/util"
)

func ForceCommit(txs map[string]*SignedTransaction) error {
	oracle := GetOracleService()

	sf := oracle.storage
	ci := oracle.chain

	sf.Touch()
	ci.Touch()

	machine := GetMachineInstance()
	var previousBlockhash string

	lastBlockHeight := storage.LastHeight()
	previousBlock, err := storage.GetChainStorageInstance().GetBlock(lastBlockHeight)

	fmt.Println("previousBlock: ", previousBlock)

	if previousBlock == nil {
		previousBlockhash = ""
	} else {
		previousBlockhash = previousBlock.BlockHash()
	}

	machine.Init(previousBlock, int64(util.Utime()))
	machine.SetTransactions(txs)
	machine.PreCommit()

	universals := machine.GetInterpreter().GetUniversals()
	DebugLog("universals: %v", universals)

	block := NewBlock(storage.LastHeight()+1, previousBlockhash)
	block.SetTimestamp(int64(util.Utime()))

	expectedBlock := machine.ExpectedBlock()

	if expectedBlock.GetTransactionCount() == 0 {
		DebugLog("no transactions to commit. invalid block.")
		return fmt.Errorf("no transactions to commit. invalid block.")
	}

	DebugLog("unv: %v", expectedBlock.UniversalUpdates)
	DebugLog("loc: %v", expectedBlock.LocalUpdates)

	err = machine.Commit(expectedBlock)
	if err != nil {
		DebugLog("Commit error: %v", err)
		return fmt.Errorf("Commit error: %v", err)
	}

	DebugLog("Commit success")
	return nil
}
