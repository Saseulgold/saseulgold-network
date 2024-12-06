package vm

import (
	"fmt"
	D "hello/pkg/core/debug"
	"reflect"
)

func OperatorLog(args ...interface{}) interface{} {
	fmt.Println(args...)
	return true
}

func OpGetType(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) > 0 {
		result := reflect.TypeOf(arr[0]).String()
		OperatorLog("OpGetType", "input:", vars, "result:", result)
		return result
	}
	OperatorLog("OpGetType", "input:", vars, "result: nil")
	return "nil"
}

func OpIsNumeric(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			switch v.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				continue
			default:
				result := fmt.Sprintf("IsNumeric:Not a numeric type: %v", v)
				OperatorLog("OpIsNumeric", "input:", vars, "result:", result)
				return result
			}
		}
		OperatorLog("OpIsNumeric", "input:", vars, "result: true")
		return true
	}
	result := "IsNumeric:Not an array type"
	OperatorLog("OpIsNumeric", "input:", vars, "result:", result)
	return result
}

func OpIsInt(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			switch v.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				continue
			default:
				result := fmt.Sprintf("IsInt:Not an integer type: %v", v)
				OperatorLog("OpIsInt", "input:", vars, "result:", result)
				return result
			}
		}
		OperatorLog("OpIsInt", "input:", vars, "result: true")
		return true
	}
	result := "IsInt:Not an array type"
	OperatorLog("OpIsInt", "input:", vars, "result:", result)
	return result
}

func OpAsString(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) > 0 {
		result := reflect.ValueOf(arr[0]).String()
		OperatorLog("OpAsString", "input:", vars, "result:", result)
		return result
	}
	OperatorLog("OpAsString", "input:", vars, "result: empty string")
	return ""
}

func OpIsString(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if _, ok := v.(string); !ok {
				i.result = fmt.Sprintf("IsString:Not a string type: %v", v)
				return false
			}
		}
		return true
	}

	D.DebugPanic("OpIsString", "input:", vars, "result: false")
	return false
}

func OpIsNull(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if v != nil {
				result := fmt.Sprintf("IsNull:Not a null type: %v", v)
				OperatorLog("OpIsNull", "input:", vars, "result:", result)
				return result
			}
		}
		OperatorLog("OpIsNull", "input:", vars, "result: true")
		return true
	}
	result := "IsNull:Not an array type"
	OperatorLog("OpIsNull", "input:", vars, "result:", result)
	return result
}

func OpIsBool(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if _, ok := v.(bool); !ok {
				result := fmt.Sprintf("IsBool:Not a boolean type: %v", v)
				OperatorLog("OpIsBool", "input:", vars, "result:", result)
				return result
			}
		}
		OperatorLog("OpIsBool", "input:", vars, "result: true")
		return true
	}
	result := "IsBool:Not an array type"
	OperatorLog("OpIsBool", "input:", vars, "result:", result)
	return result
}

func OpIsArray(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if reflect.TypeOf(v).Kind() != reflect.Slice && reflect.TypeOf(v).Kind() != reflect.Array {
				result := fmt.Sprintf("IsArray:Not an array type: %v", v)
				OperatorLog("OpIsArray", "input:", vars, "result:", result)
				return result
			}
		}
		OperatorLog("OpIsArray", "input:", vars, "result: true")
		return true
	}
	result := "IsArray:Not an array type"
	OperatorLog("OpIsArray", "input:", vars, "result:", result)
	return result
}

func OpIsDouble(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			switch v.(type) {
			case float32, float64:
				continue
			default:
				result := fmt.Sprintf("IsDouble:Not a double type: %v", v)
				OperatorLog("OpIsDouble", "input:", vars, "result:", result)
				return result
			}
		}
		OperatorLog("OpIsDouble", "input:", vars, "result: true")
		return true
	}
	result := "IsDouble:Not an array type"
	OperatorLog("OpIsDouble", "input:", vars, "result:", result)
	return result
}
