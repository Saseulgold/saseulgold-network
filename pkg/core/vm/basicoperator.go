package vm

import (
	. "hello/pkg/core/abi"
	"reflect"
)

func OpCondition(i *Interpreter, vars interface{}) interface{} {
	if i.state != StateCondition {
		return true
	}

	var abi bool
	var errMsg string

	if arr, ok := vars.([]interface{}); ok {
		if len(arr) > 0 {
			if v, ok := arr[0].(bool); ok {
				abi = v
			}
		}

		if len(arr) > 1 {
			if v, ok := arr[1].(string); ok {
				errMsg = v
			}
		}
	}

	if !abi {
		i.breakFlag = true
		if errMsg != "" {
			i.result = errMsg
		}
		return false
	}

	return true
}

func OpResponse(i *Interpreter, vars interface{}) interface{} {
	if i.state != StateExecution {
		if arr, ok := vars.([]interface{}); ok {
			return map[string][]interface{}{
				"$response": arr,
			}
		}
		return map[string][]interface{}{
			"$response": []interface{}{},
		}
	}

	i.breakFlag = true
	if arr, ok := vars.([]interface{}); ok && len(arr) > 0 {
		i.result = arr[0]
	} else {
		i.result = nil
	}
	return nil
}

func OpWeight(i *Interpreter, vars interface{}) interface{} {
	return i.weight
}

func OpIf(i *Interpreter, vars interface{}) interface{} {
	var condition bool
	var trueVal, falseVal interface{}

	if arr, ok := vars.([]interface{}); ok {
		if len(arr) > 0 {
			if v, ok := arr[0].(bool); ok {
				condition = v
			} else {
				DebugLog("OpIf: condition is not boolean")
			}
		}

		if len(arr) > 1 {
			trueVal = arr[1]
		}

		if len(arr) > 2 {
			falseVal = arr[2]
		}
	}
	DebugLog("OpIf: condition =", condition)
	DebugLog("OpIf: trueVal =", trueVal)
	DebugLog("OpIf: falseVal =", falseVal)
	if condition {
		return trueVal
	}
	return falseVal
}

func OpAnd(i *Interpreter, vars interface{}) interface{} {
	var result *bool

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if boolVal, ok := v.(bool); !ok {
				return false
			} else {
				if result == nil {
					result = &boolVal
				} else {
					newVal := *result && boolVal
					result = &newVal
				}
			}
		}
	}

	if result == nil {
		return false
	}
	return *result
}

func OpOr(i *Interpreter, vars interface{}) interface{} {
	var result *bool

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if boolVal, ok := v.(bool); !ok {
				return false
			} else {
				if result == nil {
					result = &boolVal
				} else {
					newVal := *result || boolVal
					result = &newVal
				}
			}
		}
	}

	if result == nil {
		return false
	}
	return *result
}

func OpGet(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) < 2 {
			return nil
		}

		obj, ok := arr[0].(map[string]interface{})
		if !ok {
			return nil
		}

		var key string
		switch v := arr[1].(type) {
		case string:
			key = v
		case float64:
			key = reflect.TypeOf(v).String()
		default:
			return nil
		}

		var defaultVal interface{}
		if len(arr) > 2 {
			defaultVal = arr[2]
		}

		if val, exists := obj[key]; exists {
			return val
		}
		return defaultVal
	}
	return nil
}

func OpIn(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) < 2 {
			return false
		}

		target := arr[0]
		cases, ok := arr[1].([]interface{})
		if !ok {
			return false
		}

		for _, v := range cases {
			if reflect.DeepEqual(target, v) {
				return true
			}
		}
	}
	return false
}
