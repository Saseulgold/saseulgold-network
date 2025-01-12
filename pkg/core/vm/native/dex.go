package native

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
	C "hello/pkg/core/config"
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

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "icon_url",
		"type":         "string",
		"maxlength":    256,
		"default":      "",
		"requirements": false,
	}))

	from := abi.Param("from")
	supply := abi.Param("supply")
	symbol := abi.Param("symbol")
	name := abi.Param("name")
	icon_url := abi.Param("icon_url")

	owner_balance_sg := abi.ReadUniversal("balance", from, "0")

	method.AddExecution(abi.Condition(
		abi.Gte(owner_balance_sg, MINT_FEE),
		"Balance is not enough for mint fee",
	))

	token_address := abi.HashMany("qrc_20", abi.Param("from"), abi.Param("symbol"))

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

	update_token_list := abi.WriteUniversal("tokens", token_address, symbol)
	method.AddExecution(update_token_list)

	update_owner := abi.WriteUniversal(token_address, "owner", from)
	method.AddExecution(update_owner)

	update_name := abi.WriteUniversal(token_address, "name", name)
	method.AddExecution(update_name)

	update_supply := abi.WriteUniversal(token_address, "supply", supply)
	method.AddExecution(update_supply)

	update_symbol := abi.WriteUniversal(token_address, "symbol", symbol)
	method.AddExecution(update_symbol)

	update_icon_url := abi.WriteUniversal(token_address, "icon_url", icon_url)
	method.AddExecution(update_icon_url)

	balance_address := abi.HashMany(token_address, "balance")
	update_balance := abi.WriteUniversal(balance_address, from, supply)

	method.AddExecution(update_balance)

	owner_balance_update := abi.WriteUniversal("balance", from, abi.PreciseSub(owner_balance_sg, MINT_FEE, "0"))
	method.AddExecution(owner_balance_update)

	network_fee_reserve := abi.ReadUniversal("network_fee_reserve", ZERO_ADDRESS, "0")
	network_fee_reserve_update := abi.WriteUniversal("network_fee_reserve", ZERO_ADDRESS, abi.PreciseAdd(network_fee_reserve, MINT_FEE, "0"))

	method.AddExecution(network_fee_reserve_update)

	return method
}

