package core

/* Operator is extended struct and method for ABI */


func WriteLocal(a, b Ia) *ABI {
	abi := &ABI{ name: "WriteLocal", items: []Ia{a, b}, Res: nil}
	return abi
}

func WriteUniv(a, b Ia) *ABI {
	abi := &ABI{ name: "WriteUniv", items: []Ia{a, b}, Res: nil}
	return abi
}