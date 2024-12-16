package rpc

import (
	. "hello/pkg/core/model"
	"hello/pkg/core/vm/native"
	F "hello/pkg/util"
)

type Code struct{}

func NativeRequests() map[string]map[string]*Method {
	requests := make(map[string]map[string]*Method)
	rootCid := F.RootSpaceId()

	requests[rootCid] = make(map[string]*Method)
	requests[rootCid]["GetBlock"] = native.GetBlock()
	requests[rootCid]["ListBlock"] = native.ListBlock()
	requests[rootCid]["GetBalance"] = native.GetBalance()
	return requests
}

func NativeContracts() map[string]map[string]*Method {
	contracts := make(map[string]map[string]*Method)
	rootCid := F.RootSpaceId()

	contracts[rootCid] = make(map[string]*Method)
	contracts[rootCid]["Genesis"] = native.Genesis()
	contracts[rootCid]["Register"] = native.Register()
	contracts[rootCid]["Revoke"] = native.Revoke()
	contracts[rootCid]["Faucet"] = native.Faucet()
	contracts[rootCid]["Publish"] = native.Publish()
	contracts[rootCid]["Send"] = native.Send()

	return contracts
}

func GetNativeContract(methodName string) *Method {
	return NativeContracts()[F.RootSpace()][methodName]
}

/**
func GetContract(cid string, methodName string) *Method {
	return Contracts()[cid][methodName]
}
**/
