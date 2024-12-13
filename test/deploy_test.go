package main

import (
	"fmt"
	abi "hello/pkg/core/abi"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	S "hello/pkg/core/structure"
	. "hello/pkg/core/vm"
	service "hello/pkg/service"
	"testing"
)

func TestDeploy(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "register_test"
	from := "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4"

	userMethod := &Method{
		Type:    "contract",
		Name:    "userMethod",
		Space:   "test",
		Writer:  from,
		Version: "1",
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"type":         "string",
				"maxlength":    10,
				"requirements": true,
				"default":      "10",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"type":         "string",
				"maxlength":    10,
				"requirements": true,
				"default":      "20",
			}),
		},
		Executions: []Execution{
			abi.And(
				abi.Gt(abi.Param("value1"), "5"),
				abi.Lt(abi.Param("value2"), "30"),
			),
			abi.Or(
				abi.Eq(abi.Param("value1"), "5"),
				abi.Gt(abi.Param("value2"), "15"),
			),
			abi.And(
				abi.Or(
					abi.Gt(abi.Param("value1"), "5"),
					abi.Lt(abi.Param("value2"), "10"),
				),
				abi.Eq(abi.Param("value1"), abi.Param("value2")),
			),
		},
	}

	code := userMethod.GetCode()
	if code == "" {
		t.Error("Failed to get user method code")
	}

	fmt.Println("User Method Code:", code)

	interpreter := NewInterpreter()
	interpreter.Init("transaction")

	data := S.NewOrderedMap()
	txData := S.NewOrderedMap()
	txData.Set("type", "Publish")
	txData.Set("code", code)
	txData.Set("from", from)
	txData.Set("timestamp", int64(1733211394654000))

	data.Set("transaction", txData)
	data.Set("public_key", "8860ecfe5711c9096f43411ce1ebefcb292200fbca73aa14fbf187a52cc29898")
	data.Set("signature", "1493bd19ea174751810b3fece0f23fa24c9e6d884118624e863b8fc3892f5604dba420407e2a833440d58a2e3c3e1048109f4f14f5f9fd3c9788a86e0bc5f400")

	tx, err := NewSignedTransaction(data)
	if err != nil {
		t.Errorf("NewSignedTransaction(): %s", err)
	}

	service.ForceCommit(map[string]*SignedTransaction{tx.GetTxData().Cid: &tx})

}