func LiquidityProvide() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "LiquidityProvide",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address_a",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address_b",
		"type":         "string",
		"maxlength":    80,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amount_a",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amount_b",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	tokenA := abi.Min(abi.Param("token_address_a"), abi.Param("token_address_b"))
	tokenB := abi.Max(abi.Param("token_address_a"), abi.Param("token_address_b"))

	amountA := abi.If(
		abi.Eq(tokenA, abi.Param("token_address_a")),
		abi.Param("amount_a"),
		abi.Param("amount_b"),
	)
	amountB := abi.If(
		abi.Eq(tokenB, abi.Param("token_address_b")),
		abi.Param("amount_b"),
		abi.Param("amount_a"),
	)

	from := abi.Param("from")

	// Check if tokens exist
	method.AddExecution(abi.Condition(
		abi.Or(
			abi.Eq(tokenA, ZERO_ADDRESS),
			abi.Ne(abi.ReadUniversal(tokenA, "supply", nil), nil),
		),
		"TokenA does not exist",
	))
	method.AddExecution(abi.Condition(
		abi.Or(
			abi.Eq(tokenB, ZERO_ADDRESS),
			abi.Ne(abi.ReadUniversal(tokenB, "supply", nil), nil),
		),
		"TokenB does not exist",
	))

	method.AddExecution(abi.Condition(
		abi.Ne(tokenA, tokenB),
		"Token addresses must be different",
	))

	balance_address_a := abi.If(
		abi.Eq(tokenA, ZERO_ADDRESS),
		"balance",
		abi.HashMany(tokenA, "balance"),
	)

	balance_address_b := abi.If(
		abi.Eq(tokenB, ZERO_ADDRESS),
		"balance",
		abi.HashMany(tokenB, "balance"),
	)

	userBalanceA := abi.ReadUniversal(balance_address_a, from, "0")
	userBalanceB := abi.ReadUniversal(balance_address_b, from, "0")

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
	method.AddExecution(abi.WriteUniversal(balance_address_a, from,
		abi.PreciseSub(userBalanceA, amountA, "0")))
	method.AddExecution(abi.WriteUniversal(balance_address_b, from,
		abi.PreciseSub(userBalanceB, amountB, "0")))

	pairAddress := abi.HashMany("qrc_20_pair", tokenA, tokenB)

	currentReserveA := abi.ReadUniversal(pairAddress, "reserve_a", "0")
	currentReserveB := abi.ReadUniversal(pairAddress, "reserve_b", "0")

	newReserveA := abi.PreciseAdd(currentReserveA, amountA, "0")
	newReserveB := abi.PreciseAdd(currentReserveB, amountB, "0")

	method.AddExecution(abi.WriteUniversal(pairAddress, "reserve_a", newReserveA))
	method.AddExecution(abi.WriteUniversal(pairAddress, "reserve_b", newReserveB))

	exists := abi.ReadUniversal(pairAddress, "exists", nil)

	pair_a := abi.If(
		abi.Eq(exists, nil),
		abi.WriteUniversal(pairAddress, "token_address_a", tokenA),
		abi.Check(tokenA, "tokenA"),
	)

	pair_b := abi.If(
		abi.Eq(exists, nil),
		abi.WriteUniversal(pairAddress, "token_address_b", tokenB),
		abi.Check(tokenB, "tokenB"),
	)

	method.AddExecution(pair_a)
	method.AddExecution(pair_b)

	liquidity := abi.If(
		abi.Eq(exists, nil),
		abi.PreciseSqrt(abi.PreciseMul(amountA, amountB, "0"), "0"),
		abi.Min(
			abi.PreciseDiv(abi.PreciseMul(amountA, currentReserveB, "0"), currentReserveA, "0"),
			abi.PreciseDiv(abi.PreciseMul(amountB, currentReserveA, "0"), currentReserveB, "0"),
		),
	)

	currentLiquidity := abi.ReadUniversal(pairAddress, from, "0")
	newLiquidity := abi.PreciseAdd(currentLiquidity, liquidity, "0")

	method.AddExecution(abi.WriteUniversal(pairAddress, from, newLiquidity))

	totalLiquidity := abi.ReadUniversal(pairAddress, "total_liquidity", "0")
	newTotalLiquidity := abi.PreciseAdd(totalLiquidity, liquidity, "0")

	method.AddExecution(abi.WriteUniversal(pairAddress, "total_liquidity", newTotalLiquidity))
	method.AddExecution(abi.WriteUniversal(pairAddress, "exists", true))

	symbolA := abi.ReadUniversal(tokenA, "symbol", "")
	symbolB := abi.ReadUniversal(tokenB, "symbol", "")

	method.AddExecution(abi.WriteUniversal("pairs", pairAddress, abi.Concat(symbolA, "/", symbolB)))

	avg_return := abi.ReadUniversal(pairAddress, "accumulated_reward_per_unit", "0")

	new_avg_return := abi.PreciseDiv(
		abi.PreciseMul(avg_return, totalLiquidity, "0"),
		newTotalLiquidity,
		"0",
	)

	updateAvgReturn := abi.WriteUniversal(pairAddress, "accumulated_reward_per_unit", new_avg_return)
	method.AddExecution(updateAvgReturn)

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
	balance_address := abi.HashMany(token_address, "balance")
	fromBalance := abi.ReadUniversal(balance_address, from, "0")
	method.AddExecution(abi.Condition(
		abi.Gte(fromBalance, amount),
		"Insufficient balance",
	))

	// Update balances
	method.AddExecution(abi.WriteUniversal(balance_address, from,
		abi.PreciseSub(fromBalance, amount, "0")))

	toBalance := abi.ReadUniversal(balance_address, to, "0")
	method.AddExecution(abi.WriteUniversal(balance_address, to,
		abi.PreciseAdd(toBalance, amount, "0")))

	// Deduct transfer fee from sender's balance
	newFromBalanceSG := abi.PreciseSub(fromBalanceSG, TRNF_FEE, "0")
	method.AddExecution(abi.WriteUniversal("balance", from, newFromBalanceSG))

	// Add transfer fee to network fee reserve
	network_fee_reserve := abi.ReadUniversal("network_fee_reserve", ZERO_ADDRESS, "0")
	method.AddExecution(abi.WriteUniversal("network_fee_reserve", ZERO_ADDRESS, abi.PreciseAdd(network_fee_reserve, TRNF_FEE, "0")))

	return method
}

