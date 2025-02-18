package vm

import (
	"fmt"
	. "hello/pkg/core/abi"
	D "hello/pkg/core/debug"
	"reflect"
)

func Unpack1(vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) != 1 {
		panic("OpUnpack1: vars is not an array or has 1 element")
	}

	return arr[0]
}

func Unpack1Or2(vars interface{}) (interface{}, interface{}) {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 1 || len(arr) > 2 {
		panic("OpUnpack1Or2: vars is not an array or has 2 elements")
	}

	var last interface{} = nil
	if len(arr) == 2 {
		last = arr[1]
	}

	return arr[0], last
}

func Unpack2(vars interface{}) (interface{}, interface{}) {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) != 2 {
		panic("OpUnpack2: vars is not an array or has 2 elements")
	}

	return arr[0], arr[1]
}

func Unpack3(vars interface{}) (interface{}, interface{}, interface{}) {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) != 3 {
		panic("OpUnpack3: vars is not an array or has 3 elements")
	}

	return arr[0], arr[1], arr[2]
}

func Unpack2Or3(vars interface{}) (interface{}, interface{}, interface{}) {
	arr, ok := vars.([]interface{})
	var last interface{} = nil

	if !ok || len(arr) < 2 || len(arr) > 3 {
		panic("OpUnpack2Or3: vars is not an array or has 4 elements")
	}

	if len(arr) == 3 {
		last = arr[2]
	}

	return arr[0], arr[1], last
}

func OpCondition(i *Interpreter, vars interface{}) interface{} {
	if i.state != StateCondition {
		OperatorLog("OpCondition", "input:", vars, "result: true")
		return true
	}

	tf, errMsg := Unpack2(vars)
	boolVal, _ := tf.(bool)
	errStr, _ := errMsg.(string)

	if !boolVal {
		i.breakFlag = true
		if errStr != "" {
			i.result = errStr
		}
		OperatorLog("OpCondition", "input:", vars, "result:", errStr)
		return errStr
	}

	OperatorLog("OpCondition", "input:", vars, "result: true")
	return true
}

func OpResponse(i *Interpreter, vars interface{}) interface{} {
	if i.state != StateExecution {
		response := Unpack1(vars)
		result := ABI{
			Key:   "$response",
			Value: []interface{}{response},
		}
		OperatorLog("OpResponse", "input:", vars, "result:", result)
		return result
	}

	i.breakFlag = true
	response := Unpack1(vars)
	i.result = response
	OperatorLog("OpResponse", "input:", vars, "result: nil")
	return nil
}

func OpWeight(i *Interpreter, vars interface{}) interface{} {
	OperatorLog("OpWeight", "input:", vars, "result:", i.weight)
	return i.weight
}

func OpIf(i *Interpreter, vars interface{}) interface{} {
	condition, trueVal, falseVal := Unpack3(vars)
	boolVal, _ := condition.(bool)

	var result interface{}
	if boolVal {
		result = trueVal
	} else {
		result = falseVal
	}

	return result
}

func OpAnd(i *Interpreter, vars interface{}) interface{} {
	values := vars.([]interface{})
	result := true

	for _, v := range values {
		if boolVal, ok := v.(bool); !ok || !boolVal {
			OperatorLog("OpAnd", "input:", vars, "result: false")
			return false
		}
	}

	OperatorLog("OpAnd", "input:", vars, "result:", result)
	return result
}

func OpOr(i *Interpreter, vars interface{}) interface{} {
	values := vars.([]interface{})
	result := false

	for _, v := range values {
		if boolVal, ok := v.(bool); ok && boolVal {
			result = true
			break
		}
	}

	OperatorLog("OpOr", "input:", vars, "result:", result)
	return result
}

func OpGet(i *Interpreter, vars interface{}) interface{} {
	obj, key, defaultVal := Unpack2Or3(vars)

	if obj == nil {
		return defaultVal
	}

	objMap, ok := obj.(map[string]interface{})

	if !ok {
		D.DebugPanic("OpGet", "input:", vars, "type:", reflect.TypeOf(obj), "result: nil")
		return nil
	}

	keyStr, ok := key.(string)
	if !ok {
		D.DebugPanic("OpGet", "input:", vars, "type:", reflect.TypeOf(key), "result: nil")
		return nil
	}

	if val, exists := objMap[keyStr]; exists {
		return val
	}
	return defaultVal
}

func OpIn(i *Interpreter, vars interface{}) interface{} {
	target, cases := Unpack2(vars)

	caseArray, ok := cases.([]interface{})
	if !ok {
		OperatorLog("OpIn", "input:", vars, "result: false")
		return false
	}

	for _, v := range caseArray {
		if reflect.DeepEqual(target, v) {
			OperatorLog("OpIn", "input:", vars, "result: true")
			return true
		}
	}

	OperatorLog("OpIn", "input:", vars, "result: false")
	return false
}

func OpCheck(i *Interpreter, vars interface{}) interface{} {
	v, k := Unpack1Or2(vars)

	fmt.Println("OpCheck value:", k, v, "type:", reflect.TypeOf(v))

	return vars.([]interface{})[0]
}
