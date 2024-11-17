package main

import (
	C "hello/pkg/core/config"
	. "hello/pkg/util"
	"testing"
)

func TestStatusHash(t *testing.T) {

	space := RootSpace()
	// spaceId := SpaceID(C.ZERO_ADDRESS, space)

	result := StatusHash(C.ZERO_ADDRESS, space, "rewardTime", C.ZERO_ADDRESS)
	if len(result) != 108 {
		t.Errorf("StatusHash produced incorrect length: got %d, expected 108", len(result))
	}
	t.Logf("StatusHash: %s", result)
}

func TestHash0(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}

	for _, test := range tests {
		result := Hash(test.input)
		if result != test.expected {
			t.Errorf("Hash(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestMerkleRoot(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{}, Hash("")},
		{[]string{"a"}, Hash("a")},
		{[]string{"a", "b"}, Hash(Hash("a") + Hash("b"))},
	}

	for _, test := range tests {
		result := MerkleRoot(test.input)
		if result != test.expected {
			t.Errorf("MerkleRoot(%v) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestHexTime(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{1234567890, "000000499602d2"},
		{0, "00000000000000"},
	}

	for _, test := range tests {
		result := HexTime(test.input)
		if result != test.expected {
			t.Errorf("HexTime(%d) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestHex2Bin(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"ff", []byte{255}},
		{"0000", []byte{0}},
	}

	for _, test := range tests {
		result := Hex2Bin(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("Hex2Bin(%s) length mismatch", test.input)
			continue
		}
		for i := range result {
			if result[i] != test.expected[i] {
				t.Errorf("Hex2Bin(%s)[%d] = %d; expected %d", test.input, i, result[i], test.expected[i])
			}
		}
	}
}

func TestDecBin(t *testing.T) {
	tests := []struct {
		input    int
		length   int
		expected []byte
	}{
		{255, 1, []byte{255}},
		{256, 2, []byte{1, 0}},
		{65535, 4, []byte{0, 0, 255, 255}},
	}

	for _, test := range tests {
		result := DecBin(test.input, test.length)
		if len(result) != len(test.expected) {
			t.Errorf("DecBin(%d, %d) length mismatch", test.input, test.length)
			continue
		}
		for i := range result {
			if result[i] != test.expected[i] {
				t.Errorf("DecBin(%d, %d)[%d] = %d; expected %d",
					test.input, test.length, i, result[i], test.expected[i])
			}
		}
	}
}

func TestShortHash(t *testing.T) {
	result := ShortHash("test")
	if len(result) != 40 { // RIPEMD160 produces 20 bytes = 40 hex chars
		t.Errorf("ShortHash produced incorrect length hash: got %d, expected 40", len(result))
	}
}

func TestChecksum(t *testing.T) {
	result := Checksum("test")
	if len(result) != 4 {
		t.Errorf("Checksum produced incorrect length: got %d, expected 4", len(result))
	}
}

func TestIDHash(t *testing.T) {
	result := IDHash("test")
	if len(result) != 44 { // 40 (ShortHash) + 4 (Checksum)
		t.Errorf("IDHash produced incorrect length: got %d, expected 44", len(result))
	}
}
