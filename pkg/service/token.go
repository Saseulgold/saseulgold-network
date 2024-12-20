package service

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
)

func Mint(writer string, space string) *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Mint",
		"version": "1",
		"space":   space,
		"writer":  writer,
	})

	// Add parameters
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "name",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "symbol",
		"type":         "string",
		"maxlength":    20,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amount",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	// Add executions
	method.AddExecution(abi.Condition(
		abi.Eq(abi.ReadUniversal("info", "00", nil), nil),
		"The token can only be issued once.",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(writer, abi.Param("from")),
		"You are not the contract writer.",
	))

	method.AddExecution(abi.Condition(
		abi.Gt(abi.Param("amount"), "0"),
		"The amount must be greater than 0.",
	))

	method.AddExecution(abi.Condition(
		abi.Gte(abi.Param("decimal"), 0),
		"The decimal must be greater than or equal to 0.",
	))

	// Save token info
	method.AddExecution(abi.WriteUniversal("info", "00", map[string]interface{}{
		"name":         abi.Param("name"),
		"symbol":       abi.Param("symbol"),
		"total_supply": abi.Param("amount"),
		"decimal":      abi.Param("decimal"),
	}))

	// Set initial balance
	method.AddExecution(abi.WriteUniversal("balance", abi.Param("from"), abi.Param("amount")))

	return method
}

func GetInfo(writer string, space string) *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "GetInfo",
		"version": "1",
		"space":   space,
		"writer":  writer,
	})

	info := abi.ReadUniversal("info", "00", nil)

	method.AddExecution(abi.Condition(
		abi.Ne(info, nil),
		"The token has not been issued yet.",
	))

	method.AddExecution(abi.Response(info))

	return method
}

func Send(writer string, space string) *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Send",
		"version": "1",
		"space":   space,
		"writer":  writer,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "to",
		"type":         "string",
		"maxlength":    64, // ID_HASH_SIZE
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amount",
		"type":         "string",
		"maxlength":    40,
		"requirements": true,
	}))

	from := abi.Param("from")
	to := abi.Param("to")
	amount := abi.Param("amount")

	method.AddExecution(abi.Condition(
		abi.Ne(from, to),
		"You can't send to yourself.",
	))

	method.AddExecution(abi.Condition(
		abi.Gt(amount, "0"),
		"The amount must be greater than 0.",
	))

	method.AddExecution(abi.Condition(
		abi.Gte(abi.ReadUniversal("balance", from, "0"), amount),
		"You can't send more than what you have.",
	))

	method.AddExecution(abi.WriteUniversal("balance", from,
		abi.Sub(abi.ReadUniversal("balance", from, "0"), amount)))

	method.AddExecution(abi.WriteUniversal("balance", to,
		abi.Add(abi.ReadUniversal("balance", to, "0"), amount)))

	return method
}

func GetBalance(writer string, space string) *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "GetBalance",
		"version": "1",
		"space":   space,
		"writer":  writer,
	})

	// Add address parameter
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "address",
		"type":         "string",
		"maxlength":    64, // ID_HASH_SIZE
		"requirements": true,
	}))

	// Get balance and return response
	method.AddExecution(abi.Response(map[string]interface{}{
		"balance": abi.ReadUniversal("balance", abi.Param("address"), "0"),
	}))

	return method
}
