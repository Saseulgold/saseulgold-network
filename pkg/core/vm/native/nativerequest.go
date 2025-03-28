package native

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/util"
)

func GetBlock() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "GetBlock",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name": "target",
		"type": "int",
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":      "responseType",
		"type":      "string",
		"maxlength": 5,
		"default":   "full",
	}))

	target := abi.Param("target")
	responseType := abi.Param("responseType")

	block := abi.GetBlock(target, responseType)
	response := abi.EncodeJSON(block)

	method.AddExecution(abi.Response(response))

	return method
}

func ListBlock() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "ListBlock",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":      "page",
		"type":      "int",
		"maxlength": 16,
		"default":   1,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":      "count",
		"type":      "int",
		"maxlength": 4,
		"default":   20,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":      "sort",
		"type":      "int",
		"maxlength": 2,
		"default":   -1,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":      "responseType",
		"type":      "string",
		"maxlength": 5,
		"default":   "base",
	}))

	page := abi.Param("page")
	count := abi.Param("count")
	responseType := abi.Param("responseType")

	method.AddExecution(abi.Condition(
		abi.Lte(count, 100),
		abi.EncodeJSON("The parameter 'count' must be less than or equal to 100."),
	))

	blocks := abi.ListBlock(page, count, responseType)
	response := abi.EncodeJSON(blocks)

	method.AddExecution(abi.Response(response))

	return method
}

func GetBalance() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "GetBalance",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "address",
		"type":         "string",
		"maxlength":    ID_HASH_SIZE,
		"requirements": true,
	}))

	address := abi.Param("address")
	balance := abi.ReadUniversal("balance", address, "0")
	method.AddExecution(abi.Response(balance))

	return method
}

func GetTokenInfo() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "GetTokenInfo",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))

	token_address := abi.Param("token_address")

	var response interface{}

	supply_univ := abi.ReadUniversal(token_address, "supply", nil)
	owner_univ := abi.ReadUniversal(token_address, "owner", nil)
	symbol_univ := abi.ReadUniversal(token_address, "symbol", nil)
	name_univ := abi.ReadUniversal(token_address, "name", nil)
	icon_url_univ := abi.ReadUniversal(token_address, "icon_url", nil)

	response = abi.Set(response, "token_address", token_address)
	response = abi.Set(response, "owner", owner_univ)
	response = abi.Set(response, "symbol", symbol_univ)
	response = abi.Set(response, "supply", supply_univ)
	response = abi.Set(response, "name", name_univ)
	response = abi.Set(response, "icon_url", icon_url_univ)

	response = abi.EncodeJSON(response)
	method.AddExecution(abi.Response(response))

	return method
}

func ListTransaction() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "ListTransaction",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":      "page",
		"type":      "int",
		"maxlength": 16,
		"default":   1,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":      "count",
		"type":      "int",
		"maxlength": 4,
		"default":   20,
	}))

	count := abi.Param("count")
	response := abi.ListTransaction(count)
	response = abi.EncodeJSON(response)

	method.AddExecution(abi.Response(response))

	return method
}
