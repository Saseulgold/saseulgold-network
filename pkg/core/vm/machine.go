package vm

import (
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
	"hello/pkg/rpc"
	"hello/pkg/util"
	F "hello/pkg/util"
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

func (m *Machine) Init(previousBlock *Block, roundTimestamp int64) {
	m.interpreter.Reset()
	m.interpreter.Init("transaction")

	m.previousBlock = previousBlock
	m.roundTimestamp = roundTimestamp
}

func (m *Machine) ValidateTxTimestamp(tx *SignedTransaction) bool {
	if tx.GetTimestamp() > int(util.Utime())+C.TIME_STAMP_ERROR_LIMIT {
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

func (m *Machine) PreCommit() error {
	code := &rpc.Code{}
	m.contracts = code.Contracts()

	txs := util.SortedValueK[*SignedTransaction](*m.transactions)

	for _, tx := range txs {
		txHash, err := tx.GetTxHash()
		if err != nil {
			delete(*m.transactions, txHash)
			continue
		}

		if m.ValidateTxTimestamp(tx) && tx.Validate() == nil {

			if err := m.MountContract(tx); err != nil {
				continue
			}
		}

		// Invalid transaction
		delete(*m.transactions, txHash)
	}

	return nil
}

func (m *Machine) MountContract(tx *SignedTransaction) error {
	txMap := tx.GetTx()

	cid, ok := txMap.Get("cid")
	if !ok {
		cid = F.RootSpaceId()
	}

	txType, ok := txMap.Get("type")
	if !ok {
		return fmt.Errorf("transaction type not found")
	}
	name := txType.(string)

	code, ok := m.contracts[cid.(string)][name]
	if !ok {
		return fmt.Errorf("contract not found for cid %s and method %s", cid, name)
	}

	if code == nil {
		return fmt.Errorf("contract code is nil")
	}
	/**
		if ($cid === Config::rootSpaceId() && in_array($name, Code::SYSTEM_METHODS)) {
			$this->interpreter->set($transaction, $code, new Method());
	} else {
			$this->interpreter->set($transaction, $code, $this->post_process_contract);
		}
		**/

	return nil
}
