package vm

import (
	"fmt"
	F "hello/pkg/util"
)

func OpLoadParam(i *Interpreter, vars interface{}) interface{} {
	key := Unpack1(vars)

	return i.SignedData.GetAttribute(key.(string))
}

func OpReadLocal(i *Interpreter, vars interface{}) interface{} {
	attr, key, defaultVal := Unpack2Or3(vars)

	_, ok := attr.(string)
	if !ok {
		OperatorLog("OpReadLocal", "input:", vars, "result: nil")
		return nil
	}
	_, ok = key.(string)
	if !ok {
		OperatorLog("OpReadLocal", "input:", vars, "result: nil")
		return nil
	}

	var statusHash string
	switch i.process {
	case ProcessMain:
		statusHash = F.StatusHash(i.code.GetWriter(), i.code.GetSpace(), attr.(string), key.(string))
	case ProcessPost:
		statusHash = F.StatusHash(i.postProcess.GetWriter(), i.postProcess.GetSpace(), attr.(string), key.(string))
	default:
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

	return result
}

func OpReadUniversal(i *Interpreter, vars interface{}) interface{} {
	attr, key, defaultVal := Unpack3(vars)

	_, ok := attr.(string)
	if !ok {
		OperatorLog("OpReadUniversal", "input:", vars, "result: nil")
		return nil
	}
	_, ok = key.(string)
	if !ok {
		OperatorLog("OpReadUniversal", "input:", vars, "result: nil")
		return nil
	}

	var statusHash string
	switch i.process {
	case ProcessMain:
		statusHash = F.StatusHash(i.code.GetWriter(), i.code.GetSpace(), attr.(string), key.(string))
	case ProcessPost:
		statusHash = F.StatusHash(i.postProcess.GetWriter(), i.postProcess.GetSpace(), attr.(string), key.(string))
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

	fmt.Println(fmt.Sprintf("OpReadUniversal: attr=%s; key=%s; result=%v", attr, key, result))
	return result
}
