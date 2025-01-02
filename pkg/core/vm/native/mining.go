package native

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
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

	method.AddParameter(NewParameter(map[string]interface{}{
		"name":         "mining_start",
		"type":         "string",
		"maxlength":    256,
		"requirements": true,
	}))

	// Get parameters and read states
	from := abi.Param("from")
	epoch := abi.Param("epoch")

	nonce := abi.Param("nonce")
	calculatedHash := abi.Param("calculated_hash")
	// difficulty := abi.ReadUniversal("network_difficulty", ZERO_ADDRESS, "1")

	// Validate hash
	// limit := abi.HashLimit(difficulty)
	dhash := abi.HashMany(epoch, nonce)

	method.AddExecution(abi.Condition(
		abi.Eq(calculatedHash, dhash),
		"Invalid nonce and hash.",
	))

	/**
	method.AddExecution(abi.Condition(
		abi.Lt(calculatedHash, limit),
		"Hash limit was not satisfied.",
	))
	*/

	lastRewarded := abi.ReadUniversal("lastRewarded", ZERO_ADDRESS, "0")
	lastRewarded = abi.Check(lastRewarded, "lastRewarded")

	// current := abi.AsString(util.Utime())
	// current = abi.Check(current, "current")
	current := "1735782161068000"

	timeDiff := abi.PreciseSub(current, lastRewarded, "0")
	timeDiff = abi.Min(timeDiff, "100000")
	timeDiff = abi.Check(timeDiff, "timeDiff")

	reward := abi.PreciseMul(timeDiff, REWARD_PER_SECOND, "0")
	reward = abi.PreciseDiv(reward, "1000", "0")
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

	return method
}
