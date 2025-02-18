package vm

import (
	"fmt"
	"reflect"
	"strconv"
)

func OperatorLog(args ...interface{}) interface{} {
	// fmt.Println(args...)
	return true
}

func OpGetType(i *Interpreter, vars interface{}) interface{} {
	value := Unpack1(vars)
	result := reflect.TypeOf(value).String()
	OperatorLog("OpGetType", "input:", vars, "result:", result)
	return result
}

func OpIsNumeric(i *Interpreter, vars interface{}) interface{} {
	arr := vars.([]interface{})
	for _, v := range arr {
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			continue
		case string:
			str := v.(string)
			if _, err := strconv.ParseFloat(str, 64); err == nil {
				continue
			}
			result := fmt.Sprintf("IsNumeric:Not a numeric string: %v", v)
			OperatorLog("OpIsNumeric", "input:", vars, "result:", result)
			return result
		default:
			result := fmt.Sprintf("IsNumeric:Not a numeric type: %v", v)
			OperatorLog("OpIsNumeric", "input:", vars, "result:", result)
			return result
		}
	}
	OperatorLog("OpIsNumeric", "input:", vars, "result: true")
	return true
}

func OpIsInt(i *Interpreter, vars interface{}) interface{} {
	arr := vars.([]interface{})
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

func OpAsString(i *Interpreter, vars interface{}) interface{} {
	value := Unpack1(vars)

	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		result := strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
		OperatorLog("OpAsString", "input:", vars, "result:", result)
		return result
	case string:
		OperatorLog("OpAsString", "input:", vars, "result:", v)
		return v
	default:
		result := fmt.Sprintf("%v", value)
		OperatorLog("OpAsString", "input:", vars, "result:", result)
		return result
	}
}

func OpIsString(i *Interpreter, vars interface{}) interface{} {
	arr := vars.([]interface{})
	for _, v := range arr {
		if _, ok := v.(string); !ok {
			result := fmt.Sprintf("IsString:Not a string type: %v", v)
			OperatorLog("OpIsString", "input:", vars, "result:", result)
			return result
		}
	}
	OperatorLog("OpIsString", "input:", vars, "result: true")
	return true
}

func OpIsNull(i *Interpreter, vars interface{}) interface{} {
	arr := vars.([]interface{})
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

func OpIsBool(i *Interpreter, vars interface{}) interface{} {
	arr := vars.([]interface{})
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

func OpIsArray(i *Interpreter, vars interface{}) interface{} {
	arr := vars.([]interface{})
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

func OpIsDouble(i *Interpreter, vars interface{}) interface{} {
	arr := vars.([]interface{})
	for _, v := range arr {
		switch v.(type) {
		case float32, float64:
			continue
		default:
			result := fmt.Sprintf("IsDouble:Not a double type: %v, %v", v, reflect.TypeOf(v))
			OperatorLog("OpIsDouble", "input:", vars, "result:", result)
			return result
		}
	}
	OperatorLog("OpIsDouble", "input:", vars, "result: true")
	return true
}
