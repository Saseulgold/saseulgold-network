package native

import (
    abi "hello/pkg/core/abi"
    . "hello/pkg/core/config"
    . "hello/pkg/core/model"
    . "hello/pkg/util"
)

// Genesis creates the Genesis pre-compiled contract.
func Genesis() *PreCompiledContract {
    contract := NewPreCompiledContract(map[string]interface{}{
        "type":    "contract",
        "name":    "Genesis",
        "version": "1",
        "space":   RootSpace(),
        "writer":  ZERO_ADDRESS,
    })

    genesis := abi.ReadLayerL2("genesis", "00", nil)

    contract.AddExecution(abi.Condition(
        abi.Ne(genesis, true),
        "There was already a Genesis.",
    ))

    contract.AddExecution(abi.WriteLayerL1("genesis", "00", true))

    return contract
}

// Register registers a new code in the Register pre-compiled contract.
func Register() *PreCompiledContract {
    contract := NewPreCompiledContract(map[string]interface{}{
        "type":    "contract",
        "name":    "Register",
        "version": "1",
        "space":   RootSpace(),
        "writer":  ZERO_ADDRESS,
    })

    contract.AddParameter(NewParameter(map[string]interface{}{
        "name":         "code",
        "type":         "string",
        "maxlength":    65536,
        "requirements": true,
    }))

    from := abi.Param("from")
    code := abi.Param("code")

    decodedCode := abi.DecodeJSON(code)
    codeType, name, nonce, version, writer := extractCodeDetails(decodedCode)

    codeID := abi.IDHash(name, nonce)
    contractInfo := abi.DecodeJSON(abi.ReadLayerL2("contract", codeID, nil))
    requestInfo := abi.DecodeJSON(abi.ReadLayerL2("request", codeID, nil))

    contractVersion, requestVersion := getCodeVersions(contractInfo, requestInfo)

    isNetworkManager := abi.ReadLayerL2("network_manager", from, nil)
    contract.AddExecution(abi.Condition(
        abi.Eq(isNetworkManager, true),
        "You are not network manager.",
    ))

    validateContractInputs(contract, codeType, name, writer, version, requestVersion, contractVersion)

    update := createUpdateEntry(codeType, codeID, code)
    contract.AddExecution(update)

    return contract
}

// Revoke represents the Revoke pre-compiled contract which revokes permissions.
func Revoke() *PreCompiledContract {
    contract := NewPreCompiledContract(map[string]interface{}{
        "type":    "contract",
        "name":    "Revoke",
        "version": "1",
        "space":   RootSpace(),
        "writer":  ZERO_ADDRESS,
    })

    from := abi.Param("from")
    isNetworkManager := abi.ReadLayerL2("network_manager", from, nil)

    contract.AddExecution(abi.Condition(
        abi.Eq(isNetworkManager, true),
        "You are not network manager.",
    ))

    contract.AddExecution(abi.WriteLayerL2("network_manager", from, false))

    return contract
}

// Send transfers balances between accounts in the Send pre-compiled contract.
func Send() *PreCompiledContract {
    contract := NewPreCompiledContract(map[string]interface{}{
        "type":    "contract",
        "name":    "Send",
        "version": "1",
        "space":   RootSpace(),
        "writer":  ZERO_ADDRESS,
    })

    from := abi.Param("from")
    to := abi.Param("to")
    amount := abi.Param("amount")

    fromBalance := abi.ReadLayerL1("balance", from, "0")
    toBalance := abi.ReadLayerL1("balance", to, "0")

    contract.AddExecution(abi.Condition(
        abi.Ne(from, to),
        "Sender and receiver address must be different.",
    ))

    contract.AddExecution(abi.Condition(
        abi.Gt(fromBalance, amount),
        "Balance is not enough.",
    ))

    contract.AddExecution(abi.WriteLayerL1("balance", from, abi.PreciseSub(fromBalance, amount, 0)))
    contract.AddExecution(abi.WriteLayerL1("balance", to, abi.PreciseAdd(toBalance, amount, 0)))

    return contract
}

