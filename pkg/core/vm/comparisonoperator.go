package vm

func OpEq(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return false
	}

	left := vars[0]
	right := vars[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) == toFloat64(right)
	default:
		return left == right
	}
}

func OpNeq(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return false
	}

	left := vars[0]
	right := vars[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) != toFloat64(right)
	default:
		return left != right
	}
}

func OpGt(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return false
	}

	left := vars[0]
	right := vars[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) > toFloat64(right)
	default:
		return false
	}
}

func OpGte(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return false
	}

	left := vars[0]
	right := vars[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) >= toFloat64(right)
	default:
		return false
	}
}

func OpLt(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return false
	}

	left := vars[0]
	right := vars[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) < toFloat64(right)
	default:
		return false
	}
}

func OpLte(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return false
	}

	left := vars[0]
	right := vars[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) <= toFloat64(right)
	default:
		return false
	}
}

func isNumeric(v interface{}) interface{} {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	default:
		return 0
	}
}
