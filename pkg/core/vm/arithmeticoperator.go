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
		return result
	}

	DebugLog("OpAdd: vars is not array or invalid length")
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

		return util.Sub(a, b, nil)
	}

	DebugLog("OpSub: vars is not array or invalid length")
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

		return util.Mul(a, b, nil)
	}

	DebugLog("OpMul: vars is not array or invalid length")
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
			return *divResult
		} else {
			DebugLog("OpDiv: division by zero or invalid division")
			return "0"
		}
	}

	DebugLog("OpDiv: vars is not array or invalid length")
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

		return util.Add(a, b, &scale)
	}

	DebugLog("OpPreciseAdd: vars is not array or invalid length")
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

		return util.Sub(a, b, &scale)
	}

	DebugLog("OpPreciseSub: vars is not array or invalid length")
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

		return util.Mul(a, b, &scale)
	}

	DebugLog("OpPreciseMul: vars is not array or invalid length")
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
			return *result
		}
	}

	DebugLog("OpPreciseDiv: vars is not array or invalid length")
	return "0"
}

func OpScale(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) == 0 {
			return 0
		}

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			return util.GetScale(str)
		}
	}

	DebugLog("OpScale: vars is not array or invalid length")
	return 0
}