// Publish releases new code versions within the Publish pre-compiled contract.
func Publish() *PreCompiledContract {
    contract := NewPreCompiledContract(map[string]interface{}{
        "type":    "contract",
        "name":    "Publish",
        "version": "1",
        "space":   RootSpace(),
        "writer":  ZERO_ADDRESS,
    })

    contract.AddParameter(NewParameter(map[string]interface{}{
        "name":         "code",
        "type":         "string",
        "maxlength":    65536,
        "requirements": true,
    }))

    from := abi.Param("from")
    code := abi.Param("code")
    decodedCode := abi.DecodeJSON(code)

    codeType, name, space, version, writer := extractCodeSpecs(decodedCode)

    codeID := abi.Hash(writer, space, name)
    contractInfo := abi.DecodeJSON(abi.ReadLayerL2("contract", codeID, nil))
    requestInfo := abi.DecodeJSON(abi.ReadLayerL2("request", codeID, nil))
    contractVersion, requestVersion := getCodeVersions(contractInfo, requestInfo)

    contract.AddExecution(abi.Condition(
        abi.Eq(writer, from),
        "Writer must be the same as the from address",
    ))

    validateContractInputs(contract, codeType, name, writer, version, requestVersion, contractVersion)

    contract.AddExecution(abi.If(
        abi.Eq(codeType, "contract"),
        abi.WriteLayerL2("contract", codeID, code),
        abi.If(
            abi.Eq(codeType, "request"),
            abi.WriteLayerL2("request", codeID, code),
            false,
        ),
    ))

    return contract
}

// Helper functions for contract operations

func extractCodeDetails(decodedCode interface{}) (codeType, name, nonce, version, writer interface{}) {
    codeType = abi.Get(decodedCode, "type", nil)
    name = abi.Param("name")
    nonce = abi.Get(decodedCode, "nonce", "")
    version = abi.Get(decodedCode, "version", nil)
    writer = abi.Get(decodedCode, "writer", nil)
    return
}

func getCodeVersions(contractInfo, requestInfo interface{}) (contractVersion, requestVersion string) {
    contractVersion = abi.Get(contractInfo, "version", "0")
    requestVersion = abi.Get(requestInfo, "version", "0")
    return
}

func validateContractInputs(contract *PreCompiledContract, codeType, name, writer, version, reqVersion, conVersion interface{}) {
    contract.AddExecution(abi.Condition(
        abi.Eq(writer, ZERO_ADDRESS),
        "Writer must be zero address",
    ))

    contract.AddExecution(abi.Condition(
        abi.IsString(codeType),
        "Invalid type",
    ))

    contract.AddExecution(abi.Condition(
        abi.In(codeType, []interface{}{"contract", "request"}),
        "Type must be one of the following: contract, request",
    ))

    contract.AddExecution(abi.Condition(
        abi.IsString(name),
        "Invalid name",
    ))

    contract.AddExecution(abi.Condition(
        abi.RegMatch("^[A-Za-z_0-9]+$", name),
        "The name must consist of A-Za-z_0-9.",
    ))

    contract.AddExecution(abi.Condition(
        abi.IsNumeric([]interface{}{version}),
        "Invalid version",
    ))

    versionCheck := abi.If(
        abi.Eq(codeType, "contract"),
        abi.Gt(version, conVersion),
        abi.If(
            abi.Eq(codeType, "request"),
            abi.Gt(version, reqVersion),
            false,
        ),
    )

    contract.AddExecution(abi.Condition(
        versionCheck,
        "Only new versions of code can be registered.",
    ))
}

func createUpdateEntry(codeType, codeID, code interface{}) interface{} {
    return abi.If(
        abi.Eq(codeType, "contract"),
        abi.WriteLayerL2("contract", codeID, code),
        abi.If(
            abi.Eq(codeType, "request"),
            abi.WriteLayerL2("request", codeID, code),
            nil,
        ),
    )
}

func extractCodeSpecs(decodedCode interface{}) (codeType, name, space, version, writer interface{}) {
    codeType = abi.Get(decodedCode, "t", nil)
    name = abi.Get(decodedCode, "n", nil)
    space = abi.Get(decodedCode, "s", nil)
    version = abi.Get(decodedCode, "v", nil)
    writer = abi.Get(decodedCode, "w", nil)
    return
}