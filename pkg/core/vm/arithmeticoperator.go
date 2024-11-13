package vm

import (
	. "hello/pkg/core/abi"
	"hello/pkg/util"
)

func OpAdd(i *Interpreter, vars interface{}) interface{} {
	result := "0"

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			DebugLog("OpAdd: value:", v)
			val := "0"
			if str, ok := v.(string); ok && util.IsNumeric(str) {
				val = str
			} else {
				DebugLog("OpAdd: value is not numeric")
			}

			result = util.Add(result, val, nil)
			DebugLog("OpAdd result:", result, val)
		}
	} else {
		DebugLog("OpAdd: vars is not array type")
		return "0"
	}

	return result
}

func OpSub(i *Interpreter, vars []interface{}) interface{} {
	var result string

	for _, v := range vars {
		val := "0"
		if str, ok := v.(string); ok && util.IsNumeric(str) {
			val = str
		}

		if result == "" {
			result = val
		} else {
			result = util.Sub(result, val, nil)
		}
	}

	if result == "" {
		return "0"
	}

	return result
}

func OpMul(i *Interpreter, vars []interface{}) interface{} {
	var result string

	for _, v := range vars {
		val := "0"
		if str, ok := v.(string); ok && util.IsNumeric(str) {
			val = str
		}

		if result == "" {
			result = val
		} else {
			result = util.Mul(result, val, nil)
		}
	}

	if result == "" {
		return "0"
	}

	return result
}

func OpDiv(i *Interpreter, vars []interface{}) interface{} {
	var result string

	for _, v := range vars {
		val := "0"
		if str, ok := v.(string); ok && util.IsNumeric(str) {
			val = str
		}

		if result == "" {
			result = val
		} else {
			if divResult := util.Div(result, val, nil); divResult != nil {
				result = *divResult
			} else {
				result = "0"
			}
		}
	}

	if result == "" {
		return "0"
	}

	return result
}

func OpPreciseAdd(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return "0"
	}

	left := "0"
	right := "0"
	scale := 0

	if str, ok := vars[0].(string); ok && util.IsNumeric(str) {
		left = str
	}
	if str, ok := vars[1].(string); ok && util.IsNumeric(str) {
		right = str
	}
	if len(vars) > 2 {
		if val, ok := vars[2].(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	return util.Add(left, right, &scale)
}

func OpPreciseSub(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return "0"
	}

	left := "0"
	right := "0"
	scale := 0

	if str, ok := vars[0].(string); ok && util.IsNumeric(str) {
		left = str
	}
	if str, ok := vars[1].(string); ok && util.IsNumeric(str) {
		right = str
	}
	if len(vars) > 2 {
		if val, ok := vars[2].(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	return util.Sub(left, right, &scale)
}

func OpPreciseMul(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return "0"
	}

	left := "0"
	right := "0"
	scale := 0

	if str, ok := vars[0].(string); ok && util.IsNumeric(str) {
		left = str
	}
	if str, ok := vars[1].(string); ok && util.IsNumeric(str) {
		right = str
	}
	if len(vars) > 2 {
		if val, ok := vars[2].(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	return util.Mul(left, right, &scale)
}

func OpPreciseDiv(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return "0"
	}

	left := "0"
	right := "0"
	scale := 0

	if str, ok := vars[0].(string); ok && util.IsNumeric(str) {
		left = str
	}
	if str, ok := vars[1].(string); ok && util.IsNumeric(str) {
		right = str
	}
	if len(vars) > 2 {
		if val, ok := vars[2].(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	if result := util.Div(left, right, &scale); result != nil {
		return *result
	}
	return "0"
}

func OpScale(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return 0
	}

	if str, ok := vars[0].(string); ok && util.IsNumeric(str) {
		return util.GetScale(str)
	}

	return 0
}
