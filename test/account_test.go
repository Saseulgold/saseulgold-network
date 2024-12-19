package main

import (
	. "hello/pkg/core/model"
	. "hello/pkg/crypto"
	crypto "hello/pkg/crypto"
	"testing"
)

func createValidAccountTest() struct {
	name       string
	privateKey string
	publicKey  string
	address    string
	wantEmpty  bool
} {
	privateKey, publicKey := crypto.GenerateKeyPair()
	address := crypto.GetAddress(publicKey)

	return struct {
		name       string
		privateKey string
		publicKey  string
		address    string
		wantEmpty  bool
	}{
		name:       "Create account with valid private key",
		privateKey: privateKey,
		publicKey:  publicKey,
		address:    address,
		wantEmpty:  false,
	}
}

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		publicKey  string
		address    string
		wantEmpty  bool
	}{
		createValidAccountTest(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := NewAccount(tt.privateKey)

			if tt.wantEmpty {
				if acc.GetPrivateKey() != "" || acc.GetPublicKey() != "" || acc.GetAddress() != "" {
					t.Errorf("Account should be empty but has values")
				}
			} else {
				if acc.GetPrivateKey() != tt.privateKey {
					t.Errorf("Private key does not match\nExpected: %s\nGot: %s",
						tt.privateKey, acc.GetPrivateKey())
				}
				if acc.GetPublicKey() != tt.publicKey {
					t.Errorf("Public key does not match\nExpected: %s\nGot: %s",
						tt.publicKey, acc.GetPublicKey())
				}
				if acc.GetAddress() != tt.address {
					t.Errorf("Address does not match\nExpected: %s\nGot: %s",
						tt.address, acc.GetAddress())
				}
			}
		})
	}
}

func TestSign(t *testing.T) {
	seed, pub := GenerateKeyPair()
	acc := NewAccount(seed)
	message := "Hello, World!"

	tests := []struct {
		name      string
		message   string
		account   *Account
		wantValid bool
	}{
		{
			name:      "Valid signature",
			message:   message,
			account:   acc,
			wantValid: true,
		},
		{
			name:      "Invalid signature with different message",
			message:   "Different message",
			account:   acc,
			wantValid: false,
		},
		{
			name:      "Invalid signature with empty account",
			message:   message,
			account:   NewAccount(""),
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature := Sign(message, tt.account.GetPrivateKey())
			isValid := SignatureValidity(tt.message, pub, signature)
			if isValid != tt.wantValid {
				t.Errorf("SignatureValidity() = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}