func BalanceOf() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "BalanceOf",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	// Add token_address parameter
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address",
		"type":         "string",
		"maxlength":    128,
		"requirements": true,
	}))

	// Add account parameter
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "address",
		"type":         "string",
		"maxlength":    128,
		"requirements": true,
	}))

	token_address := abi.Param("token_address")
	address := abi.Param("address")

	// Check if token exists
	method.AddExecution(abi.Condition(
		abi.Ne(abi.ReadUniversal(token_address, "supply", nil), nil),
		"Token does not exist",
	))

	// Get balance address
	balance_address := abi.HashMany(token_address, "balance")

	// Read and return balance
	balance := abi.ReadUniversal(balance_address, address, "0")
	method.AddExecution(abi.Response(balance))

	return method
}

func GetPairInfo() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "GetPairInfo",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "pair_address",
		"type":         "string",
		"maxlength":    64,
		"requirements": false,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address_a",
		"type":         "string",
		"maxlength":    64,
		"requirements": false,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address_b",
		"type":         "string",
		"maxlength":    64,
		"requirements": false,
	}))

	tokenA := abi.Min(abi.Param("token_address_a"), abi.Param("token_address_b"))
	tokenB := abi.Max(abi.Param("token_address_a"), abi.Param("token_address_b"))

	tokenA = abi.If(
		abi.Ne(abi.Param("pair_address"), nil),
		abi.ReadUniversal(abi.Param("pair_address"), "token_address_a", nil),
		tokenA,
	)

	tokenB = abi.If(
		abi.Ne(abi.Param("pair_address"), nil),
		abi.ReadUniversal(abi.Param("pair_address"), "token_address_b", nil),
		tokenB,
	)

	method.AddExecution(
		abi.Condition(
			abi.Or(
				abi.And(
					abi.Ne(tokenA, nil),
					abi.Ne(tokenB, nil),
				),
				abi.Ne(abi.Param("pair_address"), nil),
			),
			"Token addresses or pair address must be provided",
		),
	)

	pairAddress := abi.If(
		abi.Eq(abi.Param("pair_address"), nil),
		abi.HashMany("qrc_20_pair", tokenA, tokenB),
		abi.Param("pair_address"),
	)

	// Check if pair exists first
	exists := abi.ReadUniversal(pairAddress, "exists", nil)
	method.AddExecution(abi.Condition(
		abi.Eq(exists, true),
		abi.EncodeJSON("Pair does not exist"),
	))

	var response interface{}

	response = abi.Set(response, "pair_address", pairAddress)

	// Get reserves
	reserveA := abi.ReadUniversal(pairAddress, "reserve_a", "0")
	response = abi.Set(response, "reserve_a", reserveA)

	reserveB := abi.ReadUniversal(pairAddress, "reserve_b", "0")
	response = abi.Set(response, "reserve_b", reserveB)

	// Calculate swap rates
	// rate_a: How much of tokenB you get for 1 tokenA
	// rate_b: How much of tokenA you get for 1 tokenB

	addressA := abi.Min(tokenA, tokenB)
	addressB := abi.Max(tokenA, tokenB)

	// Apply fee rate to calculations
	feeRate := C.SWAP_DEDUCT_RATE

	response = abi.Set(response, "address_a", addressA)
	response = abi.Set(response, "address_b", addressB)

	// For 1 tokenA -> tokenB
	oneTokenA := MULTIPLIER
	numeratorA := abi.PreciseMul(oneTokenA, reserveB, "0")
	denominatorA := abi.PreciseAdd(reserveA, oneTokenA, "0")
	rateA := abi.PreciseMul(abi.PreciseDiv(numeratorA, denominatorA, "0"), feeRate, "0")
	response = abi.Set(response, "rate_a_to_b", rateA)

	// For 1 tokenB -> tokenA
	oneTokenB := MULTIPLIER
	numeratorB := abi.PreciseMul(oneTokenB, reserveA, "0")
	denominatorB := abi.PreciseAdd(reserveB, oneTokenB, "0")
	rateB := abi.PreciseMul(abi.PreciseDiv(numeratorB, denominatorB, "0"), feeRate, "0")
	response = abi.Set(response, "rate_b_to_a", rateB)

	// Get liquidity info
	totalLiquidity := abi.ReadUniversal(pairAddress, "total_liquidity", "0")
	response = abi.Set(response, "total_liquidity", totalLiquidity)

	// Get volume
	volume := abi.ReadUniversal(pairAddress, "volume", "0")
	response = abi.Set(response, "volume", volume)

	// Get accumulated rewards per unit
	accumulatedRewardPerUnit := abi.ReadUniversal(pairAddress, "accumulated_reward_per_unit", "0")
	response = abi.Set(response, "accumulated_reward_per_unit", accumulatedRewardPerUnit)

	// Encode and return response
	response = abi.EncodeJSON(response)
	method.AddExecution(abi.Response(response))

	return method
}

