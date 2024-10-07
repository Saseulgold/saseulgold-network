package main

import (
	"fmt"
	"hello/pkg/core/model"
	"hello/pkg/crypto"
	"hello/pkg/util"
	"testing"
)

func Assert(t *testing.T, a bool, errmsg string) {
	if !a {
		t.Errorf(errmsg)
	}
}

func TestOps23(t *testing.T) {
	s := "Hello RPEMD160!"
	h := util.ShortHash(s)
	Assert(t, h == "c0793d33fa20f601364149d8553c6a3954180ac8", "ShortHash Invalid.")

	h = util.Checksum(s)
	Assert(t, h == "a0d4", "Checksum Invalid")

	hex := "af0201"
	fmt.Println(util.Hex2Bin(hex))
	fmt.Println(hex)
	Assert(t, util.Bin2Hex(util.Hex2Bin(hex)) != hex, "Bin2Hex and Hex2Bin is invalid")
}

func TestOps03(t *testing.T) {
	private_key := "9d37f362487688d6d3c4832c4c8886f1607f122b1b4dbb88a22e71ffb2d7cb17"
	pub, _, _ := crypto.CryptoSignSeedKeypair(util.Hex2Bin(private_key))

	if util.Bin2Hex(pub) != "7586c145d04c7b767e0e09b0271eee1ca90b69a5261812baaf1c7c06a6bb3329" {
		t.Errorf("Public key from private key is invalid")
	}

	cstr := private_key + util.Bin2Hex(pub)
	if cstr != "9d37f362487688d6d3c4832c4c8886f1607f122b1b4dbb88a22e71ffb2d7cb177586c145d04c7b767e0e09b0271eee1ca90b69a5261812baaf1c7c06a6bb3329" {
		t.Errorf("Public key drived is invalid")
	}

	// data := model.TransactinData{"timestamp": 0, "msg": "hiroo"}
	b := util.Hex2Bin(private_key + util.Bin2Hex(pub))

	fmt.Println(b)
}

func TestOps13(t *testing.T) {
	const ts int64 = 1728202895902000
	data := model.TransactinData{"msg": "hiroo", "timestamp": ts}

	private_key := "9d37f362487688d6d3c4832c4c8886f1607f122b1b4dbb88a22e71ffb2d7cb17"

	sd := model.NewSignedTransaction(data)
	fmt.Println(sd.Ser())
	txhash := sd.GetTxHash()

	fmt.Println(txhash)

	if txhash != "0623ca97b5c530ea403f501671560f8614d8d4ee4b24763b01a7a878f3ecdca01e2958071916e9" {
		t.Errorf("txHash invalid for case 1")
	}

	res := sd.Sign(private_key)
	fmt.Println(res)
}

/*
func TestOps2(t *testing.T) {
	fmt.Println(sd.Ser())
	fmt.Println(sd.GetSize())
	fmt.Println(sd.GetHash())
	fmt.Println(util.Hex2Bin("29ef901d1b07ae2abe65a9c9479ae150"))
	d := util.Hex2Bin("29ef901d1b07ae2abe65a9c9479ae150")
	_d := util.Bin2Hex(d)
	fmt.Println(_d)

	seed := crypto.GetRandomKeyPairSeed()
	pub, prv, _ := crypto.CryptoSignSeedKeypair(seed)

	fmt.Println(sd.GetHash())

	fmt.Println(util.Bin2Hex(pub))
	fmt.Println(util.Bin2Hex(prv))
	fmt.Println(util.TimeHash("0", 0))
	fmt.Println(sd.GetHash())
}
*/
