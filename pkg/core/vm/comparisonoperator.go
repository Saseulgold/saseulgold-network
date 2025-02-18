package vm

import (
	"strconv"
	 "math/big"

)

func OpEq(i *Interpreter, vars interface{}) interface{} {
	left, right := Unpack2(vars)

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) == toFloat64(right)
	default:
		return left == right
	}
}

func OpNeq(i *Interpreter, vars interface{}) interface{} {
	left, right := Unpack2(vars)

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) != toFloat64(right)
	default:
		return left != right
	}
}

func OpGt(i *Interpreter, vars interface{}) interface{} {
	left, right := Unpack2(vars)

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) > toFloat64(right)
	default:
		return false
	}
}

func OpGte(i *Interpreter, vars interface{}) interface{} {
	left, right := Unpack2(vars)

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) >= toFloat64(right)
	default:
		return false
	}
}

func isValidHex(s string) bool {
    if len(s) == 0 {
        return false
    }
    for _, c := range s {
        if !((c >= '0' && c <= '9') ||
            (c >= 'a' && c <= 'f') ||
            (c >= 'A' && c <= 'F')) {
            return false
        }
    }
    return true
}


func OpHashLimitOk(i *Interpreter, vars interface{}) interface{} {
	left, right := Unpack2(vars)

    leftStr, ok1 := left.(string)
    rightStr, ok2 := right.(string)
    if !ok1 || !ok2 {
        return false
    }

    // Verify that the hash string is a valid hexadecimal number
    if !isValidHex(leftStr) || !isValidHex(rightStr) {
        return false
    }

    // Converting to a Large integer
    leftInt := new(big.Int)
    _, ok := leftInt.SetString(leftStr, 16)
    if !ok {
        return false
    }

    rightInt := new(big.Int)
    _, ok = rightInt.SetString(rightStr, 16)
    if !ok {
        return false
    }

    // compare
    return leftInt.Cmp(rightInt) < 0
}

func OpLt(i *Interpreter, vars interface{}) interface{} {
	left, right := Unpack2(vars)

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) < toFloat64(right)
	default:
		return false
	}
}

func OpLte(i *Interpreter, vars interface{}) interface{} {
	left, right := Unpack2(vars)

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) <= toFloat64(right)
	default:
		return false
	}
}

func isNumeric(v interface{}) interface{} {
	switch val := v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		// Verify that the string is in numeric format
		if len(val) == 0 {
			return false
		}
		// Allow negative numbers if the first letter is -
		start := 0
		if val[0] == '-' {
			if len(val) == 1 {
				return false
			}
			start = 1
		}
		// decimal count
		dotCount := 0
		for i := start; i < len(val); i++ {
			if val[i] == '.' {
				dotCount++
				if dotCount > 1 {
					return false
				}
				continue
			}
			if val[i] < '0' || val[i] > '9' {
				return false
			}
		}
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
	case string:
		if val == "" {
			return 0
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
		return 0

	default:
		return 0
	}
}
