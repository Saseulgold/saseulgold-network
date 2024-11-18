package vm

import (
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
)

type Machine struct {
	interpreter *Interpreter

	contracts           map[string]map[string]*Method
	requests            map[string]map[string]interface{}
	postProcessContract map[string]interface{}

	previousBlock  *Block
	roundTimestamp int64
	transactions   *map[string]*SignedTransaction

	miningNodeCount int

	currentBlockDifficulty string
	currentBlockVout       string
	currentBlockNonce      string
}

var instance *Machine

func GetInstance() *Machine {
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

func (m *Machine) Write(block *Block) {
	chain := storage.GetChainStorageInstance()
	sf := storage.GetStatusFileInstance()

	if _, err := chain.Write(block); err != nil {
		return err
	}
	if err := sf.Write(block); err != nil {
		return err
	}
	if err := sf.Update(block); err != nil {
		return err
	}

}
