package vm

import (
	. "hello/pkg/core/abi"
	"hello/pkg/util"
)

func OpAdd(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) == 2 {
		a := "0"
		b := "0"

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			a = str
		} else {
			DebugLog("OpAdd: first value is not numeric")
			return "0"
		}

		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			b = str
		} else {
			DebugLog("OpAdd: second value is not numeric")
			return "0"
		}

		result := util.Add(a, b, nil)
		DebugLog("OpAdd: result:", result)
		OperatorLog("OpAdd", "input:", vars, "result:", result)
		return result
	}

	DebugLog("OpAdd: vars is not array or invalid length")
	OperatorLog("OpAdd", "input:", vars, "result: 0")
	return "0"
}

func OpSub(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) == 2 {
		a := "0"
		b := "0"

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			a = str
		} else {
			DebugLog("OpSub: first value is not numeric")
			return "0"
		}

		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			b = str
		} else {
			DebugLog("OpSub: second value is not numeric")
			return "0"
		}

		result := util.Sub(a, b, nil)
		OperatorLog("OpSub", "input:", vars, "result:", result)
		return result
	}

	DebugLog("OpSub: vars is not array or invalid length")
	OperatorLog("OpSub", "input:", vars, "result: 0")
	return "0"
}

func OpMul(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) == 2 {
		a := "0"
		b := "0"

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			a = str
		} else {
			DebugLog("OpMul: first value is not numeric")
			return "0"
		}

		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			b = str
		} else {
			DebugLog("OpMul: second value is not numeric")
			return "0"
		}

		result := util.Mul(a, b, nil)
		OperatorLog("OpMul", "input:", vars, "result:", result)
		return result
	}

	DebugLog("OpMul: vars is not array or invalid length")
	OperatorLog("OpMul", "input:", vars, "result: 0")
	return "0"
}

func OpDiv(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) == 2 {
		a := "0"
		b := "0"

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			a = str
		} else {
			DebugLog("OpDiv: first value is not numeric")
			return "0"
		}

		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			b = str
		} else {
			DebugLog("OpDiv: second value is not numeric")
			return "0"
		}

		if divResult := util.Div(a, b, nil); divResult != nil {
			OperatorLog("OpDiv", "input:", vars, "result:", *divResult)
			return *divResult
		} else {
			DebugLog("OpDiv: division by zero or invalid division")
			OperatorLog("OpDiv", "input:", vars, "result: 0")
			return "0"
		}
	}

	DebugLog("OpDiv: vars is not array or invalid length")
	OperatorLog("OpDiv", "input:", vars, "result: 0")
	return "0"
}

func OpPreciseAdd(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && (len(arr) == 2 || len(arr) == 3) {
		a := "0"
		b := "0"
		scale := 0

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			a = str
		}
		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			b = str
		}
		if len(arr) == 3 {
			if val, ok := arr[2].(int); ok {
				scale = val
			}
		} else {
			scale = 0
		}

		if scale < 0 || scale > 10 {
			scale = 0
		}

		result := util.Add(a, b, &scale)
		OperatorLog("OpPreciseAdd", "input:", vars, "result:", result)
		return result
	}

	DebugLog("OpPreciseAdd: vars is not array or invalid length")
	OperatorLog("OpPreciseAdd", "input:", vars, "result: 0")
	return "0"
}

func OpPreciseSub(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && (len(arr) == 2 || len(arr) == 3) {
		a := "0"
		b := "0"
		scale := 0

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			a = str
		}
		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			b = str
		}
		if len(arr) == 3 {
			if val, ok := arr[2].(int); ok {
				scale = val
			}
		} else {
			scale = 0
		}

		if scale < 0 || scale > 10 {
			scale = 0
		}

		result := util.Sub(a, b, &scale)
		OperatorLog("OpPreciseSub", "input:", vars, "result:", result)
		return result
	}

	DebugLog("OpPreciseSub: vars is not array or invalid length")
	OperatorLog("OpPreciseSub", "input:", vars, "result: 0")
	return "0"
}

func OpPreciseMul(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && (len(arr) == 2 || len(arr) == 3) {
		a := "0"
		b := "0"
		scale := 0

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			a = str
		}
		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			b = str
		}
		if len(arr) == 3 {
			if val, ok := arr[2].(int); ok {
				scale = val
			}
		} else {
			scale = 0
		}

		if scale < 0 || scale > 10 {
			scale = 0
		}

		result := util.Mul(a, b, &scale)
		OperatorLog("OpPreciseMul", "input:", vars, "result:", result)
		return result
	}

	DebugLog("OpPreciseMul: vars is not array or invalid length")
	OperatorLog("OpPreciseMul", "input:", vars, "result: 0")
	return "0"
}

func OpPreciseDiv(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && (len(arr) == 2 || len(arr) == 3) {
		a := "0"
		b := "0"
		scale := 0

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			a = str
		}
		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			b = str
		}
		if len(arr) == 3 {
			if val, ok := arr[2].(int); ok {
				scale = val
			}
		} else {
			scale = 0
		}

		if scale < 0 || scale > 10 {
			scale = 0
		}

		if result := util.Div(a, b, &scale); result != nil {
			OperatorLog("OpPreciseDiv", "input:", vars, "result:", *result)
			return *result
		}
	}

	DebugLog("OpPreciseDiv: vars is not array or invalid length")
	OperatorLog("OpPreciseDiv", "input:", vars, "result: 0")
	return "0"
}

func OpScale(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) == 0 {
			OperatorLog("OpScale", "input:", vars, "result: 0")
			return 0
		}

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			result := util.GetScale(str)
			OperatorLog("OpScale", "input:", vars, "result:", result)
			return result
		}
	}

	DebugLog("OpScale: vars is not array or invalid length")
	OperatorLog("OpScale", "input:", vars, "result: 0")
	return 0
}
