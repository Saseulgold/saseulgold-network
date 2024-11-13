package vm

import (
	. "hello/pkg/core/abi"
	"hello/pkg/util"
)

func OpAdd(i *Interpreter, vars interface{}) interface{} {
	result := "0"

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			val := "0"
			if str, ok := v.(string); ok && util.IsNumeric(str) {
				val = str
			} else {
				DebugLog("OpAdd: value is not numeric")
			}

			result = util.Add(result, val, nil)
		}
	} else {
		return "0"
	}
	DebugLog("OpAdd: result:", result)
	return result
}

func OpSub(i *Interpreter, vars interface{}) interface{} {
	result := "0"

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			val := "0"
			if str, ok := v.(string); ok && util.IsNumeric(str) {
				val = str
			} else {
				DebugLog("OpSub: value is not numeric")
			}

			result = util.Sub(result, val, nil)
		}
	} else {
		return "0"
	}

	return result
}

func OpMul(i *Interpreter, vars interface{}) interface{} {
	result := "1"

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			val := "0"
			if str, ok := v.(string); ok && util.IsNumeric(str) {
				val = str
			} else {
				DebugLog("OpMul: value is not numeric")
			}
			result = util.Mul(result, val, nil)
		}
	} else {
		DebugLog("OpMul: vars is not array")
		return "0"
	}

	return result
}

func OpDiv(i *Interpreter, vars interface{}) interface{} {
	result := "0"

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			val := "0"
			if str, ok := v.(string); ok && util.IsNumeric(str) {
				val = str
			} else {
				DebugLog("OpDiv: value is not numeric")
			}

			if divResult := util.Div(result, val, nil); divResult != nil {
				result = *divResult
			} else {
				result = "0"
			}
		}
	} else {
		return "0"
	}

	return result
}

func OpPreciseAdd(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) < 2 {
			return "0"
		}

		left := "0"
		right := "0"
		scale := 0

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			left = str
		}
		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			right = str
		}
		if len(arr) > 2 {
			if val, ok := arr[2].(int); ok {
				scale = val
			}
		}

		if scale < 0 || scale > 10 {
			scale = 0
		}

		return util.Add(left, right, &scale)
	}

	return "0"
}

func OpPreciseSub(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) < 2 {
			return "0"
		}

		left := "0"
		right := "0"
		scale := 0

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			left = str
		}
		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			right = str
		}
		if len(arr) > 2 {
			if val, ok := arr[2].(int); ok {
				scale = val
			}
		}

		if scale < 0 || scale > 10 {
			scale = 0
		}

		return util.Sub(left, right, &scale)
	}

	return "0"
}

func OpPreciseMul(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) < 2 {
			return "0"
		}

		left := "0"
		right := "0"
		scale := 0

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			left = str
		}
		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			right = str
		}
		if len(arr) > 2 {
			if val, ok := arr[2].(int); ok {
				scale = val
			}
		}

		if scale < 0 || scale > 10 {
			scale = 0
		}

		return util.Mul(left, right, &scale)
	}

	return "0"
}

func OpPreciseDiv(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		if len(arr) < 2 {
			return "0"
		}

		left := "0"
		right := "0"
		scale := 0

		if str, ok := arr[0].(string); ok && util.IsNumeric(str) {
			left = str
		}
		if str, ok := arr[1].(string); ok && util.IsNumeric(str) {
			right = str
		}
		if len(arr) > 2 {
			if val, ok := arr[2].(int); ok {
				scale = val
			}
		}

		if scale < 0 || scale > 10 {
			scale = 0
		}

		if result := util.Div(left, right, &scale); result != nil {
			return *result
		}
	}

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

	return 0
}