func LiquidityWithdraw() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "LiquidityWithdraw",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	// Add parameters
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address_a",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address_b",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "lp_amount",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	// Get and sort token addresses
	tokenA := abi.Min(abi.Param("token_address_a"), abi.Param("token_address_b"))
	tokenB := abi.Max(abi.Param("token_address_a"), abi.Param("token_address_b"))
	liquidityAmount := abi.Param("lp_amount")
	from := abi.Param("from")

	// Validate withdrawal amount
	method.AddExecution(abi.Condition(
		abi.Gt(liquidityAmount, "0"),
		"Withdrawal amount must be greater than 0",
	))

	pairAddress := abi.HashMany("qrc_20_pair", tokenA, tokenB)

	// Validate pair exists
	method.AddExecution(abi.Condition(
		abi.Eq(abi.ReadUniversal(pairAddress, "exists", nil), true),
		"Liquidity pair does not exist",
	))

	// Check user's liquidity balance
	userLiquidity := abi.ReadUniversal(pairAddress, from, "0")
	method.AddExecution(abi.Condition(
		abi.Gte(userLiquidity, liquidityAmount),
		"Insufficient liquidity balance",
	))

	// Get pool state
	reserveA := abi.ReadUniversal(pairAddress, "reserve_a", "0")
	reserveB := abi.ReadUniversal(pairAddress, "reserve_b", "0")
	totalLiquidity := abi.ReadUniversal(pairAddress, "total_liquidity", "0")

	// Calculate withdrawal amounts
	amountA := abi.PreciseDiv(abi.PreciseMul(liquidityAmount, reserveA, "0"), totalLiquidity, "0")
	amountB := abi.PreciseDiv(abi.PreciseMul(liquidityAmount, reserveB, "0"), totalLiquidity, "0")

	// Calculate rewards
	accumulated_reward_per_unit := abi.ReadUniversal(pairAddress, "accumulated_reward_per_unit", "0")
	// accumulated_reward_per_unit = abi.PreciseDiv(accumulated_reward_per_unit, MULTIPLIER)

	rewards := abi.PreciseMul(liquidityAmount, accumulated_reward_per_unit, "0")
	rewards = abi.PreciseDiv(rewards, MULTIPLIER, "0")

	// Get balance addresses
	balance_address_a := abi.If(
		abi.Eq(tokenA, ZERO_ADDRESS),
		"balance",
		abi.HashMany(tokenA, "balance"),
	)
	balance_address_b := abi.If(
		abi.Eq(tokenB, ZERO_ADDRESS),
		"balance",
		abi.HashMany(tokenB, "balance"),
	)

	// Update user balances
	userBalanceA := abi.ReadUniversal(balance_address_a, from, "0")
	userBalanceB := abi.ReadUniversal(balance_address_b, from, "0")

	// Update token A balance
	method.AddExecution(abi.WriteUniversal(balance_address_a, from,
		abi.PreciseAdd(userBalanceA, amountA, "0")))

	// Update token B balance with rewards
	rewardBalanceB := abi.PreciseAdd(userBalanceB, rewards, "0")
	method.AddExecution(abi.WriteUniversal(balance_address_b, from,
		abi.PreciseAdd(rewardBalanceB, amountB, "0")))

	// Update reserves
	method.AddExecution(abi.WriteUniversal(pairAddress, "reserve_a",
		abi.PreciseSub(reserveA, amountA, "0")))
	method.AddExecution(abi.WriteUniversal(pairAddress, "reserve_b",
		abi.PreciseSub(reserveB, amountB, "0")))

	// Update liquidity state
	method.AddExecution(abi.WriteUniversal(pairAddress, from,
		abi.PreciseSub(userLiquidity, liquidityAmount, "0")))

	// Update total liquidity
	method.AddExecution(abi.WriteUniversal(pairAddress, "total_liquidity",
		abi.PreciseSub(totalLiquidity, liquidityAmount, "0")))

	return method
}

