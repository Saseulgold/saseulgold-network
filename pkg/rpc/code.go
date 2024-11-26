package rpc

import (
	. "hello/pkg/core/model"
	"hello/pkg/core/vm/native"
	F "hello/pkg/util"
)

type Code struct{}

var SystemMethods = []string{"Genesis", "Register", "Grant", "Revoke", "Oracle", "Faucet", "Publish", "Send", "Submit"}

func (c *Code) Contracts() map[string]map[string]*Method {
	contracts := make(map[string]map[string]*Method)
	rootCid := F.RootSpaceId()

	contracts[rootCid]["Genesis"] = native.Genesis()
	contracts[rootCid]["Register"] = native.Register()
	contracts[rootCid]["Revoke"] = native.Revoke()
	contracts[rootCid]["Faucet"] = native.Faucet()
	contracts[rootCid]["Publish"] = native.Publish()
	contracts[rootCid]["Send"] = native.Send()

	return contracts
}

func (c *Code) GetNativeContract(methodName string) *Method {
	return c.Contracts()[F.RootSpace()][methodName]
}

func (c *Code) GetContract(cid string, methodName string) *Method {
	return c.Contracts()[cid][methodName]
}
