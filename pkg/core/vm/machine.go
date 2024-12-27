package vm

import (
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
	"hello/pkg/rpc"
	"hello/pkg/util"
	"sync"
	"go.uber.org/zap"
)

type Machine struct {
	mu sync.RWMutex

	interpreter *Interpreter

	contracts           map[string]map[string]*Method
	requests            map[string]map[string]*Method
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
	m.interpreter.Reset(true)
	m.interpreter.Init("transaction")

	m.previousBlock = previousBlock
	m.roundTimestamp = roundTimestamp
}

func (m *Machine) ValidateTxTimestamp(tx *SignedTransaction) bool {
	if m.previousBlock != nil && tx.GetTimestamp() <= m.previousBlock.GetTimestamp() {
		return false
	}

	if tx.GetTimestamp() > int64(util.Utime())+C.TIME_STAMP_ERROR_LIMIT {
		return false
	}

	return true
}

func (m *Machine) Commit(block *Block) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	chain := storage.GetChainStorageInstance()
	sf := storage.GetStatusFileInstance()

	DebugLog("Commit block:", block.BlockHash())

	if err := chain.Write(block); err != nil {
		return err
	}

	/**
	if err := sf.Update(block); err != nil {
		return err
	}
	**/

	if err := sf.Write(block); err != nil {
		return err
	}

	m.previousBlock = block
	m.roundTimestamp = block.Timestamp_s
	m.transactions = &map[string]*SignedTransaction{}

	return nil
}

func (m *Machine) PreLoad(universalUpdates *map[string]Update) {
}

func (m *Machine) SetTransactions(txs map[string]*SignedTransaction) {
	m.transactions = &txs
}

func (m *Machine) TxValidity(tx *SignedTransaction) (bool, error) {
	logger.Info(
		fmt.Sprintf("txvalidity ", tx.GetType()), 
		zap.String("cid", tx.GetCID()),
	)

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

	m.mu.Lock()
	defer m.mu.Unlock()

	// txHash := tx.GetTxHash()
	//m.SetTransactions(map[string]*SignedTransaction{ txHash: tx, })

	if err := m.PreCommitOne(tx); err != nil {
		return false, err
	}

	if m.interpreter.GetResult() != "" {
		return false, fmt.Errorf("transaction is not valid: %s", m.interpreter.GetResult())
	}

	return true, nil
}

func (m *Machine) loadContracts() {
	m.contracts = rpc.NativeContracts()
}

func (m *Machine) deleteTransaction(txHash string, err error) {
	logger.Info(
		"delete transaction from precommit", 
		zap.String("hash", txHash), 
		zap.String("err", err.Error()),
	)

	delete(*m.transactions, txHash)
}

func (m *Machine) PreCommit() error {
	m.loadContracts()
	var err error
	txs := util.SortedValueK[*SignedTransaction](*m.transactions)

	logger.Info("precommit")

	for _, tx := range txs {
		txHash := tx.GetTxHash()

		if !(m.ValidateTxTimestamp(tx)) {
			err = fmt.Errorf("tx timestamp error: %s", txHash)
			m.deleteTransaction(txHash, err)
			continue
		}

		if err = tx.Validate(); err != nil {
			m.deleteTransaction(txHash, err)
			continue
		}

		if err = m.MountContract(*tx); err != nil {
			m.deleteTransaction(txHash, err)
			continue
		}

		if err := m.interpreter.ParameterValidate(); err != nil {
			m.deleteTransaction(txHash, err)
			continue
		}

		m.interpreter.Read()
	}

	m.interpreter.LoadUniversalStatus()
	//m.interpreter.LoadLocalStatus()
	txs = util.SortedValueK[*SignedTransaction](*m.transactions)

	for _, transaction := range txs {
		hash := transaction.GetTxHash()
		if err := m.MountContract(*transaction); err == nil {
			if result, err := m.interpreter.Execute(); err == nil {
				DebugLog("Execute ", hash, " result:", result)
				continue
			}
		}

		m.deleteTransaction(hash, err)
	}

	return err
}

func (m *Machine) MountContract(tx SignedTransaction) error {
	txMap := tx.GetTxData()
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

	m.interpreter.Set(tx.GetTxData(), code, new(Method))

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

func (m *Machine) Epoch() string {
	currentTime := util.Utime() / 1000
	timeInEpoch := currentTime % 3000

	if timeInEpoch < 2000 {
		return "txtime"
	}
	return "blocktime"
}

func (m *Machine) IsInBlockTime() bool {
	currentTime := util.Utime()
	timeInEpoch := currentTime % 5000

	return timeInEpoch >= 3000 // 마지막 2초는 블록 생성 시간
}

func (m *Machine) GetCurrentEpoch() int64 {
	return util.Utime() / 5000
}

func (m *Machine) GetPreviousBlock() *Block {
	return m.previousBlock
}

func (m *Machine) ExpectedBlock() *Block {
	previousBlock := m.GetPreviousBlock()

	var previousBlockhash string
	var Height int

	if previousBlock == nil {
		previousBlockhash = ""
		Height = 0
	} else {
		previousBlockhash = previousBlock.BlockHash()
		Height = previousBlock.Height
	}

	expectedBlock := &Block{
		Height:            Height + 1,
		Transactions:      m.transactions,
		Timestamp_s:       m.roundTimestamp,
		UniversalUpdates:  m.interpreter.GetUniversalUpdates(),
		LocalUpdates:      m.interpreter.GetLocalUpdates(),
		PreviousBlockhash: previousBlockhash,
	}

	return expectedBlock
}
func (m *Machine) Response(request SignedRequest) (interface{}, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	m.interpreter.Reset(true)
	m.interpreter.Init("request")

	m.loadRequests()

	code := m.suitedRequest(request)
	if code == nil {
		return nil, fmt.Errorf("request code not found: %s", request.GetRequestType())
	}

	m.interpreter.Set(request.GetRequestData(), code, new(Method))

	if err := m.interpreter.ParameterValidate(); err != nil {
		return nil, err
	}

	m.interpreter.Read()
	m.interpreter.LoadUniversalStatus()

	_, result := m.interpreter.Execute()

	fmt.Println("result:", result)

	return result, nil
}

func (m *Machine) loadRequests() {
	m.requests = rpc.NativeRequests()
}

func (m *Machine) suitedRequest(request SignedRequest) *Method {
	requestType := request.GetRequestType()

	if methods, ok := m.requests[request.GetRequestCID()]; ok {
		if method, exists := methods[requestType]; exists {
			return method
		}
	}
	return nil
}

func (m *Machine) PreCommitOne(tx *SignedTransaction) error {
	m.loadContracts()
	var err error

	txHash := tx.GetTxHash()

	if !(m.ValidateTxTimestamp(tx)) {
		return fmt.Errorf("tx timestamp error: %s", txHash)
	}

	if err = tx.Validate(); err != nil {
		return err
	}

	if err = m.MountContract(*tx); err != nil {
		return err
	}

	if err := m.interpreter.ParameterValidate(); err != nil {
		return err
	}

	m.interpreter.Read()
	m.interpreter.LoadUniversalStatus()

	result, err := m.interpreter.Execute()

	if result != nil && result != "" {
		return fmt.Errorf("%v", result)
	}

	return err
}
