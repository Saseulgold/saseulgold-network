package config

import (
	"hello/pkg/util"
)

const SYSTEM_NONCE = "Fiat lux. "

func ZeroAddress() string {
	return "00000000000000000000000000000000000000000000"
}

func RootSpace() string {
	return util.Hash(SYSTEM_NONCE)
}
