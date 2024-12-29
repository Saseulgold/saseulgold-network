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

	from := abi.Param("from")
	supply := abi.Param("supply")
	symbol := abi.Param("symbol")
	name := abi.Param("name")

	/**
	owner_balance_sg := abi.ReadUniversal("balance", from, "0")

	method.AddExecution(abi.Condition(
		abi.Gte(owner_balance_sg, MINT_FEE),
		"Balance is not enough for mint fee",
	))
	*/

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

	update_owner := abi.WriteUniversal(token_address, "owner", from)
	method.AddExecution(update_owner)

	update_name := abi.WriteUniversal(token_address, "name", name)
	method.AddExecution(update_name)

	update_supply := abi.WriteUniversal(token_address, "supply", supply)
	method.AddExecution(update_supply)

	update_symbol := abi.WriteUniversal(token_address, "symbol", symbol)
	method.AddExecution(update_symbol)

	balance_address := abi.HashMany(token_address, "balance")
	update_balance := abi.WriteUniversal(balance_address, from, supply)

	method.AddExecution(update_balance)

	network_fee_reserve := abi.ReadUniversal("network_fee_reserve", ZERO_ADDRESS, "0")
	network_fee_reserve_update := abi.WriteUniversal("network_fee_reserve", ZERO_ADDRESS, abi.PreciseAdd(network_fee_reserve, MINT_FEE, 0))

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
		abi.Ne(abi.ReadUniversal(tokenA, "supply", nil), nil),
		"TokenA does not exist",
	))
	method.AddExecution(abi.Condition(
		abi.Ne(abi.ReadUniversal(tokenB, "supply", nil), nil),
		"TokenB does not exist",
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
		"name":         "amount_in",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "minimum_amount_b",
		"type":         "string",
		"maxlength":    50,
		"requirements": false,
		"default":      nil,
	}))

	// Extract parameters
	tokenA := abi.Min(abi.Param("token_address_a"), abi.Param("token_address_b"))
	tokenB := abi.Max(abi.Param("token_address_a"), abi.Param("token_address_b"))

	inputAmount := abi.Param("amount_in")
	isInputTokenA := abi.Eq(abi.Param("token_address_a"), tokenA)

	amountA := abi.If(isInputTokenA, inputAmount, "0")
	amountB := abi.If(isInputTokenA, "0", inputAmount)

	from := abi.Param("from")

	// Ensure tokenA and tokenB are not the same
	method.AddExecution(abi.Condition(
		abi.Ne(tokenA, tokenB),
		"Token addresses must be different",
	))

	// Validate tokens exist
	method.AddExecution(abi.Condition(
		abi.Ne(abi.ReadUniversal(tokenA, "supply", nil), nil),
		"TokenA does not exist",
	))
	method.AddExecution(abi.Condition(
		abi.Ne(abi.ReadUniversal(tokenB, "supply", nil), nil),
		"TokenB does not exist",
	))

	// Validate reserves
	pairAddress := abi.HashMany("qrc_20_pair", tokenA, tokenB)
	totalLiquidity := abi.ReadUniversal(pairAddress, "total_liquidity", "0")

	method.AddExecution(abi.Condition(
		abi.Gt(totalLiquidity, "0"),
		"No liquidity available for fee distribution",
	))

	reserveA := abi.ReadUniversal(pairAddress, "reserve_a", "0")
	reserveB := abi.ReadUniversal(pairAddress, "reserve_b", "0")
	method.AddExecution(abi.Condition(
		abi.And(abi.Gt(reserveA, "0"), abi.Gt(reserveB, "0")),
		"Insufficient reserves for swap",
	))

	// Validate user balance for TokenA
	balance_address_a := abi.If(
		abi.Eq(tokenA, ZERO_ADDRESS),
		"balance",
		abi.HashMany(tokenA, "balance"),
	)
	balanceA := abi.ReadUniversal(balance_address_a, from, "0")
	method.AddExecution(abi.Condition(
		abi.Gte(balanceA, amountA),
		"Insufficient balance for TokenA",
	))

	// Calculate amount of TokenB to swap
	feeRate := C.SWAP_DEDUCT_RATE

	outputAmount := abi.If(isInputTokenA,
		// If input is tokenA, calculate tokenB output
		abi.PreciseDiv(
			abi.PreciseMul(amountA, reserveB, "0"),
			abi.PreciseAdd(reserveA, amountA, "0"),
			"0"),
		// If input is tokenB, calculate tokenA output
		abi.PreciseDiv(
			abi.PreciseMul(amountB, reserveA, "0"),
			abi.PreciseAdd(reserveB, amountB, "0"),
			"0"))

	// Apply fee to output amount
	outputAmountWithFee := abi.PreciseMul(outputAmount, feeRate, "0")

	// Ensure amountB meets minimum_amount_b
	/**
	method.AddExecution(abi.Condition(
		abi.Gte(amountBWithFee, minimumAmountB),
		"Output amount is less than the minimum specified",
	))
	*/

	// Update user's balances
	method.AddExecution(abi.WriteUniversal(balance_address_a, from,
		abi.PreciseSub(balanceA, amountA, "0")))

	balance_address_b := abi.If(
		abi.Eq(tokenB, ZERO_ADDRESS),
		"balance",
		abi.HashMany(tokenB, "balance"),
	)

	balanceB := abi.ReadUniversal(balance_address_b, from, "0")

	// Update user's balance with outputAmountWithFee instead of amountBWithFee
	method.AddExecution(abi.WriteUniversal(balance_address_b, from,
		abi.PreciseAdd(balanceB, outputAmountWithFee, "0")))

	// Update reserves using outputAmountWithFee
	newReserveA := abi.If(isInputTokenA,
		abi.PreciseAdd(reserveA, amountA, "0"),
		abi.PreciseSub(reserveA, outputAmountWithFee, "0"))

	newReserveB := abi.If(isInputTokenA,
		abi.PreciseSub(reserveB, outputAmountWithFee, "0"),
		abi.PreciseAdd(reserveB, amountB, "0"))

	method.AddExecution(abi.WriteUniversal(pairAddress, "reserve_a", newReserveA))
	method.AddExecution(abi.WriteUniversal(pairAddress, "reserve_b", newReserveB))

	// Calculate fee using outputAmount and outputAmountWithFee
	fee := abi.PreciseSub(outputAmount, outputAmountWithFee, "0")

	// Read current accumulated reward per unit
	current_reward_per_unit := abi.ReadUniversal(pairAddress, "accumulated_reward_per_unit", "0")

	// Calculate new reward per unit and add to accumulated total
	new_reward_per_unit := abi.PreciseDiv(fee, totalLiquidity, "0")
	accumulated_reward_per_unit := abi.PreciseAdd(current_reward_per_unit, new_reward_per_unit, "0")

	// Store accumulated reward per unit
	method.AddExecution(abi.WriteUniversal(pairAddress, "accumulated_reward_per_unit", accumulated_reward_per_unit))

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

	tokenA := abi.Param("token_address_a")
	tokenB := abi.Param("token_address_b")

	pairAddress := abi.HashMany("qrc_20_pair", abi.Min(tokenA, tokenB), abi.Max(tokenA, tokenB))

	// Check if pair exists first
	exists := abi.ReadUniversal(pairAddress, "exists", nil)
	method.AddExecution(abi.Condition(
		abi.Eq(exists, true),
		"Pair does not exist",
	))

	var response interface{}

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
	oneTokenA := "1000000000000000000" // 1 with 18 decimals
	numeratorA := abi.PreciseMul(oneTokenA, reserveB, "0")
	denominatorA := abi.PreciseAdd(reserveA, oneTokenA, "0")
	rateA := abi.PreciseMul(abi.PreciseDiv(numeratorA, denominatorA, "0"), feeRate, "0")
	response = abi.Set(response, "rate_a_to_b", rateA)

	// For 1 tokenB -> tokenA
	oneTokenB := "1000000000000000000" // 1 with 18 decimals
	numeratorB := abi.PreciseMul(oneTokenB, reserveA, "0")
	denominatorB := abi.PreciseAdd(reserveB, oneTokenB, "0")
	rateB := abi.PreciseMul(abi.PreciseDiv(numeratorB, denominatorB, "0"), feeRate, "0")
	response = abi.Set(response, "rate_b_to_a", rateB)

	// Get liquidity info
	totalLiquidity := abi.ReadUniversal(pairAddress, "total_liquidity", "0")
	response = abi.Set(response, "total_liquidity", totalLiquidity)

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
		"name":         "liquidity_amount",
		"type":         "string",
		"maxlength":    50,
		"requirements": true,
	}))

	// Get and sort token addresses
	tokenA := abi.Min(abi.Param("token_address_a"), abi.Param("token_address_b"))
	tokenB := abi.Max(abi.Param("token_address_a"), abi.Param("token_address_b"))
	liquidityAmount := abi.Param("liquidity_amount")
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
	rewards := abi.PreciseMul(liquidityAmount, accumulated_reward_per_unit, "0")

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
