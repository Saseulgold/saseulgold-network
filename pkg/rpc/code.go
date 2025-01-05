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
	requests[rootCid]["ListTransaction"] = native.ListTransaction()

	requests[rootCid]["GetBalance"] = native.GetBalance()
	requests[rootCid]["GetTokenInfo"] = native.GetTokenInfo()
	requests[rootCid]["GetPairInfo"] = native.GetPairInfo()

	requests[rootCid]["BalanceOf"] = native.BalanceOf()
	requests[rootCid]["BalanceOfLP"] = native.BalanceOfLP()

	return requests
}

func NativeContracts() map[string]map[string]*Method {
	contracts := make(map[string]map[string]*Method)
	rootCid := F.RootSpaceId()

	contracts[rootCid] = make(map[string]*Method)
	contracts[rootCid]["Genesis"] = native.Genesis()
	contracts[rootCid]["Faucet"] = native.Faucet()
	contracts[rootCid]["Publish"] = native.Publish()
	contracts[rootCid]["Send"] = native.Send()
	contracts[rootCid]["Mint"] = native.Mint()
	contracts[rootCid]["Transfer"] = native.Transfer()
	contracts[rootCid]["LiquidityProvide"] = native.LiquidityProvide()
	contracts[rootCid]["LiquidityWithdraw"] = native.LiquidityWithdraw()
	contracts[rootCid]["Swap"] = native.Swap()

	contracts[rootCid]["Mining"] = native.Mining()
	contracts[rootCid]["Count"] = native.Count()

	return contracts
}
