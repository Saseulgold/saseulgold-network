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
		"name":      "target",
		"type":      "string",
		"maxlength": TIME_HASH_SIZE,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":      "full",
		"type":      "boolean",
		"maxlength": 5,
		"default":   false,
	}))

	target := abi.Param("target")
	full := abi.Param("full")
	method.AddExecution(abi.Response(map[string]interface{}{
		"$get_block": []interface{}{target, full},
	}))

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

	page := abi.Param("page")
	count := abi.Param("count")
	sort := abi.Param("sort")

	method.AddExecution(abi.Condition(
		abi.Lte(count, 100),
		"The parameter 'count' must be less than or equal to 100.",
	))

	method.AddExecution(abi.Response(map[string]interface{}{
		"$list_block": []interface{}{page, count, sort},
	}))

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
