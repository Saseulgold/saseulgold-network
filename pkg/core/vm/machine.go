package vm

import (
	"encoding/json"
	"fmt"
	. "hello/pkg/core/model"
)

type Machine struct {
	interpreter *Interpreter

	contracts           map[string]map[string]interface{}
	requests            map[string]map[string]interface{}
	postProcessContract map[string]interface{}

	previousBlock  *Block
	roundTimestamp int64
	transactions   *map[string]SignedTransaction

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

func (m *Machine) CurrentBlockDifficulty(diff string) string {
	if diff != "" {
		m.currentBlockDifficulty = diff
	}
	return m.currentBlockDifficulty
}

func (m *Machine) SetVoutInfo(vout string, nonce string) {
	m.currentBlockVout = vout
	m.currentBlockNonce = nonce
}

func (m *Machine) SetTransactions(transactions map[string]interface{}) {
	m.transactions = make(map[string]*SignedTransaction)

	for key, tx := range transactions {
		if txMap, ok := tx.(map[string]interface{}); ok {
			m.transactions[key] = &NewSignedTransaction(txMap)
		}
	}
}

func (m *Machine) TransactionCount() int {
	return len(m.transactions)
}

func (m *Machine) LoadContracts() {
	codes := rpc.GetContracts()

	if methods, ok := codes["methods"].(map[string]map[string]interface{}); ok {
		m.contracts = methods
	}
	if postProcess, ok := codes["post_process"].(map[string]interface{}); ok {
		m.postProcessContract = postProcess
	}
}

func (m *Machine) LoadRequests() {
	m.requests = rpc.GetRequests()
}

func (m *Machine) MountContract(transaction *model.SignedTransaction) (bool, string) {
	cid := transaction.Cid
	if cid == "" {
		cid = config.RootSpaceId()
	}
	name := transaction.Type

	logger.Log("Mount Contract: cid=" + cid + "; name=" + name)
	logger.Log("rootspace=" + config.RootSpace() + " ;rootspaceid: " + config.RootSpaceId() + "; encode json: " + string(json.Marshal(m.contracts)))

	code, ok := m.contracts[cid][name]
	if !ok {
		return false, "There is no contract code: " + cid + " " + transaction.Type
	}

	if cid == config.RootSpaceId() && rpc.IsSystemMethod(name) {
		m.interpreter.Set(transaction, code, model.NewMethod())
	} else {
		m.interpreter.Set(transaction, code, m.postProcessContract)
	}

	return true, ""
}

func (m *Machine) SuitedRequest(request *model.SignedRequest) *model.Method {
	cid := request.Cid
	if cid == "" {
		cid = hasher.SpaceId(config.ZERO_ADDRESS, config.RootSpace())
	}
	name := request.Type
	if name == "" {
		name = ""
	}

	if methods, ok := m.requests[cid]; ok {
		if method, ok := methods[name]; ok {
			return method
		}
	}
	return nil
}

func (m *Machine) Chunk() *model.Chunk {
	chunk := &model.Chunk{
		PreviousBlockhash: m.previousBlock.Blockhash,
		STimestamp:        m.roundTimestamp,
		Transactions:      getMapKeys(m.transactions),
	}

	chunk.SignChunk(env.Node())

	return chunk
}

func (m *Machine) Hypothesis(chunks []interface{}) *model.Hypothesis {
	hypothesis := &model.Hypothesis{
		PreviousBlockhash: m.previousBlock.Blockhash,
		Chunks:            chunks,
		STimestamp:        m.roundTimestamp,
		Thashs:            getMapKeys(m.transactions),
	}

	hypothesis.SignHypothesis(env.Node())

	return hypothesis
}

func (m *Machine) ExpectedBlock(seal []interface{}) *model.MainBlock {
	expectedBlock := &Block{
		Height:            m.previousBlock.Height + 1,
		Transactions:      m.transactions,
		Timestamp_s:       m.roundTimestamp,
		UniversalUpdates:  m.interpreter.GetUniversalUpdates(),
		LocalUpdates:      m.interpreter.GetLocalUpdates(),
		PreviousBlockhash: m.previousBlock.Blockhash,
		Vout:              m.currentBlockVout,
		Nonce:             m.currentBlockNonce,
	}

	expectedBlock.MakeBlockhash()
	logger.Log("expectedBlock: " + json.Marshal(expectedBlock))
	return expectedBlock
}

func (m *Machine) TimeValidity(transaction *SignedTransaction, timestamp int64) (bool, string) {
	// min < time <= max
	if m.previousBlock.GetTimestamp() < transaction.GetTimestamp() && transaction.GetTimestamp() <= timestamp {
		return true, ""
	}

	errMsg := fmt.Sprintf("Timestamp must be greater than %d and less than %d", m.previousBlock.GetTimestamp(), timestamp)
	return false, errMsg
}

func (m *Machine) PreLoad(universalUpdates map[string]interface{}, localUpdates map[string]interface{}) {
	for key, update := range universalUpdates {
		if updateMap, ok := update.(map[string]interface{}); ok {
			if old, exists := updateMap["old"]; exists {
				m.interpreter.SetUniversalLoads(key, old)
			}
		}
	}

	for key, update := range localUpdates {
		if updateMap, ok := update.(map[string]interface{}); ok {
			if old, exists := updateMap["old"]; exists {
				m.interpreter.SetLocalLoads(key, old)
			}
		}
	}
}

func (m *Machine) Commit(block *Block) bool {
	// if block.GetTimestamp() < clock.Utime()+config.TimestampErrorLimit {
	if true {
		if m.Write(block) {
			return true
		}
	}

	return false
}
