package native

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/util"
)

func Genesis() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Genesis",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "network_address",
		"type":         "string",
		"maxlength":    ID_HASH_SIZE,
		"requirements": true,
	}))

	genesis := abi.ReadLocal("genesis", "00", nil)

	method.AddExecution(abi.Condition(
		abi.Ne(genesis, true),
		"There was already a Genesis.",
	))

	return method
}

func Faucet() *Method {
	// For testing
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Faucet",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	from := abi.Param("from")
	balance := "100000000000000000"
	method.AddExecution(abi.WriteUniversal("balance", from, balance))

	return method
}

func Register() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Register",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "code",
		"type":         "string",
		"maxlength":    65536,
		"requirements": true,
	}))

	from := abi.Param("from")
	code := abi.Param("code")

	decodedCode := abi.DecodeJSON(code)

	codeType := abi.Get(decodedCode, "type")
	name := abi.Param("name")
	nonce := abi.Get(decodedCode, "nonce")
	version := abi.Get(decodedCode, "version")
	writer := abi.Get(decodedCode, "writer")

	codeID := abi.IDHash([]interface{}{name, nonce})

	contractInfo := abi.DecodeJSON(abi.ReadLocal("contract", codeID, nil))
	requestInfo := abi.DecodeJSON(abi.ReadLocal("request", codeID, nil))

	contractVersion := abi.Get(contractInfo, "version")
	requestVersion := abi.Get(requestInfo, "version")

	isNetworkManager := abi.ReadLocal("network_manager", from, nil)
	method.AddExecution(abi.Condition(
		abi.Eq(isNetworkManager, true),
		"You are not network manager.",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(writer, ZERO_ADDRESS),
		"Writer must be zero address",
	))

	method.AddExecution(abi.Condition(
		abi.IsString([]interface{}{codeType}),
		"Invalid type",
	))

	method.AddExecution(abi.Condition(
		abi.In(codeType, []interface{}{"contract", "request"}),
		"Type must be one of the following: contract, request",
	))

	method.AddExecution(abi.Condition(
		abi.IsString([]interface{}{name}),
		"Invalid name",
	))

	method.AddExecution(abi.Condition(
		abi.RegMatch("^[A-Za-z_0-9]+$", name),
		"The name must consist of A-Za-z_0-9.",
	))

	method.AddExecution(abi.Condition(
		abi.IsNumeric([]interface{}{version}),
		"Invalid version",
	))

	versionCheck := abi.If(
		abi.Eq(codeType, "contract"),
		abi.Gt(version, contractVersion),
		abi.If(
			abi.Eq(codeType, "request"),
			abi.Gt(version, requestVersion),
			false,
		),
	)

	method.AddExecution(abi.Condition(
		versionCheck,
		"Only new versions of code can be registered.",
	))

	update := abi.If(
		abi.Eq(codeType, "contract"),
		abi.WriteLocal("contract", codeID, code),
		abi.If(
			abi.Eq(codeType, "request"),
			abi.WriteLocal("request", codeID, code),
			nil,
		),
	)
	method.AddExecution(update)

	return method
}

func Revoke() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Revoke",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	from := abi.Param("from")
	isNetworkManager := abi.ReadLocal("network_manager", from, nil)

	method.AddExecution(abi.Condition(
		abi.Eq(isNetworkManager, true),
		"You are not network manager.",
	))

	method.AddExecution(abi.WriteLocal("network_manager", from, false))

	return method
}

func Send() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Send",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	from := abi.Param("from")
	to := abi.Param("to")
	amount := abi.Param("amount")

	fromBalance := abi.ReadUniversal("balance", from, "0")
	toBalance := abi.ReadUniversal("balance", to, "0")

	method.AddExecution(abi.Condition(
		abi.Ne(from, to),
		"Sender and receiver address must be different.",
	))

	method.AddExecution(abi.Condition(
		abi.Gt(fromBalance, amount),
		"Balance is not enough.",
	))

	method.AddExecution(abi.WriteUniversal("balance", from, abi.PreciseSub(fromBalance, amount, 0)))
	method.AddExecution(abi.WriteUniversal("balance", to, abi.PreciseAdd(toBalance, amount, 0)))

	return method
}

func Publish() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Publish",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "code",
		"type":         "string",
		"maxlength":    65536,
		"requirements": true,
	}))

	from := abi.Param("from")
	code := abi.Param("code")

	decodedCode := abi.DecodeJson(code)

	codeType := abi.Get(decodedCode, "t")
	name := abi.Get(decodedCode, "n")
	space := abi.Get(decodedCode, "s")
	version := abi.Get(decodedCode, "v")
	writer := abi.Get(decodedCode, "w")

	codeID := abi.Hash([]interface{}{writer, space, name})

	contractInfo := abi.ReadLocal("contract", codeID, nil)
	contractInfo = abi.DecodeJson(contractInfo)

	requestInfo := abi.ReadLocal("request", codeID, nil)
	requestInfo = abi.DecodeJson(requestInfo)

	contractVersion := abi.Get(contractInfo, "v")
	requestVersion := abi.Get(requestInfo, "v")

	method.AddExecution(abi.Condition(
		abi.Eq(writer, from),
		"Writer must be the same as the from address",
	))

	method.AddExecution(abi.Condition(
		abi.IsString([]interface{}{codeType}),
		abi.Concat([]interface{}{"Invalid type: ", codeType}),
	))

	method.AddExecution(abi.Condition(
		abi.In(codeType, []interface{}{"contract", "request"}),
		"Type must be one of the following: contract, request",
	))

	method.AddExecution(abi.Condition(
		abi.IsString([]interface{}{name}),
		abi.Concat([]interface{}{"Invalid name: ", name}),
	))

	method.AddExecution(abi.Condition(
		abi.RegMatch("/^[A-Za-z_0-9]+$/", name),
		"The name must consist of A-Za-z_0-9",
	))

	method.AddExecution(abi.Condition(
		abi.IsNumeric([]interface{}{version}),
		abi.Concat([]interface{}{"Invalid version: ", version}),
	))

	method.AddExecution(abi.Condition(
		abi.IsString([]interface{}{space}),
		"invalid nonce",
		// abi.Concat([]interface{}{"Invalid nonce: ", space}),
	))

	method.AddExecution(abi.Condition(
		abi.If(
			abi.Eq(codeType, "contract"),
			abi.Gt(version, contractVersion),
			abi.If(
				abi.Eq(codeType, "request"),
				abi.Gt(version, requestVersion),
				false,
			),
		),
		"Only new versions of code can be registered",
	))

	method.AddExecution(abi.If(
		abi.Eq(codeType, "contract"),
		abi.WriteLocal("contract", codeID, code),
		abi.If(
			abi.Eq(codeType, "request"),
			abi.WriteLocal("request", codeID, code),
			false,
		),
	))

	return method
}
