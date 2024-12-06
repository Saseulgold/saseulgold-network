package vm

import (
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	"hello/pkg/core/model"
	S "hello/pkg/core/structure"
	. "hello/pkg/core/vm"
	"hello/pkg/util"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreCommit(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "precommit_test"

	machine := GetMachineInstance()
	machine.GetInterpreter().Reset()

	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", "Faucet")
	txData.Set("from", "60c3a6cd858c90574bcdc35b2da5dbc7225275f50efd")
	txData.Set("timestamp", util.Utime()-1000)

	data.Set("transaction", txData)
	data.Set("public_key", "test_public_key")
	data.Set("signature", "test_signature")

	tx0, err := model.NewSignedTransaction(data)

	// 테스트용 트랜잭션 생성
	data = S.NewOrderedMap()
	txData = S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50edf")
	txData.Set("from", "60c3a6cd858c90574bcdc35b2da5dbc7225275f50efd")
	txData.Set("amount", 1000)
	txData.Set("timestamp", util.Utime())

	data.Set("transaction", txData)
	data.Set("public_key", "test_public_key")
	data.Set("signature", "test_signature")

	tx1, err := model.NewSignedTransaction(data)
	assert.NoError(t, err)

	if err != nil {
		t.Errorf("NewSignedTransaction error: %v", err)
	}

	txHash0, err := tx0.GetTxHash()
	_, err = tx1.GetTxHash()
	assert.NoError(t, err)

	txHash1, err := tx1.GetTxHash()
	assert.NoError(t, err)

	if err != nil {
		t.Errorf("GetTxHash error: %v", err)
	}

	// Machine에 트랜잭션 설정
	txs := map[string]*model.SignedTransaction{
		txHash0: &tx0,
		txHash1: &tx1,
	}

	machine.Init(nil, 0)
	machine.SetTransactions(txs)

	// PreCommit 실행
	err = machine.PreCommit()
	log.Println(machine.GetInterpreter().GetUniversalUpdates())

	if err != nil {
		t.Errorf("PreCommit error: %v", err)
	}

	assert.Equal(t, 2, len(*machine.GetInterpreter().GetUniversalUpdates()))
}

func TestPreCommit2(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "precommit_test"

	machine := GetMachineInstance()
	machine.GetInterpreter().Reset()

	// 테스트용 트랜잭션 생성
	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "50c3a6cd858c90574bcdc35b2da5dbc7225275f50efd")
	txData.Set("from", "60c3a6cd858c90574bcdc35b2da5dbc7225275f50edd")
	txData.Set("amount", "10000000")
	txData.Set("timestamp", util.Utime())

	data.Set("transaction", txData)
	data.Set("public_key", "test_public_key")
	data.Set("signature", "test_signature")

	tx1, err := model.NewSignedTransaction(data)
	assert.NoError(t, err)

	if err != nil {
		t.Errorf("NewSignedTransaction error: %v", err)
	}

	txHash1, err := tx1.GetTxHash()
	assert.NoError(t, err)

	if err != nil {
		t.Errorf("GetTxHash error: %v", err)
	}

	// Machine에 트랜잭션 설정
	txs := map[string]*model.SignedTransaction{
		txHash1: &tx1,
	}

	DebugLog("BeforeInit: ", machine.GetInterpreter().GetUniversalUpdates())
	machine.Init(nil, 0)
	DebugLog("AfterInit: ", machine.GetInterpreter().GetUniversalUpdates())
	machine.SetTransactions(txs)

	// PreCommit 실행
	err = machine.PreCommit()
	log.Println(machine.GetInterpreter().GetUniversalUpdates())
	assert.Equal(t, 0, len(*machine.GetInterpreter().GetUniversalUpdates()))
}