func BalanceOfLP() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "request",
		"name":    "BalanceOfLP",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	// Add pair_address parameter
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "pair_address",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))

	// Add account parameter
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "address",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))

	pair_address := abi.Param("pair_address")
	address := abi.Param("address")

	// Check if token exists
	method.AddExecution(abi.Condition(
		abi.Ne(abi.ReadUniversal(pair_address, "total_liquidity", nil), nil),
		"Token Pair does not exist",
	))

	// Read and return lp token balance
	balance := abi.ReadUniversal(pair_address, address, "0")
	method.AddExecution(abi.Response(balance))

	return method
}

func Swap() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Swap",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	// Parameters for the swap
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address_in",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "token_address_out",
		"type":         "string",
		"maxlength":    64,
		"requirements": true,
	}))
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "amount_in",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "minimum_amount_out",
		"type":         "string",
		"maxlength":    50,
		"requirements": false,
		"default":      nil,
	}))

	// Extract parameters
	tokenAddressIn := abi.Param("token_address_in")
	tokenAddressOut := abi.Param("token_address_out")

	tokenAddressA := abi.Min(tokenAddressIn, tokenAddressOut)
	tokenAddressB := abi.Max(tokenAddressIn, tokenAddressOut)

	amountIn := abi.Param("amount_in")

	minimumAmountOut := abi.Param("minimum_amount_out")
	from := abi.Param("from")

	// Determine input and output tokens based on token addresses
	isAToB := abi.If(abi.Eq(tokenAddressIn, abi.Min(tokenAddressA, tokenAddressB)), true, false) // Define a consistent ordering

	inputTokenAddress := abi.If(isAToB, tokenAddressA, tokenAddressB)
	outputTokenAddress := abi.If(isAToB, tokenAddressB, tokenAddressA)

	reserveInputAttr := abi.If(isAToB, "reserve_a", "reserve_b")
	reserveOutputAttr := abi.If(isAToB, "reserve_b", "reserve_a")

	// Ensure input and output tokens are different
	method.AddExecution(abi.Condition(
		abi.Ne(inputTokenAddress, outputTokenAddress),
		"Input and output token addresses must be different",
	))

	// Validate tokens exist
	method.AddExecution(abi.Condition(
		abi.Or(
			abi.Eq(inputTokenAddress, ZERO_ADDRESS),
			abi.Ne(abi.ReadUniversal(inputTokenAddress, "supply", nil), nil),
		),
		"Input token does not exist",
	))
	method.AddExecution(abi.Condition(
		abi.Or(
			abi.Eq(outputTokenAddress, ZERO_ADDRESS),
			abi.Ne(abi.ReadUniversal(outputTokenAddress, "supply", nil), nil),
		),
		"Output token does not exist",
	))

	// Validate reserves
	pairAddress := abi.HashMany("qrc_20_pair", tokenAddressA, tokenAddressB)
	totalLiquidity := abi.ReadUniversal(pairAddress, "total_liquidity", "0")

	method.AddExecution(abi.Condition(
		abi.Gt(totalLiquidity, "0"),
		"No liquidity available for fee distribution",
	))

	reserveInput := abi.ReadUniversal(pairAddress, reserveInputAttr, "0")
	reserveOutput := abi.ReadUniversal(pairAddress, reserveOutputAttr, "0")

	method.AddExecution(abi.Condition(
		abi.And(abi.Gt(reserveInput, "0"), abi.Gt(reserveOutput, "0")),
		"Insufficient reserves for swap",
	))

	// Validate user balance for Input Token
	balanceInputAddress := abi.If(
		abi.Eq(inputTokenAddress, ZERO_ADDRESS),
		"balance",
		abi.HashMany(inputTokenAddress, "balance"),
	)
	balanceInput := abi.ReadUniversal(balanceInputAddress, from, "0")
	method.AddExecution(abi.Condition(
		abi.Gte(balanceInput, amountIn),
		"Insufficient balance for Input Token",
	))

	// Calculate amount of Output Token to swap
	feeRateNumerator := "997"
	feeRateDenominator := "1000"

	// denominator = reserveInput + amountInWithFee
	denominator := abi.PreciseAdd(reserveInput, amountIn, "0")

	// calculatedOutputAmount = (amountInWithFee * reserveOutput) / denominator
	outputAmount := abi.PreciseDiv(
		abi.PreciseMul(amountIn, reserveOutput, "0"),
		denominator,
		"0",
	)

	outputAmountWithFee := abi.PreciseDiv(abi.PreciseMul(outputAmount, feeRateNumerator, "0"), feeRateDenominator, "0")
	outputAmountWithFee = abi.Check(outputAmountWithFee, "outputAmountWithFee")

	method.AddExecution(abi.Condition(
		abi.Or(
			abi.Eq(minimumAmountOut, nil),
			abi.Gte(outputAmountWithFee, minimumAmountOut),
		),
		"Output amount is less than the minimum specified",
	))

	// Update user's input token balance
	method.AddExecution(abi.WriteUniversal(balanceInputAddress, from,
		abi.PreciseSub(balanceInput, amountIn, "0"),
	))

	// Update user's output token balance
	balanceOutputAddress := abi.If(
		abi.Eq(outputTokenAddress, ZERO_ADDRESS),
		"balance",
		abi.HashMany(outputTokenAddress, "balance"),
	)
	balanceOutput := abi.ReadUniversal(balanceOutputAddress, from, "0")
	method.AddExecution(abi.WriteUniversal(balanceOutputAddress, from,
		abi.PreciseAdd(balanceOutput, outputAmountWithFee, "0"),
	))

	// Update reserves based on swap direction
	newReserveInput := abi.PreciseAdd(reserveInput, amountIn, "0")
	newReserveOutput := abi.PreciseSub(reserveOutput, outputAmountWithFee, "0")

	method.AddExecution(abi.WriteUniversal(pairAddress, reserveInputAttr, newReserveInput))
	method.AddExecution(abi.WriteUniversal(pairAddress, reserveOutputAttr, newReserveOutput))

	// Calculate and distribute fee
	fee := abi.PreciseSub(outputAmount, outputAmountWithFee, "0")

	currentRewardPerUnit := abi.ReadUniversal(pairAddress, "accumulated_reward_per_unit", "0")

	newRewardPerUnit := abi.PreciseDiv(
		abi.PreciseMul(fee, MULTIPLIER, "0"),
		totalLiquidity,
		"0",
	)

	accumulatedRewardPerUnit := abi.PreciseAdd(currentRewardPerUnit, newRewardPerUnit, "0")

	// Store accumulated reward per unit
	method.AddExecution(abi.WriteUniversal(pairAddress, "accumulated_reward_per_unit", accumulatedRewardPerUnit))

	volume := abi.If(
		isAToB,
		amountIn,
		outputAmount,
	)

	acc := abi.ReadUniversal(pairAddress, "volume", "0")
	acc = abi.PreciseAdd(acc, volume, "0")

	method.AddExecution(abi.WriteUniversal(pairAddress, "volume", acc))
	return method
}
