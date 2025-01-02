package vm

import (
	"encoding/json"
	"fmt"
	"hello/pkg/util"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	C "hello/pkg/core/config"
)

func OpSet(i *Interpreter, vars interface{}) interface{} {
	origin, key, value := Unpack3(vars)
	var originObj map[string]interface{}
	var ok bool

	if origin == nil {
		originObj = make(map[string]interface{})
		ok = true
	} else {
		originObj, ok = origin.(map[string]interface{})
	}

	if !ok {
		OperatorLog("OpObjectSet", "input:", vars, "result:", map[string]interface{}{})
		return map[string]interface{}{}
	}

	keyStr, ok := key.(string)
	if !ok {
		OperatorLog("OpObjectSet", "input:", vars, "result:", originObj)
		return originObj
	}

	originObj[keyStr] = value
	OperatorLog("OpObjectSet", "input:", vars, "result:", originObj)

	return originObj
}

func OpArrayPush(i *Interpreter, vars interface{}) interface{} {
	origin, key, value := Unpack3(vars)

	originArr, ok := origin.([]interface{})
	if !ok {
		OperatorLog("OpArrayPush", "input:", vars, "result:", []interface{}{})
		return []interface{}{}
	}

	keyStr, ok := key.(string)
	if !ok {
		OperatorLog("OpArrayPush", "input:", vars, "result:", []interface{}{})
		return []interface{}{}
	}

	originArr = append(originArr, map[string]interface{}{keyStr: value})
	OperatorLog("OpArrayPush", "input:", vars, "result:", originArr)
	return originArr
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
	arr := Unpack1(vars)

	if arr2, ok := arr.([]interface{}); ok {
		result := len(arr2)
		OperatorLog("OpCount", "input:", vars, "result:", result)
		return result
	}

	OperatorLog("OpCount", "input:", vars, "result:", 0)
	return 0
}

func OpStrlen(i *Interpreter, vars interface{}) interface{} {
	str := Unpack1(vars)

	if strVal, ok := str.(string); ok {
		result := len(strVal)
		OperatorLog("OpStrlen", "input:", vars, "result:", result)
		return result
	}

	OperatorLog("OpStrlen", "input:", vars, "result:", 0)
	return 0
}

func OpRegMatch(i *Interpreter, vars interface{}) interface{} {
	pattern, value := Unpack2(vars)

	patternStr, ok1 := pattern.(string)
	valueStr, ok2 := value.(string)

	if !ok1 || !ok2 {
		OperatorLog("OpRegMatch", "input:", vars, "result:", false)
		return false
	}

	if len(patternStr) >= 2 && patternStr[0] == '/' && patternStr[len(patternStr)-1] == '/' {
		patternStr = patternStr[1 : len(patternStr)-1]
	}

	matched, err := regexp.MatchString(patternStr, valueStr)
	if err != nil {
		OperatorLog("OpRegMatch", "input:", vars, "result:", false)
		return false
	}

	OperatorLog("OpRegMatch", "input:", vars, "result:", matched)
	return matched
}

func OpEncodeJson(i *Interpreter, vars interface{}) interface{} {
	value := Unpack1(vars)

	bytes, err := json.Marshal(value)
	if err != nil {
		OperatorLog("OpEncodeJson", "input:", vars, "result:", "")
		return ""
	}

	result := string(bytes)
	OperatorLog("OpEncodeJson", "input:", vars, "result:", result)
	return result
}

func OpDecodeJson(i *Interpreter, item interface{}) interface{} {
	fmt.Println("OpDecodeJson", "input:", item)
	arr, ok := item.([]interface{})
	if !ok {
		return nil
	}

	str, ok := arr[0].(string)
	if !ok {
		return nil
	}

	var result interface{}
	err := json.Unmarshal([]byte(str), &result)

	if err != nil {
		return nil
	}

	return result
}

func OpHashLimit(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}

	difficulty, ok := arr[0].(string)
	if !ok {
		return nil
	}

	diff, err := strconv.ParseInt(difficulty, 16, 64)
	if err != nil {
		return nil
	}

	hashCount := new(big.Int)
	hashCount.SetString(C.HASH_COUNT, 16)

	result := hashCount.Div(hashCount, big.NewInt(int64(diff)))

	// Convert to hexadecimal and pad with zeros to 64 characters
	hexStr := fmt.Sprintf("%x", result)
	paddedHexStr := fmt.Sprintf("%078s", hexStr)
	return strings.ReplaceAll(paddedHexStr, " ", "0")
}

func OpHashMany(i *Interpreter, vars interface{}) interface{} {
	var result strings.Builder

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			if str, ok := v.(string); ok {
				result.WriteString(str)
				result.WriteString(",")
			}
		}
	} else {
		panic("OpHashMany: input is not an array")
	}

	hashResult := util.Hash(result.String())
	OperatorLog("OpHashMany", "input:", vars, "string: ", result.String(), "result:", hashResult)
	return hashResult
}

func OpHash(i *Interpreter, vars interface{}) interface{} {
	// TODO
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

func OpLen(i *Interpreter, vars interface{}) interface{} {
	arr := Unpack1(vars)

	switch v := arr.(type) {
	case string:
		return len(v)
	case []interface{}:
		return len(v)
	default:
		return 0
	}
}

func OpEra(i *Interpreter, vars interface{}) interface{} {
	mined := Unpack1(vars)
	// Total supply: 3.7 billion (3,700,000,000)

	era0 := "1850000000"
	era1 := "2775000000"
	era2 := "3237500000"
	era3 := "3700000000"

	minedInt, ok := mined.(string)
	if !ok {
		OperatorLog("OpEra", "input:", vars, "result:", 0)
		return 0
	}

	var era int
	switch {
	case minedInt < era0:
		era = 0
	case minedInt < era1:
		era = 1
	case minedInt < era2:
		era = 2
	case minedInt < era3:
		era = 3
	default:
		era = 4
	}

	OperatorLog("OpEra", "input:", vars, "result:", era)
	return era
}
