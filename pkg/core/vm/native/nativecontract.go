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

	genesis := abi.ReadUniversal("genesis", "00", nil)

	method.AddExecution(abi.Condition(
		abi.Ne(genesis, true),
		"There was already a Genesis.",
	))

	method.AddExecution(abi.WriteUniversal("genesis", "00", true))

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
	balance := "100000000000000000000000000"

	method.AddExecution(abi.Condition(
		abi.Eq(IS_TESTNET, true),
		abi.EncodeJSON("faucet is not supported in mainnet"),
	))

	method.AddExecution(abi.WriteUniversal("balance", from, balance))

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

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "from",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "to",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amount",
		"type":         "string",
		"maxlength":    256,
		"requirements": true,
	}))

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
		abi.Gt(amount, "0"),
		"Amount must be greater than zero.",
	))

	method.AddExecution(abi.Condition(
		abi.Gt(fromBalance, abi.PreciseAdd(amount, SEND_FEE, 0)),
		"Balance is not enough.",
	))

	method.AddExecution(abi.WriteUniversal("balance", from, abi.PreciseSub(fromBalance, abi.PreciseAdd(amount, SEND_FEE, 0), 0)))
	method.AddExecution(abi.WriteUniversal("balance", to, abi.PreciseAdd(toBalance, amount, 0)))

	network_fee_reserve := abi.ReadUniversal("network_fee_reserve", ZERO_ADDRESS, "0")
	method.AddExecution(abi.WriteUniversal("network_fee_reserve", ZERO_ADDRESS, abi.PreciseAdd(network_fee_reserve, SEND_FEE, 0)))

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
		"maxlength":    16384,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "from",
		"type":         "string",
		"maxlength":    44,
		"requirements": true,
	}))

	from := abi.Param("from")
	code := abi.Param("code")
	decodedCode := abi.DecodeJSON(code)

	codeType := abi.Get(decodedCode, "t", nil)
	name := abi.Get(decodedCode, "n", nil)
	space := abi.Get(decodedCode, "s", nil)
	version := abi.Get(decodedCode, "v", nil)
	writer := abi.Get(decodedCode, "w", nil)

	spaceID := abi.SpaceID(writer, abi.Hash(space))
	contractID := abi.HashMany(spaceID, name)

	exists := abi.If(
		abi.Eq(codeType, "contract"),
		abi.ReadUniversal("contract", contractID, nil),
		abi.ReadUniversal("request", contractID, nil),
	)

	method.AddExecution(abi.Condition(
		abi.Eq(writer, from),
		abi.EncodeJSON("Writer must be the same as the from address"),
	))

	method.AddExecution(abi.Condition(
		abi.IsString(codeType),
		abi.EncodeJSON(abi.Concat("Invalid type: ", codeType)),
	))

	method.AddExecution(abi.Condition(
		abi.Or(
			abi.Eq(codeType, "contract"),
			abi.Eq(codeType, "request"),
		),
		abi.EncodeJSON("Type must be one of the following: contract, request"),
	))

	method.AddExecution(abi.Condition(
		abi.IsString(name),
		abi.EncodeJSON(abi.Concat("Invalid name: ", name)),
	))

	method.AddExecution(abi.Condition(
		abi.RegMatch("/^[A-Za-z_0-9]+$/", name),
		abi.EncodeJSON("The name must consist of A-Za-z_0-9"),
	))

	method.AddExecution(abi.Condition(
		abi.IsNumeric(version),
		abi.EncodeJSON(abi.Concat("Invalid version: ", version)),
	))

	method.AddExecution(abi.Condition(
		abi.IsString(space),
		abi.EncodeJSON("invalid space"),
	))

	fee := abi.PreciseMul(
		abi.AsString(abi.Len(code)),
		PUBLISH_FEE_PER_BYTE,
		"0",
	)

	fee = abi.Check(fee, "publish_fee")
	userBalance := abi.ReadUniversal("balance", from, "0")

	method.AddExecution(
		abi.Condition(
			abi.Gte(userBalance, fee),
			abi.EncodeJSON("Balance is not enough"),
		),
	)

	info := abi.ReadUniversal(spaceID, name, nil)
	info = abi.DecodeJSON(info)
	deployedVersion := abi.Get(info, "v", "0")

	method.AddExecution(abi.Condition(
		abi.Or(
			abi.Eq(exists, nil),
			abi.If(
				abi.Eq(codeType, "contract"),
				abi.Gt(version, deployedVersion),
				abi.If(
					abi.Eq(codeType, "request"),
					abi.Gt(version, deployedVersion),
					false,
				),
			),
		),
		abi.EncodeJSON("Only new versions of code can be registered"),
	))

	method.AddExecution(abi.If(
		abi.Eq(codeType, "contract"),
		abi.WriteUniversal("contract", contractID, from),
		abi.If(
			abi.Eq(codeType, "request"),
			abi.WriteUniversal("request", contractID, from),
			false,
		),
	))

	method.AddExecution(abi.WriteUniversal(spaceID, name, code))
	method.AddExecution(abi.WriteUniversal("balance", from, abi.PreciseSub(userBalance, fee, 0)))
	method.AddExecution(abi.WriteUniversal("network_fee_reserve", ZERO_ADDRESS, fee))

	difficulty := abi.ReadUniversal("network_difficulty", ZERO_ADDRESS, "0")
	difficulty = abi.PreciseSub(difficulty, "8", "0")
	method.AddExecution(abi.WriteUniversal("network_difficulty", ZERO_ADDRESS, difficulty))

	return method
}

func Count() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Count",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	count := abi.ReadUniversal("count", "count", "0")
	update := abi.PreciseAdd(count, "1", 0)
	wr := abi.WriteUniversal("count", "count", update)
	method.AddExecution(wr)

	return method
}
