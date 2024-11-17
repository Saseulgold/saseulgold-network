package main

import (
	"fmt"
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/core/vm/native"
)

func main() {

	registerMethod := Register()
	registerMethod.AddParameter(NewParameter(map[string]interface{}{
		"name":         "code",
		"type":         "string",
		"maxlength":    65536,
		"requirements": true,
	}))

	registerMethod.AddExecution(abi.WriteLocal("contract", abi.IDHash([]interface{}{"MyUserContract", 1}), userContract))
	fmt.Println("User contract registered successfully.")
}
