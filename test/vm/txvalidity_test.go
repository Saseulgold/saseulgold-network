package vm

import (
	C "hello/pkg/core/config"
	"hello/pkg/core/model"
	S "hello/pkg/core/structure"
	. "hello/pkg/core/vm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxValidity(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "genesis_test_2"

	machine := GetMachineInstance()
	machine.GetInterpreter().Reset()

	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", "Send")
	txData.Set("to", "464708fbb8d8c3af1b3470b2eecd538f23be08b18b1f")
	txData.Set("amount", "100")
	txData.Set("from", "4a8fd2ebb308370a689c3ef47cb83ba182683def3d4f")
	txData.Set("timestamp", 1733220569173000)

	data.Set("transaction", txData)
	data.Set("public_key", "bdfcefde7c536e8342f1ec65c69373f9ff46f33c18acf0f5848c71e037eca9f2")
	data.Set("signature", "f2d8bb8a72f8d77caa44e0e153e431b60e3dfd997e0d12d47ae6b026d6b5c6d9c3d3daf7565e40b5329f955706baf15a32fb144eda52aaeebb1715460586230b")

	tx0, err := model.NewSignedTransaction(data)
	// 테스트용 트랜잭션 생성

	valid, err := machine.TxValidity(&tx0)
	assert.NoError(t, err)
	assert.False(t, valid)

}
