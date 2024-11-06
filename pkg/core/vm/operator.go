package core

import (
	. "hello/pkg/core"
)

func OpWriteLocal(a, b Ia) *ABI {
	abi := &ABI{name: "WriteLocal", items: []Ia{a, b}, Res: nil}
	return abi
}

func OpWriteUniv(a, b Ia) *ABI {
	abi := &ABI{name: "WriteUniv", items: []Ia{a, b}, Res: nil}
	return abi
}
