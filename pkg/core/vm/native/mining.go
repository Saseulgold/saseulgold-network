package native

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/util"
)

func Mining() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Mining",
		"version": "1",
		"space":   RootSpace(),
		"writer":  ZERO_ADDRESS,
	})

	// Add parameters
	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "nonce",
		"type":         "string",
		"maxlength":    256,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "epoch",
		"type":         "string",
		"maxlength":    256,
		"requirements": true,
	}))

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "calculated_hash",
		"type":         "string",
		"maxlength":    256,
		"requirements": true,
	}))

	// Get parameters and read states
	from := abi.Param("from")
	epoch := abi.Param("epoch")

	// Add block hash verification for current and previous 5 blocks
	blockHeight := abi.AsString(epoch)
	block := abi.GetBlock(blockHeight, "full")
	blockHash := abi.Get(block, "hash", nil)

	method.AddExecution(abi.Condition(
		abi.Ne(blockHash, nil),
		"Invalid block height or block not found.",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(blockHash, epoch),
		"Block hash does not match with epoch.",
	))

	// Verify previous 5 blocks
	prevBlock1 := abi.GetBlock(abi.PreciseSub(blockHeight, "1", "0"), "full")
	prevBlock2 := abi.GetBlock(abi.PreciseSub(blockHeight, "2", "0"), "full")
	prevBlock3 := abi.GetBlock(abi.PreciseSub(blockHeight, "3", "0"), "full")
	prevBlock4 := abi.GetBlock(abi.PreciseSub(blockHeight, "4", "0"), "full")
	prevBlock5 := abi.GetBlock(abi.PreciseSub(blockHeight, "5", "0"), "full")

	prevHash1 := abi.Get(prevBlock1, "hash", nil)
	prevHash2 := abi.Get(prevBlock2, "hash", nil)
	prevHash3 := abi.Get(prevBlock3, "hash", nil)
	prevHash4 := abi.Get(prevBlock4, "hash", nil)
	prevHash5 := abi.Get(prevBlock5, "hash", nil)

	method.AddExecution(abi.Condition(
		abi.And(
			abi.Ne(prevHash1, nil),
			abi.Ne(prevHash2, nil),
			abi.Ne(prevHash3, nil),
			abi.Ne(prevHash4, nil),
			abi.Ne(prevHash5, nil),
		),
		"Previous blocks not found.",
	))

	// Verify block chain continuity
	method.AddExecution(abi.Condition(
		abi.Eq(abi.Get(block, "prev_hash", nil), prevHash1),
		"Block chain continuity broken at height -1",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(abi.Get(prevBlock1, "prev_hash", nil), prevHash2),
		"Block chain continuity broken at height -2",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(abi.Get(prevBlock2, "prev_hash", nil), prevHash3),
		"Block chain continuity broken at height -3",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(abi.Get(prevBlock3, "prev_hash", nil), prevHash4),
		"Block chain continuity broken at height -4",
	))

	method.AddExecution(abi.Condition(
		abi.Eq(abi.Get(prevBlock4, "prev_hash", nil), prevHash5),
		"Block chain continuity broken at height -5",
	))

	nonce := abi.Param("nonce")
	calculatedHash := abi.Param("calculated_hash")
	// difficulty := abi.ReadUniversal("network_difficulty", ZERO_ADDRESS, "1")

	limit := C.NETWORK_DIFF
	dhash := abi.HashMany(epoch, nonce)

	method.AddExecution(abi.Condition(
		abi.Eq(calculatedHash, dhash),
		"Invalid nonce and hash.",
	))

	method.AddExecution(abi.Condition(
		abi.HashLimitOk(calculatedHash, limit),
		"Hash limit was not satisfied.",
	))

	lastRewarded := abi.ReadUniversal("lastRewarded", ZERO_ADDRESS, "0")
	lastRewarded = abi.Check(lastRewarded, "lastRewarded")

	current := abi.SUtime()
	current = abi.Check(current, "current")

	addressLastRewarded := abi.Check(abi.ReadUniversal("lastRewarded", from, nil), "lastrewardfrom")

	method.AddExecution(abi.Condition(
		abi.Or(
			abi.Eq(addressLastRewarded, nil),
			abi.Gt(abi.Check(abi.PreciseSub(current, addressLastRewarded, "0"), "ftd"), "120000000"),
		),
		"Queue is empty.",
	))

	timeDiff := abi.PreciseSub(current, lastRewarded, "0")
	timeDiff = abi.PreciseDiv(timeDiff, "1000", "0")
	// timeDiff = abi.Min(timeDiff, "10000000")

	timeDiff = abi.Check(timeDiff, "timeDiff")
	total_supply := abi.ReadUniversal("network_supply", ZERO_ADDRESS, C.INITIAL_SUPPLY)

	era := abi.Era(abi.PreciseDiv(total_supply, MULTIPLIER, "0"))
	era = abi.AsString(era)
	era = abi.Check(era, "era")

	unit := abi.PreciseDiv(REWARD_PER_SECOND, abi.PrecisePow("2", era, "0"), "0")
	unit = abi.Check(unit, "unit")

	method.AddExecution(abi.Check(unit, "unit"))

	reward := abi.PreciseMul(timeDiff, unit, "0")
	reward = abi.PreciseDiv(reward, "1000", "0")
	reward = abi.Min(reward, "1680000000000000000000000")

	reward = abi.Check(reward, "reward")

	method.AddExecution(
		abi.Condition(
			abi.Gt(reward, "0"),
			"Reward is 0.",
		),
	)

	balance := abi.ReadUniversal("balance", from, "0")
	newBalance := abi.PreciseAdd(balance, reward, "0")

	method.AddExecution(abi.WriteUniversal("balance", from, newBalance))
	method.AddExecution(abi.WriteUniversal("lastRewarded", ZERO_ADDRESS, current))

	method.AddExecution(abi.WriteUniversal("lastRewarded", ZERO_ADDRESS, current))
	method.AddExecution(abi.WriteUniversal("lastRewarded", from, current))

	new_total_supply := abi.PreciseAdd(total_supply, reward, "0")

	method.AddExecution(abi.WriteUniversal("network_supply", ZERO_ADDRESS, new_total_supply))

	difficulty := abi.ReadUniversal("network_difficulty", ZERO_ADDRESS, "2000")
	difficulty = abi.Min(abi.PreciseAdd(difficulty, "12", "0"), "4875")
	method.AddExecution(abi.WriteUniversal("network_difficulty", ZERO_ADDRESS, difficulty))

	era = abi.Check(era, "era")
	method.AddExecution(abi.Check(era, "era"))

	return method
}
