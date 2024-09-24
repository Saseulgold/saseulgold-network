package core

/* Operator is extended struct and method for ABI */


func OpWriteLocal(a, b Ia) *ABI {
	abi := &ABI{ name: "WriteLocal", items: []Ia{a, b}, Res: nil}
	return abi
}

func OpWriteUniv(a, b Ia) *ABI {
	abi := &ABI{ name: "WriteUniv", items: []Ia{a, b}, Res: nil}
	return abi
}