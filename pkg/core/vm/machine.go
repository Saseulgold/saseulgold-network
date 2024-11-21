package vm

import (
	"errors"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
	"hello/pkg/util"
)

type Machine struct {
	interpreter *Interpreter

	contracts           map[string]map[string]*Method
	requests            map[string]map[string]interface{}
	postProcessContract map[string]interface{}

	previousBlock  *Block
	roundTimestamp int64
	transactions   *map[string]*SignedTransaction

	currentBlockDifficulty string
	currentBlockVout       string
	currentBlockNonce      string
}

var instance *Machine

func GetMachineInstance() *Machine {
	if instance == nil {
		instance = &Machine{
			interpreter:    NewInterpreter(),
			previousBlock:  nil,
			roundTimestamp: 0,
		}
	}
	return instance
}

func (m *Machine) Init(previousBlock *Block, roundTimestamp int64) {
	m.interpreter.Reset()
	m.interpreter.Init("transaction")

	m.previousBlock = previousBlock
	m.roundTimestamp = roundTimestamp
	m.currentBlockDifficulty = ""
	m.currentBlockVout = ""
	m.currentBlockNonce = ""
}

func (m *Machine) ValidateBlockTimestamp(block *Block) error {
	DebugLog("ValidateBlockTimestamp", "block.Timestamp_s", block.Timestamp_s, "roundTimestamp", int(util.Utime())+C.TIME_STAMP_ERROR_LIMIT)

	if block.Timestamp_s > int(util.Utime())+C.TIME_STAMP_ERROR_LIMIT {
		return errors.New("block timestamp is greater than current round timestamp")
	}

	return nil
}

func (m *Machine) Commit(block *Block) error {
	/**
	if err := m.ValidateBlockTimestamp(block); err != nil {
		return err
	}
	**/

	chain := storage.GetChainStorageInstance()
	sf := storage.GetStatusFileInstance()

	if err := chain.Write(block); err != nil {
		return err
	}
	if err := sf.Write(block); err != nil {
		return err
	}
	/**
	if err := sf.Update(block); err != nil {
		return err
	}
	**/

	return nil
}

func (m *Machine) PreLoad(universalUpdates map[string]map[string]interface{}, localUpdates map[string]map[string]interface{}) {
	for key, update := range universalUpdates {
		old, exists := update["old"]
		if exists {
			m.interpreter.SetUniversalLoads(key, old)
		}
	}

	for key, update := range localUpdates {
		old, exists := update["old"]
		if exists {
			m.interpreter.SetLocalLoads(key, old)
		}
	}
}
