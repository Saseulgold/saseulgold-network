package vm

import (
	"reflect"
)

func OpGetType(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) > 0 {
		return reflect.TypeOf(vars[0]).String()
	}
	return "nil"
}

func OpIsNumeric(i *Interpreter, vars []interface{}) interface{} {
	for _, v := range vars {
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			continue
		default:
			return false
		}
	}
	return true
}

func OpIsInt(i *Interpreter, vars []interface{}) interface{} {
	for _, v := range vars {
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			continue
		default:
			return false
		}
	}
	return true
}

func OpAsString(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) > 0 {
		return reflect.ValueOf(vars[0]).String()
	}
	return ""
}

func OpIsString(i *Interpreter, vars []interface{}) interface{} {
	for _, v := range vars {
		if _, ok := v.(string); !ok {
			return false
		}
	}
	return true
}

func OpIsNull(i *Interpreter, vars []interface{}) interface{} {
	for _, v := range vars {
		if v != nil {
			return false
		}
	}
	return true
}

func OpIsBool(i *Interpreter, vars []interface{}) interface{} {
	for _, v := range vars {
		if _, ok := v.(bool); !ok {
			return false
		}
	}
	return true
}

func OpIsArray(i *Interpreter, vars []interface{}) interface{} {
	for _, v := range vars {
		if reflect.TypeOf(v).Kind() != reflect.Slice && reflect.TypeOf(v).Kind() != reflect.Array {
			return false
		}
	}
	return true
}

func OpIsDouble(i *Interpreter, vars []interface{}) interface{} {
	for _, v := range vars {
		switch v.(type) {
		case float32, float64:
			continue
		default:
			return false
		}
	}
	return true
}
