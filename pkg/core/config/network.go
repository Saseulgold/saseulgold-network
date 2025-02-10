package config

import (
	"log"

	"gopkg.in/ini.v1"
)

var (
	VERSION                  string
	SWIFT_PORT               string
	SWIFT_HOST               string
	CLI_DEFAULT_REQUEST      string
	NETWORK_DIFF             string
	SG_HARDFORK_START_HEIGHT int64
	INITIAL_SUPPLY           string
	IS_REPLICA               bool
	IS_TESTNET               bool
	BIN_EXEC_ALIAS           string
)

func init() {
	cfg, err := ini.Load("network.ini")
	if err != nil {
		log.Fatalf("Failed to load network.ini: %v", err)
	}

	network := cfg.Section("network")
	VERSION = network.Key("version").MustString("1.0.0")
	SWIFT_PORT = network.Key("swift_port").MustString("9001")
	SWIFT_HOST = network.Key("swift_host").MustString("localhost")
	CLI_DEFAULT_REQUEST = network.Key("cli_default_request").MustString("localhost:9001")
	NETWORK_DIFF = network.Key("network_diff").MustString("1000")
	SG_HARDFORK_START_HEIGHT = network.Key("sg_hardfork_start_height").MustInt64(0)
	INITIAL_SUPPLY = network.Key("initial_supply").MustString("1000000000000000000")
	IS_REPLICA = network.Key("is_replica").MustBool(false)
	IS_TESTNET = network.Key("is_testnet").MustBool(false)
	BIN_EXEC_ALIAS = network.Key("bin_exec_alias").MustString("sg")
}

const SEND_FEE = "5000000000000"
const TRNF_FEE = "5000000000000"
const MINT_FEE = "5000000000000"
const SWAP_DEDUCT_RATE = "0.97"
const SWAP_FEE_RATE = "0.03"
const REWARD_DIFFICULTY = "2000"
const REWARD_PER_SECOND = "19422000000000000000"
const MULTIPLIER = "1000000000000000000"
const PUBLISH_FEE_PER_BYTE = "1000000000000000000"
const HASH_COUNT = "115792089237316195423570985008687907853269984665640564039457584007913129639936"
