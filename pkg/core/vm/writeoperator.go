package vm

import (
	. "hello/pkg/core/abi"
	C "hello/pkg/core/config"
	"hello/pkg/core/debug"
	F "hello/pkg/util"
)

func OpWriteUniversal(i *Interpreter, vars interface{}) interface{} {
	OperatorLog("op start", i.process.String(), i.state.String())

	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 3 {
		debug.DebugPanic("OpWriteUniversal", "input:", vars, "result:", nil)
		return nil
	}

	attr, ok1 := arr[0].(string)
	key, ok2 := arr[1].(string)
	if !ok1 || !ok2 {
		debug.DebugPanic("OpWriteUniversal", "input:", vars, "result:", nil)
		return nil
	}

	var statusHash string
	if i.process == ProcessMain {
		statusHash = F.StatusHash(i.code.GetWriter(), i.code.GetSpace(), attr, key)
		debug.DebugLog("write_universal main: attr=%s; key=%s; writer=%s; space=%s; status_hash=%s",
			attr, key, i.code.GetWriter(), i.code.GetSpace(), statusHash)
	} else if i.process == ProcessPost {
		statusHash = F.StatusHash(i.postProcess.GetWriter(), i.postProcess.GetSpace(), attr, key)
	}

	if statusHash == "" {
		return nil
	}

	var result interface{}
	switch i.state {
	case StateRead:
		OperatorLog("OpWriteUniversal READ ", statusHash)
		i.AddUniversalLoads(statusHash)
		result = nil
	case StateCondition:
		OperatorLog("OpWriteUniversal CONDITION ", statusHash, arr[2])
		value := arr[2]
		length := len(F.String(value))

		if length > C.STATUS_SIZE_LIMIT {
			i.breakFlag = true
			i.result = "Too long status values. maximum size: " + string(C.STATUS_SIZE_LIMIT)
		}

		if i.process == ProcessMain {
			i.SignedData.SetCachedUniversal(statusHash, value)
			i.weight += int64(len(statusHash) + length)
		}

		return ABI{
			Key:   "$write_universal",
			Value: arr,
		}
	case StateExecution:
		value := arr[2]
		result = i.SetUniversalStatus(statusHash, value)
	}

	return result
}

func OpWriteLocal(i *Interpreter, vars interface{}) interface{} {
	arr, ok := vars.([]interface{})
	if !ok || len(arr) < 3 {
		OperatorLog("OpWriteLocal", "input:", vars, "result:", nil)
		return nil
	}

	attr, ok1 := arr[0].(string)
	key, ok2 := arr[1].(string)
	if !ok1 || !ok2 {
		OperatorLog("OpWriteLocal", "input:", vars, "result:", nil)
		return nil
	}

	var statusHash string
	if i.process == ProcessMain {
		statusHash = F.StatusHash(i.code.GetWriter(), i.code.GetSpace(), attr, key)
		debug.DebugLog("write_local main: key=%s; writer=%s; space=%s; value=%s",
			key, i.code.GetWriter(), i.code.GetSpace(), F.String(arr[2])[:30])
	} else if i.process == ProcessPost {
		statusHash = F.StatusHash(i.postProcess.GetWriter(), i.postProcess.GetSpace(), attr, key)
	}

	if statusHash == "" {
		OperatorLog("OpWriteLocal", "input:", vars, "result:", nil)
		return nil
	}

	var result interface{}
	switch i.state {
	case StateRead:
		i.AddLocalLoads(statusHash)
		result = nil
	case StateCondition:
		value := arr[2]
		length := len(F.String(value))

		if length > C.STATUS_SIZE_LIMIT {
			i.breakFlag = true
			i.result = "Too long status values. maximum size: " + string(C.STATUS_SIZE_LIMIT)
		}

		if i.process == ProcessMain {
			i.SignedData.SetCachedLocal(statusHash, value)
			i.weight += int64(len(statusHash) + length*1000000000)
		}
		result = map[string]interface{}{
			"$write_local": arr,
		}
	case StateExecution:
		value := arr[2]
		result = i.SetLocalStatus(statusHash, value)
	}

	OperatorLog("OpWriteLocal", "input:", vars, "result:", result)
	return result
}
