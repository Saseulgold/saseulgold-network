package main

import (
	"fmt"
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
	Assert(t, util.Bin2Hex(util.Hex2Bin(hex)) == hex, "Bin2Hex and Hex2Bin is invalid")

	private_key := "264d444fea3fa0fab2acbd2fe65188781688fb5458077c3bc006238a0634e6da"
	address_expected := "a54b31040d2cb66f44098f9d7b8fb89761d73c13eba5"
	fmt.Println("adddddd")
	fmt.Println(util.IDHash(crypto.GetXpub(private_key)))
	Assert(t, util.IDHash(crypto.GetXpub(private_key)) == address_expected, "IdHash is invalid")
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
