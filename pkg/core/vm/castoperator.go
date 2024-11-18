package vm

import (
	"fmt"
	"reflect"
)

func OpGetType(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) > 0 {
		return reflect.TypeOf(arr[0]).String()
	}
	return "nil"
}

func OpIsNumeric(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			switch v.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				continue
			default:
				return fmt.Sprintf("IsNumeric:Not a numeric type: %v", v)
			}
		}
		return true
	}
	return "IsNumeric:Not an array type"
}

func OpIsInt(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			switch v.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				continue
			default:
				return fmt.Sprintf("IsInt:Not an integer type: %v", v)
			}
		}
		return true
	}
	return "IsInt:Not an array type"
}

func OpAsString(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) > 0 {
		return reflect.ValueOf(arr[0]).String()
	}
	return ""
}

func OpIsString(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if _, ok := v.(string); !ok {
				return fmt.Sprintf("IsString:Not a string type: %v", v)
			}
		}
		return true
	}
	return "IsString:Not an array type"
}

func OpIsNull(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if v != nil {
				return fmt.Sprintf("IsNull:Not a null type: %v", v)
			}
		}
		return true
	}
	return "IsNull:Not an array type"
}

func OpIsBool(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if _, ok := v.(bool); !ok {
				return fmt.Sprintf("IsBool:Not a boolean type: %v", v)
			}
		}
		return true
	}
	return "IsBool:Not an array type"
}

func OpIsArray(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if reflect.TypeOf(v).Kind() != reflect.Slice && reflect.TypeOf(v).Kind() != reflect.Array {
				return fmt.Sprintf("IsArray:Not an array type: %v", v)
			}
		}
		return true
	}
	return "IsArray:Not an array type"
}

func OpIsDouble(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			switch v.(type) {
			case float32, float64:
				continue
			default:
				return fmt.Sprintf("IsDouble:Not a double type: %v", v)
			}
		}
		return true
	}
	return "IsDouble:Not an array type"
}
