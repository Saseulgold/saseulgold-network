package vm

import (
	"fmt"
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
		return *divResult
	} else {
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

	if result == nil {
		return "0"
	}
	fmt.Println("OpPreciseSqrt", "input:", vars, "result:", *result)
	return *result
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
		return *result
	}

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

func OpMax(i *Interpreter, vars interface{}) interface{} {
	a, b := Unpack2(vars)

	// Check if both are strings
	if aStr, aOk := a.(string); aOk {
		if bStr, bOk := b.(string); bOk {
			// If both are numeric strings, compare as numbers
			if util.IsNumeric(aStr) && util.IsNumeric(bStr) {
				if result := util.Compare(aStr, bStr, 0); result >= 0 {
					fmt.Println("OpMax: ", aStr)
					return aStr
				}
				fmt.Println("OpMax: ", bStr)
				return bStr
			}
			// If both are regular strings, compare lexicographically
			if aStr >= bStr {
				fmt.Println("OpMax: ", aStr)
				return aStr
			}
			fmt.Println("OpMax: ", bStr)
			return bStr
		}
	}

	// Check if both are numbers
	if aNum, aOk := a.(float64); aOk {
		if bNum, bOk := b.(float64); bOk {
			if aNum >= bNum {
				return aNum
			}
			return bNum
		}
	}

	DebugLog("OpMax: incompatible types or invalid values")
	return nil
}

func OpMin(i *Interpreter, vars interface{}) interface{} {
	a, b := Unpack2(vars)

	// Check if both are strings
	if aStr, aOk := a.(string); aOk {
		if bStr, bOk := b.(string); bOk {
			// If both are numeric strings, compare as numbers
			if util.IsNumeric(aStr) && util.IsNumeric(bStr) {
				if result := util.Compare(aStr, bStr, 0); result <= 0 {
					return aStr
				}
				fmt.Println("OpMin: ", bStr)
				return bStr
			}
			// If both are regular strings, compare lexicographically
			if aStr <= bStr {
				fmt.Println("OpMin: ", aStr)
				return aStr
			}
			fmt.Println("OpMin: ", bStr)
			return bStr
		}
	}

	// Check if both are numbers
	if aNum, aOk := a.(float64); aOk {
		if bNum, bOk := b.(float64); bOk {
			if aNum <= bNum {
				return aNum
			}
			return bNum
		}
	}

	DebugLog("OpMin: incompatible types or invalid values")
	return nil
}
