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
