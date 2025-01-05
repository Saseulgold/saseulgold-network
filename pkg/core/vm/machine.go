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

	previousBlock *Block
	transactions  *map[string]*SignedTransaction
}

var instance *Machine

func GetMachineInstance() *Machine {
	if instance == nil {
		instance = &Machine{
			interpreter:   NewInterpreter(),
			previousBlock: nil,
		}
	}
	return instance
}

func NewMachine(previousBlock *Block) *Machine {
	m := &Machine{
		interpreter:   NewInterpreter(),
		previousBlock: previousBlock,
	}

	m.Init(previousBlock)
	return m
}

func (m *Machine) GetInterpreter() *Interpreter {
	return m.interpreter
}

func (m *Machine) Init(previousBlock *Block) {
	m.interpreter.Reset(true)
	m.interpreter.Init("transaction")

	m.previousBlock = previousBlock

	m.loadContracts()
	m.loadRequests()
}

func (m *Machine) ValidateTxTimestamp(tx *SignedTransaction) bool {
	if m.previousBlock != nil && tx.GetTimestamp() <= m.previousBlock.GetTimestamp()-C.TIME_STAMP_FORWARD_ERROR_LIMIT {
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
	block.Timestamp_s = util.Utime()
	m.transactions = &map[string]*SignedTransaction{}

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
	logger.Info("txValidity", zap.String("tx", fmt.Sprintf("%v", tx)))

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
	m.Init(lastBlock)

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

func (m *Machine) loadRequests() {
	m.requests = rpc.NativeRequests()
}

func (m *Machine) loadUserDefinedContract(cid string, name string) *Method {
	contractKey := cid + util.FillHashSuffix(name)
	si := storage.GetStatusFileInstance()

	cursor, ok := si.CachedUniversalIndexes[contractKey]
	fmt.Println("contractKey:", contractKey)
	fmt.Println("cursor:", cursor)
	if !ok {
		return nil
	}

	data, err := si.ReadUniversalStatus(cursor)
	if err != nil {
		return nil
	}

	contract := ParseMethod(data.(string))

	return contract
}

func (m *Machine) deleteTransaction(txHash string, err error) {
	delete(*m.transactions, txHash)

	if err != nil {
		logger.Error("deleteTransaction", zap.String("txHash", txHash), zap.Error(err))
	}
}

func (m *Machine) PreCommit() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	logger.Info("preCommit", zap.String("transactions", fmt.Sprintf("%v", *m.transactions)))

	m.loadContracts()
	var err error
	txs := util.SortedValueK[*SignedTransaction](*m.transactions)

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
			logger.Info("excute", zap.String("hash", hash))
			if result, err := m.interpreter.Execute(); err == nil {
				logger.Info("excuted", zap.String("hash", hash), zap.Any("result", result))
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
	var code *Method

	if cid == util.RootSpaceId() {
		code = m.contracts[cid][name]
		if code == nil {
			return fmt.Errorf("contract not found for cid %s and method %s", cid, name)
		}

		m.interpreter.Set(tx.GetTxData(), code.Copy(), new(Method))
		return nil
	} else {
		code = m.loadUserDefinedContract(cid, name)
	}

	if code == nil {
		return fmt.Errorf("contract code is nil")
	}

	m.interpreter.Set(tx.GetTxData(), code.Copy(), new(Method))

	return nil
}

func (m *Machine) NextBlock() *Block {
	return &Block{
		Height:            m.previousBlock.Height + 1,
		Transactions:      m.transactions,
		Timestamp_s:       util.Utime(),
		UniversalUpdates:  m.interpreter.GetUniversalUpdates(),
		LocalUpdates:      m.interpreter.GetLocalUpdates(),
		PreviousBlockhash: m.previousBlock.BlockHash(),
	}
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
		Timestamp_s:       util.Utime(),
		UniversalUpdates:  m.interpreter.GetUniversalUpdates(),
		LocalUpdates:      m.interpreter.GetLocalUpdates(),
		PreviousBlockhash: previousBlockhash,
	}

	return expectedBlock
}
func (m *Machine) Response(request SignedRequest) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	logger.Info("Response", zap.String("request", fmt.Sprintf("%v", request)))

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
	m.mu.Lock()
	defer m.mu.Unlock()

	logger.Info("preCommitOne", zap.String("tx", fmt.Sprintf("%v", tx)))

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
