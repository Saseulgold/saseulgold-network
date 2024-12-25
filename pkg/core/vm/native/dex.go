package native

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/util"
)

func Mint() *Method {

	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Mint",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "name",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "symbol",
		"type":         "string",
		"maxlength":    5,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "supply",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	from := abi.Param("from")
	supply := abi.Param("supply")
	symbol := abi.Param("symbol")

	// Add fee check condition
	method.AddExecution(abi.Condition(
		abi.Gte(supply, MINT_FEE),
		"Supply amount must be greater than mint fee",
	))

	token_address := abi.HashMany([]interface{}{"qrc_20", abi.Param("from"), abi.Param("symbol")})

	method.AddExecution(abi.Condition(
		abi.Eq(abi.ReadUniversal(token_address, "owner", nil), nil),
		"The token can only be issued once.",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(abi.ReadUniversal(token_address, "supply", nil), nil),
		"The token can only be issued once.",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(abi.ReadUniversal(token_address, "symbol", nil), nil),
		"The token can only be issued once.",
	))

	method.AddExecution(abi.Condition(
		abi.Gt(supply, "0"),
		"The supply amount must be greater than 0.",
	))

	cond1 := abi.Condition(
		abi.Gte(abi.Len(symbol), 3),
		"The symbol string`s length must be greater than 2",
	)

	method.AddExecution(cond1)

	update_owner := abi.WriteUniversal(token_address, "owner", from)
	method.AddExecution(update_owner)

	update_supply := abi.WriteUniversal(token_address, "supply", supply)
	method.AddExecution(update_supply)

	update_symbol := abi.WriteUniversal(token_address, "symbol", symbol)
	method.AddExecution(update_symbol)

	owner_address := abi.HashMany([]interface{}{token_address, from})
	update_owner_balance := abi.WriteUniversal(owner_address, "balance", supply)

	method.AddExecution(update_owner_balance)

	network_fee_reserve := abi.ReadUniversal("network_fee_reserve", ZERO_ADDRESS, "0")
	rewserve_update := abi.WriteUniversal("network_fee_reserve", ZERO_ADDRESS, abi.PreciseAdd(network_fee_reserve, MINT_FEE, 0))

	method.AddExecution(rewserve_update)

	return method
}

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

	// Check if tokens exist
	method.AddExecution(abi.Condition(
		abi.Ne(abi.ReadUniversal(tokenA, "supply", nil), nil),
		"TokenA does not exist",
	))
	method.AddExecution(abi.Condition(
		abi.Ne(abi.ReadUniversal(tokenB, "supply", nil), nil),
		"TokenB does not exist",
	))

	// Check user's token balances
	fromBalanceA := abi.HashMany([]interface{}{tokenA, from})
	fromBalanceB := abi.HashMany([]interface{}{tokenB, from})

	userBalanceA := abi.ReadUniversal(fromBalanceA, "balance", "0")
	userBalanceB := abi.ReadUniversal(fromBalanceB, "balance", "0")

	// Verify sufficient balances
	method.AddExecution(abi.Condition(
		abi.Gte(userBalanceA, amountA),
		"Insufficient balance for tokenA",
	))
	method.AddExecution(abi.Condition(
		abi.Gte(userBalanceB, amountB),
		"Insufficient balance for tokenB",
	))

	// Deduct tokens from user's balance
	method.AddExecution(abi.WriteUniversal(fromBalanceA, "balance",
		abi.PreciseSub(userBalanceA, amountA, "0")))
	method.AddExecution(abi.WriteUniversal(fromBalanceB, "balance",
		abi.PreciseSub(userBalanceB, amountB, "0")))

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

func Transfer() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Transfer",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	// Add parameters
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "to",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amount",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	// Get parameters
	token_address := abi.Param("token_address")
	from := abi.Param("from")
	to := abi.Param("to")
	amount := abi.Param("amount")

	// Check if sender has enough balance for transfer fee
	fromBalanceSG := abi.ReadUniversal("balance", from, "0")
	method.AddExecution(abi.Condition(
		abi.Gt(fromBalanceSG, TRNF_FEE),
		"Balance is not enough for transfer fee",
	))

	// Check if sender and receiver are different
	method.AddExecution(abi.Condition(
		abi.Ne(from, to),
		"Sender and receiver address must be different",
	))

	// Check if token exists
	method.AddExecution(abi.Condition(
		abi.Ne(abi.ReadUniversal(token_address, "supply", nil), nil),
		"Token does not exist",
	))

	// Check if amount is valid
	method.AddExecution(abi.Condition(
		abi.Gt(amount, "0"),
		"Transfer amount must be greater than 0",
	))

	// Check if sender has enough balance
	fromAddressBalance := abi.HashMany([]interface{}{token_address, from})
	fromBalance := abi.ReadUniversal(fromAddressBalance, "balance", "0")
	method.AddExecution(abi.Condition(
		abi.Gte(fromBalance, amount),
		"Insufficient balance",
	))

	// Update balances
	newFromBalance := abi.PreciseSub(fromBalance, amount, "0")
	method.AddExecution(abi.WriteUniversal(fromAddressBalance, "balance", newFromBalance))

	toAddressBalance := abi.HashMany([]interface{}{token_address, to})
	toBalance := abi.ReadUniversal(toAddressBalance, "balance", "0")
	newToBalance := abi.PreciseAdd(toBalance, amount, "0")
	method.AddExecution(abi.WriteUniversal(toAddressBalance, "balance", newToBalance))

	// Deduct transfer fee from sender's balance
	newFromBalanceSG := abi.PreciseSub(fromBalance, TRNF_FEE, "0")
	method.AddExecution(abi.WriteUniversal("balance", from, newFromBalanceSG))

	// Add transfer fee to network fee reserve
	network_fee_reserve := abi.ReadUniversal("network_fee_reserve", ZERO_ADDRESS, "0")
	method.AddExecution(abi.WriteUniversal("network_fee_reserve", ZERO_ADDRESS, abi.PreciseAdd(network_fee_reserve, TRNF_FEE, "0")))

	return method
}
