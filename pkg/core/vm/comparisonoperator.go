package vm

import (
	. "hello/pkg/core/abi"
	"strconv"
)

func OpEq(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 2 {
		return false
	}

	left := arr[0]
	right := arr[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) == toFloat64(right)
	default:
		return left == right
	}
}

func OpNeq(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 2 {
		return false
	}

	left := arr[0]
	right := arr[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) != toFloat64(right)
	default:
		return left != right
	}
}

func OpGt(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	DebugLog("OpGt: arr =", arr)

	if !ok || len(arr) < 2 {
		DebugLog("OpGt: vars is not an array or has less than 2 elements")
		return false
	}

	left := arr[0]
	right := arr[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		DebugLog("OpGt: left =", toFloat64(left))
		DebugLog("OpGt: right =", toFloat64(right))
		return toFloat64(left) > toFloat64(right)
	default:
		DebugLog("OpGt: left or right is not numeric")
		return false
	}
}

func OpGte(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 2 {
		return false
	}

	left := arr[0]
	right := arr[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) >= toFloat64(right)
	default:
		return false
	}
}

func OpLt(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 2 {
		return false
	}

	left := arr[0]
	right := arr[1]

	switch {
	case isNumeric(left).(bool) && isNumeric(right).(bool):
		return toFloat64(left) < toFloat64(right)
	default:
		return false
	}
}

func OpLte(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 2 {
		return false
	}

	left := arr[0]
	right := arr[1]

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
		// 문자열이 숫자 형식인지 확인
		if len(val) == 0 {
			return false
		}
		// 첫 문자가 - 인 경우 음수 허용
		start := 0
		if val[0] == '-' {
			if len(val) == 1 {
				return false
			}
			start = 1
		}
		// 소수점 카운트
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
