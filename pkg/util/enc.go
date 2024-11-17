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

func Bin2Hex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// DecBin converts a decimal number to a binary string representation with specified length
// length=1: uint8, length=2: uint16 (big endian), length=3/4: uint32 (big endian), else: uint64 (big endian)
func DecBin(dec int, length int) []byte {
	switch length {
	case 1:
		return []byte{byte(dec)}
	case 2:
		b := make([]byte, 2)
		b[0] = byte(dec >> 8)
		b[1] = byte(dec)
		return b
	case 3, 4:
		b := make([]byte, 4)
		b[0] = byte(dec >> 24)
		b[1] = byte(dec >> 16)
		b[2] = byte(dec >> 8)
		b[3] = byte(dec)
		return b
	default:
		b := make([]byte, 8)
		b[0] = byte(dec >> 56)
		b[1] = byte(dec >> 48)
		b[2] = byte(dec >> 40)
		b[3] = byte(dec >> 32)
		b[4] = byte(dec >> 24)
		b[5] = byte(dec >> 16)
		b[6] = byte(dec >> 8)
		b[7] = byte(dec)
		return b
	}
}

func BinDec(bin []byte) int {
	return HexDec(Bin2Hex(bin))
}

func HexDec(hex string) int {
	if len(hex) < 2 {
		val, _ := strconv.ParseInt(hex, 16, 64)
		return int(val)
	}

	lastDigit, _ := strconv.ParseInt(hex[len(hex)-1:], 16, 64)
	rest := HexDec(hex[:len(hex)-1])

	return 16*rest + int(lastDigit)
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

func SpaceID(owner string, space string) string {
	return Hash(Concat(owner, space))
}

func IsHex(hex string) bool {
	if len(hex) == 0 {
		return false
	}

	// Check if string length is even
	if len(hex)%2 != 0 {
		return false
	}

	// Check if string only contains valid hex characters
	for _, c := range hex {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}

	return true
}
