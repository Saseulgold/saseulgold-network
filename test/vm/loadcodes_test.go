package vm

import (
	"hello/pkg/rpc"
	F "hello/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSystemContracts(t *testing.T) {
	contracts := rpc.NativeContracts()

	expectedMethods := []string{
		"Genesis",
		"Register",
		"Revoke",
		"Faucet",
		"Publish",
		"Send",
	}

	for _, methodName := range expectedMethods {
		method, exists := contracts[F.RootSpaceId()][methodName]
		assert.True(t, exists, methodName+" method must exist")
		assert.NotNil(t, method, methodName+" method must not be nil")
	}

	assert.Equal(t, len(expectedMethods), len(contracts[F.RootSpaceId()]), "Number of contract methods must match")
}
