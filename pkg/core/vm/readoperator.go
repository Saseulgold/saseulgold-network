package vm

import (
	. "hello/pkg/core/abi"
	F "hello/pkg/util"
)

func OpLoadParam(i *Interpreter, vars interface{}) interface{} {
	var result interface{}

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			str, ok := v.(string)
			if !ok {
				DebugLog("OpLoadParam: value is not string")
				return nil
			}

			if result == nil {
				result = i.SignedData.GetAttribute(str)
			} else if arr, ok := result.(map[string]interface{}); ok {
				if val, exists := arr[str]; exists {
					result = val
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	return result
}

func OpReadLocal(i *Interpreter, vars interface{}) interface{} {
	var attr, key string
	var defaultVal interface{}

	if arr, ok := vars.([]interface{}); ok {
		if len(arr) > 0 {
			if v, ok := arr[0].(string); ok {
				attr = v
			}
		}
		if len(arr) > 1 {
			if v, ok := arr[1].(string); ok {
				key = v
			}
		}
		if len(arr) > 2 {
			defaultVal = arr[2]
		}
	}

	var statusHash string
	switch i.process {
	case ProcessMain:
		statusHash = F.StatusHash(i.code.GetWriter(), i.code.GetSpace(), attr, key)
	case ProcessPost:
		statusHash = F.StatusHash(i.postProcess.GetWriter(), i.postProcess.GetSpace(), attr, key)
	default:
		return nil
	}

	if statusHash == "" {
		return nil
	}

	switch i.state {
	case StateRead:
		i.AddLocalLoads(statusHash)
	case StateCondition:
		if i.process == ProcessPost {
			cachedData := i.SignedData.GetCachedLocal(statusHash)
			if cachedData != nil {
				return cachedData
			}
		}
		return i.GetLocalStatus(statusHash, defaultVal)
	case StateExecution:
		return i.GetLocalStatus(statusHash, defaultVal)
	}

	return nil
}

func OpReadUniversal(i *Interpreter, vars interface{}) interface{} {
	var attr, key string
	var defaultVal interface{}

	if arr, ok := vars.([]interface{}); ok {
		if len(arr) > 0 {
			if v, ok := arr[0].(string); ok {
				attr = v
			}
		}
		if len(arr) > 1 {
			if v, ok := arr[1].(string); ok {
				key = v
			}
		}
		if len(arr) > 2 {
			defaultVal = arr[2]
		}
	}

	var statusHash string
	switch i.process {
	case ProcessMain:
		statusHash = F.StatusHash(i.code.GetWriter(), i.code.GetSpace(), attr, key)
	case ProcessPost:
		statusHash = F.StatusHash(i.postProcess.GetWriter(), i.postProcess.GetSpace(), attr, key)
	default:
		return nil
	}

	if statusHash == "" {
		return nil
	}

	switch i.state {
	case StateRead:
		i.AddUniversalLoads(statusHash)
	case StateCondition:
		if i.process == ProcessPost {
			cachedData := i.SignedData.GetCachedUniversal(statusHash)
			if cachedData != nil {
				return cachedData
			}
		}
		return i.GetUniversalStatus(statusHash, defaultVal)
	case StateExecution:
		return i.GetUniversalStatus(statusHash, defaultVal)
	}

	return nil
}
