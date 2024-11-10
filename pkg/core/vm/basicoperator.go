package vm

import (
	"reflect"
)

type BasicOperator struct {
	state  string
	break_ bool
	result interface{}
	weight int64
}

func (b *BasicOperator) Condition(vars []interface{}) bool {
	if b.state != "CONDITION" {
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
		b.break_ = true
		if errMsg != "" {
			b.result = errMsg
		}
		return false
	}

	return true
}

func (b *BasicOperator) Response(vars []interface{}) interface{} {
	if b.state != "EXECUTION" {
		return map[string][]interface{}{
			"$response": vars,
		}
	}

	b.break_ = true
	if len(vars) > 0 {
		b.result = vars[0]
	} else {
		b.result = nil
	}
	return nil
}

func (b *BasicOperator) Weight(vars []interface{}) int64 {
	return b.weight
}

func (b *BasicOperator) If(vars []interface{}) interface{} {
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

func (b *BasicOperator) And(vars []interface{}) bool {
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

func (b *BasicOperator) Or(vars []interface{}) bool {
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

func (b *BasicOperator) Get(vars []interface{}) interface{} {
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

func (b *BasicOperator) In(vars []interface{}) bool {
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
