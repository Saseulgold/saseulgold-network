package vm

import (
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
				return false
			}
		}
		return true
	}
	return false
}

func OpIsInt(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			switch v.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				continue
			default:
				return false
			}
		}
		return true
	}
	return false
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
				return false
			}
		}
		return true
	}
	return false
}

func OpIsNull(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if v != nil {
				return false
			}
		}
		return true
	}
	return false
}

func OpIsBool(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if _, ok := v.(bool); !ok {
				return false
			}
		}
		return true
	}
	return false
}

func OpIsArray(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if reflect.TypeOf(v).Kind() != reflect.Slice && reflect.TypeOf(v).Kind() != reflect.Array {
				return false
			}
		}
		return true
	}
	return false
}

func OpIsDouble(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			switch v.(type) {
			case float32, float64:
				continue
			default:
				return false
			}
		}
		return true
	}
	return false
}
