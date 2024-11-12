package vm

import (
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	. "hello/pkg/crypto"
	"reflect"
)

type State int

const (
	StateNull State = iota
	StateRead
	StateCondition
	StateExecution
	StateMain
	StatePost
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
	result           interface{}
	weight           int64
	methods          map[string]interface{}
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
	println("Init mode:", mode)

	if mode == "" {
		mode = "transaction"
	}

	if i.mode != mode {
		i.mode = mode
		i.methods = make(map[string]interface{})

		i.loadMethod("BasicOperator")
		i.loadMethod("ArithmeticOperator")
		i.loadMethod("CastOperator")
		i.loadMethod("UtilOperator")
		/*
			if mode == "transaction" {
				i.loadMethod("WriteOperator")
			} else {
				i.loadMethod("ChainOperator")
			}
		*/
	}
}

func (i *Interpreter) loadMethod(name string) {
	for _, methodName := range OperatorFunctions[name] {
		println("Load method:", methodName)
		i.methods[methodName] = reflect.ValueOf(i).MethodByName(methodName).Interface()
	}
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
	println("Process ABI:", abi)

	if abiArr, ok := abi.([]interface{}); ok {
		for idx, item := range abiArr {
			processed := i.Process(item)
			abiArr[idx] = processed
		}
		return abiArr
	}
	return abi
}

func (i *Interpreter) setDefaultValue() {
	if i.signedData.GetAttribute("version") == nil {
		i.signedData.SetAttribute("version", C.VERSION)
	}

	// Set parameter defaults
	for name, param := range i.code.GetParameters() {
		requirements := param.GetRequirements()
		if !requirements && i.signedData.GetAttribute(name) == nil {
			defaultVal := param.GetDefault()
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

func (i *Interpreter) Execute() (interface{}, bool) {
	executions := i.code.GetExecutions()
	postExecutions := i.postProcess.GetExecutions()

	i.state = StateCondition
	i.process = ProcessMain

	// main, condition
	for key, execution := range executions {
		executions[key] = i.Process(execution)

		if i.breakFlag {
			return i.result, false
		}
	}

	// TODO: Hash(executions)
	processLength := len(executions)

	switch i.mode {
	case "transaction":
		if processLength > C.TX_SIZE_LIMIT {
			return "Too long processing.", false
		}
	default:
		if processLength > C.BLOCK_TX_SIZE_LIMIT {
			return "Too long processing.", false
		}
	}

	// post, condition
	i.process = ProcessPost

	for key, execution := range postExecutions {
		postExecutions[key] = i.Process(execution)

		if i.breakFlag {
			return i.result, false
		}
	}

	// main, execution
	i.state = StateExecution
	i.process = ProcessMain

	for key, execution := range executions {
		executions[key] = i.Process(execution)

		if i.breakFlag {
			return i.result, true
		}
	}

	// post, execution
	i.process = ProcessPost

	for key, execution := range postExecutions {
		postExecutions[key] = i.Process(execution)

		if i.breakFlag {
			return i.result, true
		}
	}

	return "", true
}

func (i *Interpreter) SetCode(code *Method) {
	i.code = code
}

func (i *Interpreter) SetPostProcess(postProcess *Method) {
	i.postProcess = postProcess
}
