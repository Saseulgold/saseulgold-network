package main

import (
	"fmt"
	abi "hello/pkg/core/abi"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/core/vm"
	. "hello/pkg/core/vm/native"
	"testing"
)

func TestDeploy(t *testing.T) {
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "register_test"

	userMethod := &Method{
		Type:    "contract",
		Name:    "userMethod",
		Space:   "test",
		Writer:  "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4",
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

	publishMethod := Publish()

	interpreter := NewInterpreter()
	interpreter.Init("transaction")

	signedData := NewSignedData()
	signedData.SetAttribute("code", code)
	signedData.SetAttribute("from", userMethod.Writer)

	post := &Method{
		Type:       "contract",
		Parameters: Parameters{},
		Executions: []Execution{},
	}

	interpreter.Reset()
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(publishMethod)
	interpreter.SetPostProcess(post)

	_, msg := interpreter.Execute()
	result := interpreter.GetResult()
	t.Logf("Error message: %v", msg)
	t.Logf("Result: %v", result)

	updates := interpreter.GetLocalUpdates()
	t.Logf("Updates: %v", updates)
	abi.DebugLog("Writer:", userMethod.Writer)

	/**
	prefix := "contract"
	statusHash := F.StatusHash(publishMethod.GetWriter(), publishMethod.GetSpace(), prefix, publishMethod.GetName())
	t.Logf("Status hash: %v", interpreter.GetLocalStatus(statusHash, nil))
	*/

}
