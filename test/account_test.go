package main

import (
	. "hello/pkg/core/model"
	"testing"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		publicKey  string
		address    string
		wantEmpty  bool
	}{
		{
			name:       "Create account with valid private key",
			privateKey: "41de21adfea9ba36a29513ff4277adc8ecce3fa22c2c57bb0006243eab48e821",
			publicKey:  "391e87c9ceedb34ecd7f74d4536a33851ce54dbb0c2dfbf1a529816f8ed78afd",
			address:    "f43808a3998233c4336d873880fe4a22fdd7eafdd90e",
			wantEmpty:  false,
		},
		{
			name:       "Create account with empty private key",
			privateKey: "",
			publicKey:  "",
			address:    "",
			wantEmpty:  true,
		},
		{
			name:       "Create account with invalid private key",
			privateKey: "invalid_private_key",
			publicKey:  "",
			address:    "",
			wantEmpty:  true,
		},
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
