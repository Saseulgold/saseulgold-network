package config

import (
	"os"
)

var IS_TEST = false

// const QUANTUM_ROOT_DIR = "/Users/louis/qcn"
var QUANTUM_ROOT_DIR = func() string {
	dir := os.Getenv("QUANTUM_ROOT_DIR")
	println("QUANTUM_ROOT_DIR env value:", dir)
	if dir != "" {
		return dir
	}
	return ""
}()

var DATA_ROOT_DIR = func() string {
	return QUANTUM_ROOT_DIR + "/data"
}()
var DATA_ROOT_TEST_DIR = func() string {
	return QUANTUM_ROOT_DIR + "/testdata"
}()

const LEDGER_FILESIZE_LIMIT = 268435456

const HEX_TIME_BYTES = 7
const HASH_BYTES = 32
const TIME_HASH_BYTES = HEX_TIME_BYTES + HASH_BYTES
const STATUS_HASH_BYTES = 64

const HASH_SIZE = HASH_BYTES * 2
const HEX_TIME_SIZE = HEX_TIME_BYTES * 2
const TIME_HASH_SIZE = TIME_HASH_BYTES * 2
const ID_HASH_SIZE = 44

const STATUS_PREFIX_SIZE = 64
const STATUS_KEY_SIZE = 64
const STATUS_HASH_SIZE = 128
const STATUS_KEY_BYTES = STATUS_HASH_BYTES
const STATUS_HEAP_BYTES = STATUS_HASH_BYTES + 10

const DATA_ID_BYTES = 2
const SEEK_BYTES = 4
const LENGTH_BYTES = 4

const CHAIN_KEY_BYTES = TIME_HASH_BYTES
const CHAIN_HEADER_BYTES = 4
const CHAIN_HEIGHT_BYTES = 4
const CHAIN_HEAP_BYTES = TIME_HASH_BYTES + 14
