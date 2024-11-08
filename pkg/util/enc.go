package util

import (
	"crypto/sha256"
	"encoding/hex"
	_ "fmt"
	"strconv"

	"golang.org/x/crypto/ripemd160"

	C "hello/pkg/core/config"
)

func Hash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func FillHash(hash string) string {
	if len(hash) < C.STATUS_HASH_BYTES {
		return PadRight(hash, "0", C.STATUS_HASH_BYTES)
	}
	return hash
}

func MerkleRoot(data []string) string {
	if len(data) == 0 {
		return Hash("")
	}

	parent := Map(data, func(s string) string {
		return Hash(s)
	})

	for len(parent) > 1 {
		child := []string{}
		for i := 0; i < len(parent); i += 2 {
			if i+1 < len(parent) {
				s := Hash(Concat(parent[i], parent[i+1]))
				child = append(child, s)
			} else {
				child = append(child, parent[i])
			}
		}
		parent = child
	}

	return parent[0]
}

func HexTime(utime int64) string {
	hexTime := strconv.FormatInt(utime, 16)
	return PadLeft(hexTime, "0", C.HEX_TIME_SIZE)
}

func TimeHash(obj string, t int64) string {
	return HexTime(t) + Hash(obj)
}

func Hex2Bin(hexStr string) []byte {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		panic("Hex2Bin failed")
	}
	return decoded
}

func Hex2UInt64(hexStr string) uint64 {
	decoded, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		panic("Hex2UInt64 failed")
	}
	return decoded
}

func Hex2Int64(hexStr string) int64 {
	decoded, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		panic("Hex2Int64 failed")
	}
	return decoded
}

func Bin2Hex(byte []byte) string {
	return hex.EncodeToString(byte)
}

func StringToByte(str string) []byte {
	byteArray := make([]byte, len(str))

	for i := 0; i < len(str); i++ {
		byteArray[i] = str[i]
	}

	return byteArray
}

func ShortHash(message string) string {
	hasher := ripemd160.New()
	hasher.Write([]byte(Hash(message)))
	hash := hasher.Sum(nil)
	return Bin2Hex(hash)
}

func Checksum(message string) string {
	return Hash(Hash(message))[:4]
}

func IDHash(meessage string) string {
	sh := ShortHash(meessage)
	return sh + Checksum(sh)
}

func StatusHash(owner string, space string, attr string, key string) string {
	return StatusPrefix(owner, space, attr) + key
}

func StatusPrefix(owner string, space string, attr string) string {
	return Hash(Concat(owner, space, attr))
}

func RootSpace() string {
	return Hash(C.SYSTEM_NONCE)
}
