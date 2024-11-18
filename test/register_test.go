package main

import (
	"fmt"
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	. "hello/pkg/core/vm/native"
)

func main() {

	userMethod := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "10",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
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

	registerMethod := Register()
	registerMethod.AddParameter(NewParameter(map[string]interface{}{
		"name":         "code",
		"type":         "string",
		"maxlength":    65536,
		"requirements": true,
	}))

	registerMethod.AddExecution(abi.WriteLocal("contract", abi.IDHash([]interface{}{"userMethod", 1}), userMethod.GetCodeRaw()))
	fmt.Println("User contract registered successfully.")
}
