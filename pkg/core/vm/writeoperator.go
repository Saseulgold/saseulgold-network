package vm

import (
	. "hello/pkg/core/abi"
	C "hello/pkg/core/config"
	"hello/pkg/core/debug"
	F "hello/pkg/util"
)

func OpWriteUniversal(i *Interpreter, vars interface{}) interface{} {
	OperatorLog("op start", i.process.String(), i.state.String())

	attr, key, value := Unpack3(vars)

	_, ok := attr.(string)
	if !ok {
		OperatorLog("OpWriteUniversal", "input:", vars, "result: nil")
		return nil
	}
	_, ok = key.(string)
	if !ok {
		OperatorLog("OpWriteUniversal", "input:", vars, "result: nil")
		return nil
	}

	var statusHash string
	if i.process == ProcessMain {
		statusHash = F.StatusHash(i.code.GetWriter(), i.code.GetSpace(), attr.(string), key.(string))
		debug.DebugLog("write_universal main: attr=%s; key=%s; writer=%s; space=%s; status_hash=%s",
			attr, key, i.code.GetWriter(), i.code.GetSpace(), statusHash)
	} else if i.process == ProcessPost {
		statusHash = F.StatusHash(i.postProcess.GetWriter(), i.postProcess.GetSpace(), attr.(string), key.(string))
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
		OperatorLog("OpWriteUniversal CONDITION ", statusHash, value)
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
			Value: vars,
		}
	case StateExecution:
		result = i.SetUniversalStatus(statusHash, value)
	}

	return result
}

func OpWriteLocal(i *Interpreter, vars interface{}) interface{} {
	attr, key, value := Unpack3(vars)

	_, ok := attr.(string)
	if !ok {
		OperatorLog("OpWriteLocal", "input:", vars, "result: nil")
		return nil
	}
	_, ok = key.(string)
	if !ok {
		OperatorLog("OpWriteLocal", "input:", vars, "result: nil")
		return nil
	}

	var statusHash string
	if i.process == ProcessMain {
		statusHash = F.StatusHash(i.code.GetWriter(), i.code.GetSpace(), attr.(string), key.(string))
		debug.DebugLog("write_local main: key=%s; writer=%s; space=%s; value=%s",
			key, i.code.GetWriter(), i.code.GetSpace(), F.String(value)[:30])
	} else if i.process == ProcessPost {
		statusHash = F.StatusHash(i.postProcess.GetWriter(), i.postProcess.GetSpace(), attr.(string), key.(string))
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
		length := len(F.String(value))

		if length > C.STATUS_SIZE_LIMIT {
			i.breakFlag = true
			i.result = "Too long status values. maximum size: " + string(C.STATUS_SIZE_LIMIT)
		}

		if i.process == ProcessMain {
			i.SignedData.SetCachedLocal(statusHash, value)
			i.weight += int64(len(statusHash) + length*1000000000)
		}
		result = ABI{
			Key:   "$write_local",
			Value: vars,
		}
	case StateExecution:
		result = i.SetLocalStatus(statusHash, value)
	}

	OperatorLog("OpWriteLocal", "input:", vars, "result:", result)
	return result
}
