package vm

import (
	"encoding/json"
	"hello/pkg/util"
	"regexp"
	"strings"
)

func OpArrayPush(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 3 {
		OperatorLog("OpArrayPush", "input:", vars, "result:", []interface{}{})
		return []interface{}{}
	}

	origin, ok := arr[0].([]interface{})
	if !ok {
		OperatorLog("OpArrayPush", "input:", vars, "result:", []interface{}{})
		return []interface{}{}
	}

	key := arr[1]
	value := arr[2]
	origin = append(origin, map[string]interface{}{key.(string): value})

	OperatorLog("OpArrayPush", "input:", vars, "result:", origin)
	return origin
}

func OpConcat(i *Interpreter, vars interface{}) interface{} {
	var result strings.Builder

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if str, ok := v.(string); ok {
				result.WriteString(str)
			}
		}
	}

	finalResult := result.String()
	OperatorLog("OpConcat", "input:", vars, "result:", finalResult)
	return finalResult
}

func OpCount(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) == 0 {
		OperatorLog("OpCount", "input:", vars, "result:", 0)
		return 0
	}

	if arr2, ok := arr[0].([]interface{}); ok {
		result := len(arr2)
		OperatorLog("OpCount", "input:", vars, "result:", result)
		return result
	}

	OperatorLog("OpCount", "input:", vars, "result:", 0)
	return 0
}

func OpStrlen(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) == 0 {
		OperatorLog("OpStrlen", "input:", vars, "result:", 0)
		return 0
	}

	if str, ok := arr[0].(string); ok {
		result := len(str)
		OperatorLog("OpStrlen", "input:", vars, "result:", result)
		return result
	}

	OperatorLog("OpStrlen", "input:", vars, "result:", 0)
	return 0
}

func OpRegMatch(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 2 {
		OperatorLog("OpRegMatch", "input:", vars, "result:", false)
		return false
	}

	pattern, ok1 := arr[0].(string)
	value, ok2 := arr[1].(string)

	if !ok1 || !ok2 {
		OperatorLog("OpRegMatch", "input:", vars, "result:", false)
		return false
	}

	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		OperatorLog("OpRegMatch", "input:", vars, "result:", false)
		return false
	}

	OperatorLog("OpRegMatch", "input:", vars, "result:", matched)
	return matched
}

func OpEncodeJson(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) == 0 {
		OperatorLog("OpEncodeJson", "input:", vars, "result:", "")
		return ""
	}

	bytes, err := json.Marshal(arr[0])
	if err != nil {
		OperatorLog("OpEncodeJson", "input:", vars, "result:", "")
		return ""
	}

	result := string(bytes)
	OperatorLog("OpEncodeJson", "input:", vars, "result:", result)
	return result
}

func OpDecodeJson(i *Interpreter, vars interface{}) interface{} {
	if arr, ok := vars.([]interface{}); ok && len(arr) > 0 {
		// 첫 번째 요소가 이미 map인 경우 직접 반환
		if m, ok := arr[0].(map[string]interface{}); ok {
			OperatorLog("OpDecodeJson", "input:", vars, "result:", m)
			return m
		}

		// 문자열인 경우 JSON 디코딩 시도
		if jsonStr, ok := arr[0].(string); ok {
			var result interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err == nil {
				OperatorLog("OpDecodeJson", "input:", vars, "result:", result)
				return result
			}
		}

		// 그 외의 경우 원본 값 반환
		OperatorLog("OpDecodeJson", "input:", vars, "result:", arr[0])
		return arr[0]
	}

	OperatorLog("OpDecodeJson", "input:", vars, "result:", nil)
	return nil
}

func OpHashLimit(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) == 0 {
		OperatorLog("OpHashLimit", "input:", vars, "result:", "")
		return ""
	}

	difficulty, ok := arr[0].(string)
	if !ok {
		OperatorLog("OpHashLimit", "input:", vars, "result:", "")
		return ""
	}

	// TODO: Implement hash limit logic
	OperatorLog("OpHashLimit", "input:", vars, "result:", difficulty)
	return difficulty
	// return crypto.HashLimit(difficulty)
}

func OpHashMany(i *Interpreter, vars interface{}) interface{} {
	var result strings.Builder

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if str, ok := v.(string); ok {
				result.WriteString(str)
			}
		}
	}

	hashResult := util.Hash(result.String())
	OperatorLog("OpHashMany", "input:", vars, "result:", hashResult)
	return hashResult
}

func OpHash(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) == 0 {
		OperatorLog("OpHash", "input:", vars, "result:", "")
		return ""
	}

	if str, ok := arr[0].(string); ok {
		result := util.Hash(str)
		OperatorLog("OpHash", "input:", vars, "result:", result)
		return result
	}

	OperatorLog("OpHash", "input:", vars, "result:", "")
	return ""
}

func OpShortHash(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) == 0 {
		OperatorLog("OpShortHash", "input:", vars, "result:", "")
		return ""
	}

	if str, ok := arr[0].(string); ok {
		result := util.ShortHash(str)
		OperatorLog("OpShortHash", "input:", vars, "result:", result)
		return result
	}

	OperatorLog("OpShortHash", "input:", vars, "result:", "")
	return ""
}

func OpIdHash(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) == 0 {
		OperatorLog("OpIdHash", "input:", vars, "result:", "")
		return ""
	}

	if str, ok := arr[0].(string); ok {
		result := util.IDHash(str)
		OperatorLog("OpIdHash", "input:", vars, "result:", result)
		return result
	}

	OperatorLog("OpIdHash", "input:", vars, "result:", "")
	return ""
}

func OpSignVerify(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 3 {
		OperatorLog("OpSignVerify", "input:", vars, "result:", false)
		return false
	}

	// obj := arr[0]
	// publicKey, ok1 := arr[1].(string)
	// signature, ok2 := arr[2].(string)

	// if !ok1 || !ok2 {
	// 	return false
	// }

	// return crypto.VerifySignature(obj, publicKey, signature)
	OperatorLog("OpSignVerify", "input:", vars, "result:", true)
	return true
}
