package vm

import (
	"encoding/json"
	"hello/pkg/util"
	"regexp"
	"strings"
)

func OpArrayPush(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 3 {
		return []interface{}{}
	}

	origin, ok := vars[0].([]interface{})
	if !ok {
		return []interface{}{}
	}

	key := vars[1]
	value := vars[2]
	origin = append(origin, map[string]interface{}{key.(string): value})

	return origin
}

func OpConcat(i *Interpreter, vars []interface{}) interface{} {
	var result strings.Builder

	for _, v := range vars {
		if str, ok := v.(string); ok {
			result.WriteString(str)
		}
	}

	return result.String()
}

func OpCount(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return 0
	}

	if arr, ok := vars[0].([]interface{}); ok {
		return len(arr)
	}

	return 0
}

func OpStrlen(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return 0
	}

	if str, ok := vars[0].(string); ok {
		return len(str)
	}

	return 0
}

func OpRegMatch(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 2 {
		return false
	}

	pattern, ok1 := vars[0].(string)
	value, ok2 := vars[1].(string)

	if !ok1 || !ok2 {
		return false
	}

	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}

	return matched
}

func OpEncodeJson(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return ""
	}

	bytes, err := json.Marshal(vars[0])
	if err != nil {
		return ""
	}

	return string(bytes)
}

func OpDecodeJson(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return []interface{}{}
	}

	str, ok := vars[0].(string)
	if !ok {
		return []interface{}{}
	}

	var result []interface{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return []interface{}{}
	}

	return result
}

func OpHashLimit(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return ""
	}

	difficulty, ok := vars[0].(string)
	if !ok {
		return ""
	}

	// TODO: Implement hash limit logic
	return difficulty
	// return crypto.HashLimit(difficulty)
}

func OpHashMany(i *Interpreter, vars []interface{}) interface{} {
	var result strings.Builder
	for _, v := range vars {
		if str, ok := v.(string); ok {
			result.WriteString(str)
		}
	}

	return util.Hash(result.String())
}

func OpHash(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return ""
	}

	if str, ok := vars[0].(string); ok {
		return util.Hash(str)
	}

	return ""
}

func OpShortHash(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return ""
	}

	if str, ok := vars[0].(string); ok {
		return util.ShortHash(str)
	}

	return ""
}

func OpIdHash(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) == 0 {
		return ""
	}

	if str, ok := vars[0].(string); ok {
		return util.IDHash(str)
	}

	return ""
}

func OpSignVerify(i *Interpreter, vars []interface{}) interface{} {
	if len(vars) < 3 {
		return false
	}

	// obj := vars[0]
	// publicKey, ok1 := vars[1].(string)
	// signature, ok2 := vars[2].(string)

	// if !ok1 || !ok2 {
	// 	return false
	// }

	// return crypto.VerifySignature(obj, publicKey, signature)
	return true
}
