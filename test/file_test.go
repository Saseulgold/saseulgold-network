package main

import (
	"hello/pkg/core/config"
	F "hello/pkg/util"
	"testing"
)

func TestHash(t *testing.T) {
	owner := config.ZeroAddress()
	space := F.RootSpace()
	sender := "a53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4"
	expected := "b3c1ed9ce9df9d2531bb6e2945f044590974408f547f3574d56075e13394770da53ac0f003a3507e0d8fa7fb40ac6fa591f91c7227c4"

	hash := F.StatusHash(owner, space, "balance", sender)

	if hash != expected {
		t.Errorf("Expected %s, got %s", expected, hash)
	}
}
