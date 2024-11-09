package model

import (
	. "hello/pkg/crypto"
)

type Account struct {
	privateKey string
	publicKey  string
	address    string
}

func NewAccount(privateKey string) *Account {
	acc := &Account{}
	if privateKey != "" && KeyValidity(privateKey) {
		acc.SetPrivateKey(privateKey)
		acc.SetPublicKey(GetXpub(privateKey))
		acc.SetAddress(GetAddress(acc.GetPublicKey()))
	}
	return acc
}

func (a *Account) GetPrivateKey() string {
	return a.privateKey
}

func (a *Account) SetPrivateKey(key string) {
	a.privateKey = key
}

func (a *Account) GetPublicKey() string {
	return a.publicKey
}

func (a *Account) SetPublicKey(key string) {
	a.publicKey = key
}

func (a *Account) GetAddress() string {
	return a.address
}

func (a *Account) SetAddress(addr string) {
	a.address = addr
}

func (a *Account) ToMap() map[string]string {
	return map[string]string{
		"private_key": a.GetPrivateKey(),
		"public_key":  a.GetPublicKey(),
		"address":     a.GetAddress(),
	}
}
