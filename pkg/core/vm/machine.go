package vm

import (
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
	"hello/pkg/rpc"
	"hello/pkg/util"
	F "hello/pkg/util"
	"slices"
)

type Machine struct {
	interpreter *Interpreter

	contracts           map[string]map[string]*Method
	requests            map[string]map[string]interface{}
	postProcessContract map[string]interface{}

	previousBlock  *Block
	roundTimestamp int64
	transactions   *map[string]*SignedTransaction
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

func (m *Machine) GetInterpreter() *Interpreter {
	return m.interpreter
}

func (m *Machine) Init(previousBlock *Block, roundTimestamp int64) {
	m.interpreter.Reset()
	m.interpreter.Init("transaction")

	m.previousBlock = previousBlock
	m.roundTimestamp = roundTimestamp
}

func (m *Machine) ValidateTxTimestamp(tx *SignedTransaction) bool {
	if tx.GetTimestamp() > int64(util.Utime())+C.TIME_STAMP_ERROR_LIMIT {
		return false
	}

	return true
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

func (m *Machine) SetTransactions(txs map[string]*SignedTransaction) {
	m.transactions = &txs
}

func (m *Machine) TxValidity(tx *SignedTransaction) (bool, error) {
	size, err := tx.GetSize()
	if err != nil {
		return false, err
	}

	if size > C.TX_SIZE_LIMIT {
		return false, fmt.Errorf("The length of the signed transaction must be less than %d characters", C.TX_SIZE_LIMIT)
	}

	if err := tx.Validate(); err != nil {
		return false, err
	}

	chain := storage.GetChainStorageInstance()
	lastBlock, err := chain.LastBlock()
	if err != nil {
		return false, err
	}
	roundTimestamp := util.Utime() + C.TIME_STAMP_ERROR_LIMIT

	m.Init(lastBlock, roundTimestamp)

	txHash, err := tx.GetTxHash()
	if err != nil {
		return false, err
	}

	m.SetTransactions(map[string]*SignedTransaction{
		txHash: tx,
	})

	if err := m.PreCommit(); err != nil {
		return false, err
	}

	return true, nil
}

func (m *Machine) loadContracts() {
	m.contracts = rpc.NativeContracts()
}

func (m *Machine) PreCommit() error {
	m.loadContracts()

	txs := util.SortedValueK[*SignedTransaction](*m.transactions)

	for _, tx := range txs {
		txHash, err := tx.GetTxHash()
		DebugLog("preCommitRead	:", txHash)
		if err != nil {
			delete(*m.transactions, txHash)
			continue
		}

		if m.ValidateTxTimestamp(tx) && tx.Validate() == nil {
			if err := m.MountContract(tx); err == nil {
				if err := m.interpreter.ParameterValidate(); err == nil {
					m.interpreter.Read()
				} else {
					DebugPanic("ParameterValidate error: %v", err)
				}
			} else {
				DebugPanic("MountContract error: %v", err)
			}
		}

		delete(*m.transactions, txHash)
	}

	for _, transaction := range txs {
		hash, _ := transaction.GetTxHash()
		if err := m.MountContract(transaction); err == nil {
			if result, err := m.interpreter.Execute(); err == nil {
				DebugLog("Execute ", hash, " result:", result)
				continue
			}
		} else {
			DebugPanic("MountContract error: %v", err)
		}

		delete(*m.transactions, hash)
	}

	return nil
}

func (m *Machine) MountContract(tx *SignedTransaction) error {
	txMap := tx.GetTxData()
	DebugLog("MountContract tx:", string(txMap.Data.Ser()))
	DebugLog("MountContract tx type:", txMap.Type)
	DebugLog("MountContract tx cid:", tx.GetCID())

	cid := tx.GetCID()

	txType, ok := txMap.Data.Get("type")
	if !ok {
		return fmt.Errorf("transaction type not found")
	}
	name := txType.(string)

	code, ok := m.contracts[cid][name]
	if !ok {
		return fmt.Errorf("contract not found for cid %s and method %s", cid, name)
	}

	if code == nil {
		return fmt.Errorf("contract code is nil")
	}

	if cid == F.RootSpaceId() && slices.Contains(rpc.SystemMethods, name) {
		m.interpreter.Set(tx.GetTxData(), code, new(Method))
	}

	return nil
}

func (m *Machine) NextBlock() *Block {
	return &Block{
		Height:            m.previousBlock.Height + 1,
		Transactions:      m.transactions,
		Timestamp_s:       m.roundTimestamp,
		UniversalUpdates:  m.interpreter.GetUniversalUpdates(),
		LocalUpdates:      m.interpreter.GetLocalUpdates(),
		PreviousBlockhash: m.previousBlock.BlockHash(),
	}
}

func (m *Machine) TimeValidity(tx *SignedTransaction, timestamp int64) (bool, error) {
	if m.previousBlock.Timestamp_s < tx.GetTimestamp() && tx.GetTimestamp() <= int64(timestamp) {
		return true, nil
	}

	return false, fmt.Errorf("Timestamp must be greater than %d and less than %d", m.previousBlock.Timestamp_s, timestamp)
}
