package vm

import (
	. "hello/pkg/core/abi"
	"hello/pkg/util"
)

func OpAdd(i *Interpreter, vars interface{}) interface{} {
	a, b := Unpack2(vars)

	aStr := "0"
	bStr := "0"

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	} else {
		DebugLog("OpAdd: first value is not numeric")
		return "0"
	}

	if str, ok := b.(string); ok && util.IsNumeric(str) {
		bStr = str
	} else {
		DebugLog("OpAdd: second value is not numeric")
		return "0"
	}

	result := util.Add(aStr, bStr, nil)
	OperatorLog("OpAdd", "input:", vars, "result:", result)
	return result
}

func OpSub(i *Interpreter, vars interface{}) interface{} {
	a, b := Unpack2(vars)

	aStr := "0"
	bStr := "0"

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	} else {
		DebugLog("OpSub: first value is not numeric")
		return "0"
	}

	if str, ok := b.(string); ok && util.IsNumeric(str) {
		bStr = str
	} else {
		DebugLog("OpSub: second value is not numeric")
		return "0"
	}

	result := util.Sub(aStr, bStr, nil)
	OperatorLog("OpSub", "input:", vars, "result:", result)
	return result
}

func OpMul(i *Interpreter, vars interface{}) interface{} {
	a, b := Unpack2(vars)

	aStr := "0"
	bStr := "0"

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	} else {
		DebugLog("OpMul: first value is not numeric")
		return "0"
	}

	if str, ok := b.(string); ok && util.IsNumeric(str) {
		bStr = str
	} else {
		DebugLog("OpMul: second value is not numeric")
		return "0"
	}

	result := util.Mul(aStr, bStr, nil)
	OperatorLog("OpMul", "input:", vars, "result:", result)
	return result
}

func OpDiv(i *Interpreter, vars interface{}) interface{} {
	a, b := Unpack2(vars)

	aStr := "0"
	bStr := "0"

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	} else {
		DebugLog("OpDiv: first value is not numeric")
		return "0"
	}

	if str, ok := b.(string); ok && util.IsNumeric(str) {
		bStr = str
	} else {
		DebugLog("OpDiv: second value is not numeric")
		return "0"
	}

	if divResult := util.Div(aStr, bStr, nil); divResult != nil {
		OperatorLog("OpDiv", "input:", vars, "result:", *divResult)
		return *divResult
	} else {
		DebugLog("OpDiv: division by zero or invalid division")
		OperatorLog("OpDiv", "input:", vars, "result: 0")
		return "0"
	}
}

func OpPreciseAdd(i *Interpreter, vars interface{}) interface{} {
	a, b, scaleVal := Unpack2Or3(vars)

	aStr := "0"
	bStr := "0"
	scale := 0

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	}
	if str, ok := b.(string); ok && util.IsNumeric(str) {
		bStr = str
	}
	if scaleVal != nil {
		if val, ok := scaleVal.(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	result := util.Add(aStr, bStr, &scale)
	OperatorLog("OpPreciseAdd", "input:", vars, "result:", result)
	return result
}

func OpPreciseSub(i *Interpreter, vars interface{}) interface{} {
	a, b, scaleVal := Unpack2Or3(vars)

	aStr := "0"
	bStr := "0"
	scale := 0

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	}
	if str, ok := b.(string); ok && util.IsNumeric(str) {
		bStr = str
	}
	if scaleVal != nil {
		if val, ok := scaleVal.(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	result := util.Sub(aStr, bStr, &scale)
	OperatorLog("OpPreciseSub", "input:", vars, "result:", result)
	return result
}

func OpPreciseMul(i *Interpreter, vars interface{}) interface{} {
	a, b, scaleVal := Unpack2Or3(vars)

	aStr := "0"
	bStr := "0"
	scale := 0

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	}
	if str, ok := b.(string); ok && util.IsNumeric(str) {
		bStr = str
	}
	if scaleVal != nil {
		if val, ok := scaleVal.(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	result := util.Mul(aStr, bStr, &scale)
	OperatorLog("OpPreciseMul", "input:", vars, "result:", result)
	return result
}

func OpPreciseSqrt(i *Interpreter, vars interface{}) interface{} {
	a, scaleVal := Unpack1Or2(vars)

	aStr := "0"
	scale := 0

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	}
	if scaleVal != nil {
		if val, ok := scaleVal.(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	result := util.Sqrt(aStr, &scale)
	OperatorLog("OpPreciseSqrt", "input:", vars, "result:", result)
	return result
}

func OpPreciseDiv(i *Interpreter, vars interface{}) interface{} {
	a, b, scaleVal := Unpack2Or3(vars)

	aStr := "0"
	bStr := "0"
	scale := 0

	if str, ok := a.(string); ok && util.IsNumeric(str) {
		aStr = str
	}
	if str, ok := b.(string); ok && util.IsNumeric(str) {
		bStr = str
	}
	if scaleVal != nil {
		if val, ok := scaleVal.(int); ok {
			scale = val
		}
	}

	if scale < 0 || scale > 10 {
		scale = 0
	}

	if result := util.Div(aStr, bStr, &scale); result != nil {
		OperatorLog("OpPreciseDiv", "input:", vars, "result:", *result)
		return *result
	}

	DebugLog("OpPreciseDiv: vars is not array or invalid length")
	OperatorLog("OpPreciseDiv", "input:", vars, "result: 0")
	return "0"
}

func OpScale(i *Interpreter, vars interface{}) interface{} {
	value := Unpack1(vars)

	if str, ok := value.(string); ok && util.IsNumeric(str) {
		result := util.GetScale(str)
		OperatorLog("OpScale", "input:", vars, "result:", result)
		return result
	}

	DebugLog("OpScale: value is not numeric")
	OperatorLog("OpScale", "input:", vars, "result: 0")
	return 0
}
