package native

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/util"
)

func LiquidityProvider() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "LiquidityProvider",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "tokenA",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "tokenB",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amountA",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amountB",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	tokenA := abi.Param("tokenA")
	tokenB := abi.Param("tokenB")
	amountA := abi.Param("amountA")
	amountB := abi.Param("amountB")
	from := abi.Param("from")

	pairAddress := abi.HashMany([]interface{}{"pair", tokenA, tokenB})

	method.AddExecution(abi.Condition(
		abi.Eq(abi.ReadUniversal(pairAddress, "exists", nil), nil),
		"Liquidity pair does not exist.",
	))
	method.AddExecution(abi.WriteUniversal(pairAddress, "exists", true))

	currentReserveA := abi.ReadUniversal(pairAddress, "reserveA", "0")
	currentReserveB := abi.ReadUniversal(pairAddress, "reserveB", "0")

	newReserveA := abi.PreciseAdd(currentReserveA, amountA, "0")
	newReserveB := abi.PreciseAdd(currentReserveB, amountB, "0")

	method.AddExecution(abi.WriteUniversal(pairAddress, "reserveA", newReserveA))
	method.AddExecution(abi.WriteUniversal(pairAddress, "reserveB", newReserveB))

	liquidity := abi.PreciseSqrt(abi.PreciseMul(amountA, amountB, "0"), "0")

	liquidityToken := abi.HashMany([]interface{}{"liquidity", pairAddress, from})
	currentLiquidity := abi.ReadUniversal(liquidityToken, "balance", "0")
	newLiquidity := abi.PreciseAdd(currentLiquidity, liquidity, "0")

	method.AddExecution(abi.WriteUniversal(liquidityToken, "balance", newLiquidity))

	return method
}
