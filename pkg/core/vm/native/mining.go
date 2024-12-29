package native

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/config"
	. "hello/pkg/core/model"
	"hello/pkg/util"
	. "hello/pkg/util"
)

func Submit() *Method {
	method := NewMethod(map[string]interface{}{
		"type":    "contract",
		"name":    "Submit",
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
		"name":         "calculated_hash",
		"type":         "string",
		"maxlength":    256,
		"requirements": true,
	}))

	// Get parameters and read states
	from := abi.Param("from")
	epoch := abi.Param("epoch")

	fromBalance := abi.ReadUniversal("balance", from, "0")

	nonce := abi.Param("nonce")
	calculatedHash := abi.Param("calculated_hash")
	difficulty := abi.ReadLocal("rewardDifficulty", ZERO_ADDRESS, "2000")

	method.AddExecution(fromBalance)

	// Validate hash
	limit := abi.HashLimit(difficulty)
	dhash := abi.HashMany(epoch, nonce)

	method.AddExecution(abi.Condition(
		abi.Eq(calculatedHash, dhash),
		"Invalid nonce and hash.",
	))

	method.AddExecution(abi.Condition(
		abi.Lt(calculatedHash, limit),
		"Hash limit was not satisfied.",
	))

	lastRewarded := abi.ReadUniversal("lastRewarded", from, "0")
	abi.PreciseMul(lastRewarded, REWARD_PER_SECOND, "0")
	method.AddExecution(abi.WriteUniversal("lastRewarded", ZERO_ADDRESS, util.UTime()))

	return method
}
