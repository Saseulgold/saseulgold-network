package vm

import (
	F "hello/pkg/util"
)

func OpLoadParam(i *Interpreter, vars interface{}) interface{} {
	var result interface{}

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			str, ok := v.(string)
			if !ok {
				OperatorLog("OpLoadParam", "input:", vars, "result: nil")
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

	OperatorLog("OpLoadParam", "input:", vars, "result:", result)
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
		OperatorLog("OpReadLocal", "input:", vars, "result: nil")
		return nil
	}

	if statusHash == "" {
		OperatorLog("OpReadLocal", "input:", vars, "result: nil")
		return nil
	}

	var result interface{}
	switch i.state {
	case StateRead:
		i.AddLocalLoads(statusHash)
	case StateCondition:
		if i.process == ProcessPost {
			cachedData := i.SignedData.GetCachedLocal(statusHash)
			if cachedData != nil {
				result = cachedData
				break
			}
		}
		result = i.GetLocalStatus(statusHash, defaultVal)
	case StateExecution:
		result = i.GetLocalStatus(statusHash, defaultVal)
	}

	OperatorLog("OpReadLocal", "input:", vars, "result:", result)
	return result
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
		OperatorLog("OpReadUniversal", "input:", vars, "result: nil")
		return nil
	}

	if statusHash == "" {
		OperatorLog("OpReadUniversal", "input:", vars, "result: nil")
		return nil
	}

	var result interface{}
	switch i.state {
	case StateRead:
		i.AddUniversalLoads(statusHash)
	case StateCondition:
		if i.process == ProcessPost {
			cachedData := i.SignedData.GetCachedUniversal(statusHash)
			if cachedData != nil {
				result = cachedData
				break
			}
		}
		result = i.GetUniversalStatus(statusHash, defaultVal)
	case StateExecution:
		result = i.GetUniversalStatus(statusHash, defaultVal)
	}

	OperatorLog("OpReadUniversal", "input:", vars, "result:", result)
	return result
}
