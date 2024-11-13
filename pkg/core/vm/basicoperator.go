package vm

import (
	"reflect"
)

func OpCondition(i *Interpreter, vars []interface{}) interface{} {
	if i.state != StateCondition {
		return true
	}

	var abi bool
	var errMsg string

	if len(vars) > 0 {
		if v, ok := vars[0].(bool); ok {
			abi = v
		}
	}

	if len(vars) > 1 {
		if v, ok := vars[1].(string); ok {
			errMsg = v
		}
	}

	if !abi {
		i.breakFlag = true
		if errMsg != "" {
			i.result = errMsg
		}
		return false
	}

	return true
}

func OpResponse(i *Interpreter, vars []interface{}) interface{} {
	if i.state != StateExecution {
		return map[string][]interface{}{
			"$response": vars,
		}
	}

	i.breakFlag = true
	if len(vars) > 0 {
		i.result = vars[0]
	} else {
		i.result = nil
	}
	return nil
}

func OpWeight(i *Interpreter, vars []interface{}) interface{} {
	return i.weight
}

func OpIf(i *Interpreter, vars []interface{}) interface{} {
	var condition bool
	var trueVal, falseVal interface{}

	if len(vars) > 0 {
		if v, ok := vars[0].(bool); ok {
			condition = v
		}
	}

	if len(vars) > 1 {
		trueVal = vars[1]
	}

	if len(vars) > 2 {
		falseVal = vars[2]
	}

	if condition {
		return trueVal
	}
	return falseVal
}

func OpAnd(i *Interpreter, vars []interface{}) interface{} {
	var result *bool

	for _, v := range vars {
		if boolVal, ok := v.(bool); !ok {
			return false
		} else {
			if result == nil {
				result = &boolVal
			} else {
				newVal := *result && boolVal
				result = &newVal
			}
		}
	}

	if result == nil {
		return false
	}
	return *result
}

func OpOr(i *Interpreter, vars []interface{}) interface{} {
	var result *bool

	for _, v := range vars {
		if boolVal, ok := v.(bool); !ok {
			return false
		} else {
			if result == nil {
				result = &boolVal
			} else {
				newVal := *result || boolVal
				result = &newVal
			}
		}
	}

	if result == nil {
		return false
	}
	return *result
}

func OpGet(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return nil
	}

	obj, ok := vars[0].(map[string]interface{})
	if !ok {
		return nil
	}

	var key string
	switch v := vars[1].(type) {
	case string:
		key = v
	case float64:
		key = reflect.TypeOf(v).String()
	default:
		return nil
	}

	var defaultVal interface{}
	if len(vars) > 2 {
		defaultVal = vars[2]
	}

	if val, exists := obj[key]; exists {
		return val
	}
	return defaultVal
}

func OpIn(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return false
	}

	target := vars[0]
	cases, ok := vars[1].([]interface{})
	if !ok {
		return false
	}

	for _, v := range cases {
		if reflect.DeepEqual(target, v) {
			return true
		}
	}
	return false
}
