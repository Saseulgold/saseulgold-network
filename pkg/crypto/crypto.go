package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"hello/pkg/util"
)

// sodium_crypto_sign_seed_keypair
func CryptoSignSeedKeypair(seed []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, nil, fmt.Errorf("invalid seed size: expected %d bytes, got %d", ed25519.SeedSize, len(seed))
	}

	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return publicKey, privateKey, nil
}

func GetRandomKeyPairSeed() []byte {
	seed := make([]byte, ed25519.SeedSize)
	rand.Read(seed)
	return seed
}

func GetXpub(privateKey string) string {
	pub, _, _ := CryptoSignSeedKeypair(util.Hex2Bin(privateKey))
	return util.Bin2Hex(pub)
}

func GetAddress(publicKey string) string {
	if !KeyValidity(publicKey) {
		return ""
	}
	return util.IDHash(publicKey)
}

func AddressValidity(address string) bool {
	return util.IsHex(address)
}

func Sign(message string, privateKey string) string {
	p0 := util.StringToByte(message)

	xpub := GetXpub(privateKey)
	p1 := util.Hex2Bin(privateKey + xpub)
	return util.Bin2Hex(ed25519.Sign(p1, p0))
}

func SignatureValidity(obj string, publicKey string, signature string) bool {
	if !KeyValidity(publicKey) || len(signature) != SIGNATURE_SIZE {
		return false
	}

	return VerifySignature(util.Hex2Bin(signature), util.StringToByte(obj), util.Hex2Bin(publicKey))
}

func VerifySignature(signature, message, publicKey []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}

const (
	KEY_SIZE       = 64  // Equivalent to SODIUM_CRYPTO_AUTH_BYTES * 2
	SIGNATURE_SIZE = 128 // Equivalent to SODIUM_CRYPTO_SIGN_BYTES * 2
)

func KeyValidity(key string) bool {
	return len(key) == KEY_SIZE && util.IsHex(key)
}
