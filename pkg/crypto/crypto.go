package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"hello/pkg/util"

	"golang.org/x/crypto/ripemd160"
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

func Sign(message string, privateKey string) string {
	p0 := util.StringToByte(message)

	xpub := GetXpub(privateKey)
	p1 := util.Hex2Bin(privateKey + xpub)
	return util.Bin2Hex(ed25519.Sign(p1, p0))
}

func Ripemd160(message string) string {
	hasher := ripemd160.New()
	hasher.Write([]byte(message))
	hash := hasher.Sum(nil)
	return util.Bin2Hex(hash)
}
