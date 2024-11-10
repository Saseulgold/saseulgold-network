package vm

import (
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/crypto"
)

type State int

const (
	StateNull State = iota
	StateRead
	StateCondition
	StateExecution
)

type Process int

const (
	ProcessNull Process = iota
	ProcessMain
	ProcessPost
)

type Interpreter struct {
	mode string

	signedData       *SignedData
	code             *Method
	postProcess      *Method
	breakFlag        bool
	result           string
	weight           int64
	methods          []string
	state            State
	process          Process
	universals       map[string]interface{}
	locals           map[string]interface{}
	universalUpdates map[string]map[string]interface{}
	localUpdates     map[string]map[string]interface{}
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		signedData:       nil,
		code:             nil,
		postProcess:      nil,
		universals:       make(map[string]interface{}),
		locals:           make(map[string]interface{}),
		universalUpdates: make(map[string]map[string]interface{}),
		localUpdates:     make(map[string]map[string]interface{}),
	}
}

func (i *Interpreter) Reset() {
	i.signedData = nil
	i.code = nil
	i.postProcess = nil
	i.breakFlag = false
	i.result = ""
	i.weight = 0

	i.universalUpdates = make(map[string]map[string]interface{})
	i.localUpdates = make(map[string]map[string]interface{})
}

func (i *Interpreter) Init(mode string) {
	if mode == "" {
		mode = "transaction"
	}

	if i.mode != mode {
		i.mode = mode
		i.methods = make([]string, 0)

		i.loadMethod("BasicOperator")
		i.loadMethod("ArithmeticOperator")
		i.loadMethod("ComparisonOperator")
		i.loadMethod("UtilOperator")
		i.loadMethod("CastOperator")
		i.loadMethod("ReadOperator")

		if mode == "transaction" {
			i.loadMethod("WriteOperator")
		} else {
			i.loadMethod("ChainOperator")
		}
	}
}

func (i *Interpreter) loadMethod(name string) {
	// TODO: Implement method loading logic
}

func (i *Interpreter) Set(data *SignedData, code *Method, postProcess *Method) {
	i.signedData = data
	i.code = code
	i.postProcess = postProcess
	i.breakFlag = false
	i.weight = 0
	i.result = "Conditional Error"
	i.setDefaultValue()
}

func (i *Interpreter) Process(abi interface{}) interface{} {
	if abiMap, ok := abi.(map[string]interface{}); ok {
		for key, item := range abiMap {
			if len(key) > 0 {
				prefix := key[0:1]
				vars := i.Process(item)

				if prefix == "$" {
					return nil
				} else {
					abiMap[key] = vars
				}
			}
		}
	}
	return abi
}

func (i *Interpreter) setDefaultValue() {
	if i.signedData.GetAttribute("version") == nil {
		i.signedData.SetAttribute("version", C.VERSION)
	}

	// Set parameter defaults
	for name, param := range i.code.Parameters() {
		requirements := param.GetAttribute("requirements")
		if !requirements && i.signedData.GetAttribute(name) == nil {
			defaultVal := param.Get("default")
			i.signedData.SetAttribute(name, defaultVal)
		}
	}

	// Contract specific defaults
	if i.mode == "transaction" {
		if i.signedData.GetAttribute("from") == nil {
			i.signedData.SetAttribute("from", GetAddress(i.signedData.PublicKey))
		}

		if i.signedData.GetAttribute("hash") == nil {
			i.signedData.SetAttribute("hash", i.signedData.Hash)
		}

		if i.signedData.GetAttribute("size") == nil {
			i.signedData.SetAttribute("size", i.signedData.Size())
		}

		i.weight += i.signedData.GetInt64("size")
	}
}

// Additional methods would follow...
