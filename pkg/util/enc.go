package util

import (
	"crypto/sha256"
	"encoding/hex"
	_ "fmt"
	"strconv"

	"golang.org/x/crypto/ripemd160"
)

const HEX_TIME_BYTES = 7
const HEX_TIME_SIZE = HEX_TIME_BYTES * 2

func Hash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
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
		i := 0

		for i+1 < len(parent) {
			if len(parent) > i+1 {
				s := Hash(Concat(parent[i], parent[i+1]))
				child = append(child, s)
			} else {
				child = append(child, parent[i])
			}
			i += 2
		}

		parent = child // 업데이트된 자식을 부모로 설정
	}

	return parent[0]
}

func HexTime(utime int64) string {
	hexTime := strconv.FormatInt(utime, 16)
	return PadLeft(hexTime, "0", HEX_TIME_SIZE)
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

func StateHash(owner string, space string, attr string, key string) string {
	return StatePrefix(owner, space, attr) + key
}

func StatePrefix(owner string, space string, attr string) string {
	return Hash(Concat(owner, space, attr))
}
