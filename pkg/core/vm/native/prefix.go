package native

import (
	. "hello/pkg/core/config"
	. "hello/pkg/util"
)

func ContractPrefix() string {
	return StatusPrefix(ZERO_ADDRESS, RootSpace(), "contract")
}

func RequestPrefix() string {
	return StatusPrefix(ZERO_ADDRESS, RootSpace(), "request")
}

func BalancePrefix() string {
	return StatusPrefix(ZERO_ADDRESS, RootSpace(), "balance")
}
