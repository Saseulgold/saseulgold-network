package vm

import (
	. "hello/pkg/core/abi"
	D "hello/pkg/core/debug"
	"reflect"
)

func OpCondition(i *Interpreter, vars interface{}) interface{} {
	if i.state != StateCondition {
		OperatorLog("OpCondition", "input:", vars, "result: true")
		return true
	}

	var tf bool
	var errMsg string

	if arr, ok := vars.([]interface{}); ok {
		if len(arr) > 0 {
			if v, ok := arr[0].(bool); ok {
				tf = v
			}
		}

		if len(arr) > 1 {
			if v, ok := arr[1].(string); ok {
				errMsg = v
			}
		}
	}

	if !tf {
		i.breakFlag = true
		if errMsg != "" {
			i.result = errMsg
		}
		OperatorLog("OpCondition", "input:", vars, "result:", errMsg)
		return errMsg
	}

	OperatorLog("OpCondition", "input:", vars, "result: true")
	return true
}

func OpResponse(i *Interpreter, vars interface{}) interface{} {
	if i.state != StateExecution {
		var result ABI
		if arr, ok := vars.([]interface{}); ok {
			result = ABI{
				Key:   "$response",
				Value: arr,
			}
		} else {
			result = ABI{
				Key:   "$response",
				Value: []interface{}{},
			}
		}
		OperatorLog("OpResponse", "input:", vars, "result:", result)
		return result
	}

	i.breakFlag = true
	if arr, ok := vars.([]interface{}); ok && len(arr) > 0 {
		i.result = arr[0]
	} else {
		i.result = nil
	}
	OperatorLog("OpResponse", "input:", vars, "result: nil")
	return nil
}

func OpWeight(i *Interpreter, vars interface{}) interface{} {
	OperatorLog("OpWeight", "input:", vars, "result:", i.weight)
	return i.weight
}

func OpIf(i *Interpreter, vars interface{}) interface{} {
	var condition bool
	var trueVal, falseVal interface{}

	if arr, ok := vars.([]interface{}); ok {
		if len(arr) > 0 {
			if v, ok := arr[0].(bool); ok {
				condition = v
			}
		}

		if len(arr) > 1 {
			trueVal = arr[1]
		}

		if len(arr) > 2 {
			falseVal = arr[2]
		}
	}

	var result interface{}
	if condition {
		result = trueVal
	} else {
		result = falseVal
	}
	OperatorLog("OpIf", "input:", vars, "result:", result)
	return result
}

func OpAnd(i *Interpreter, vars interface{}) interface{} {
	var result *bool

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if boolVal, ok := v.(bool); !ok {
				OperatorLog("OpAnd", "input:", vars, "result: false")
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
		OperatorLog("OpAnd", "input:", vars, "result: false")
		return false
	}
	OperatorLog("OpAnd", "input:", vars, "result:", *result)
	return *result
}

func OpOr(i *Interpreter, vars interface{}) interface{} {
	var result *bool

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if boolVal, ok := v.(bool); !ok {
				OperatorLog("OpOr", "input:", vars, "result: false")
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
		OperatorLog("OpOr", "input:", vars, "result: false")
		return false
	}
	OperatorLog("OpOr", "input:", vars, "result:", *result)
	return *result
}

func OpGet(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) < 2 {
			OperatorLog("OpGet1", "input:", vars, "type:", reflect.TypeOf(arr[0]), "result: nil")
			return nil
		}

		var obj map[string]interface{}
		switch v := arr[0].(type) {
		case map[string]interface{}:
			obj = v
		default:
			DebugLog("OpGet2", "input:", vars, "type:", reflect.TypeOf(v))
			D.DebugPanic("OpGet2", "input:", vars, "type:", reflect.TypeOf(v), "result: nil")
			return nil
		}

		var key string
		switch v := arr[1].(type) {
		case string:
			key = v
		default:
			D.DebugPanic("OpGet4", "input:", vars, "type:", reflect.TypeOf(arr[1]), "result: nil")
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

	DebugLog("OpGet", "input:", vars, "result: nil")
	D.DebugPanic("OpGet", "input:", vars, "result: nil")
	return nil
}

func OpIn(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) < 2 {
			OperatorLog("OpIn", "input:", vars, "result: false")
			return false
		}

		target := arr[0]
		cases, ok := arr[1].([]interface{})
		if !ok {
			OperatorLog("OpIn", "input:", vars, "result: false")
			return false
		}

		for _, v := range cases {
			if reflect.DeepEqual(target, v) {
				OperatorLog("OpIn", "input:", vars, "result: true")
				return true
			}
		}
	}
	OperatorLog("OpIn", "input:", vars, "result: false")
	return false
}
