package vm

import (
	. "hello/pkg/core/abi"
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
	StateMain
	StatePost
)

type Process int

const (
	ProcessNull Process = iota
	ProcessMain
	ProcessPost
)

type MethodFunc func(*Interpreter, interface{}) interface{}

type Interpreter struct {
	mode string

	SignedData  *SignedData
	code        *Method
	postProcess *Method
	breakFlag   bool

	result           interface{}
	weight           int64
	methods          map[string]MethodFunc
	state            State
	process          Process
	universals       map[string]interface{}
	locals           map[string]interface{}
	universalUpdates map[string]map[string]interface{}
	localUpdates     map[string]map[string]interface{}
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		SignedData:       nil,
		code:             nil,
		postProcess:      nil,
		universals:       make(map[string]interface{}),
		locals:           make(map[string]interface{}),
		universalUpdates: make(map[string]map[string]interface{}),
		localUpdates:     make(map[string]map[string]interface{}),
	}
}

func (i *Interpreter) Reset() {
	i.SignedData = nil
	i.code = nil
	i.postProcess = nil
	i.breakFlag = false
	i.result = ""
	i.weight = 0

	i.universalUpdates = make(map[string]map[string]interface{})
	i.localUpdates = make(map[string]map[string]interface{})
}

func (i *Interpreter) SetSignedData(signedData *SignedData) {
	i.SignedData = signedData
}

func (i *Interpreter) Init(mode string) {
	println("Init mode:", mode)

	if mode == "" {
		mode = "transaction"
	}

	if i.mode != mode {
		i.mode = mode
		i.methods = make(map[string]MethodFunc)

		i.loadMethod("BasicOperator")
		i.loadMethod("ComparisonOperator")
		i.loadMethod("ArithmeticOperator")
		i.loadMethod("CastOperator")
		i.loadMethod("UtilOperator")
		i.loadMethod("ReadOperator")
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
	for methodName, method := range OperatorFunctions[name] {
		i.methods[methodName] = method
	}
}

func (i *Interpreter) Set(data *SignedData, code *Method, postProcess *Method) {
	i.SignedData = data
	i.code = code
	i.postProcess = postProcess
	i.breakFlag = false
	i.weight = 0
	i.result = "Conditional Error"
	i.setDefaultValue()
}

func (i *Interpreter) Process(abi interface{}) interface{} {

	switch op := abi.(type) {
	case ABI:
		method, ok := i.methods[op.Key]
		if !ok {
			panic("Method not found: " + op.Key)
		}
		DebugLog("Process method:", i.mode, "method:", op.Key, "value:", op.Value)

		if arr, ok := op.Value.([]interface{}); ok {
			for index, v := range arr {
				if _, isABI := v.(ABI); isABI {
					arr[index] = i.Process(v)
				}
			}
			op.Value = arr
		}
		return method(i, op.Value)
	default:
		return op
	}
}

func (i *Interpreter) setDefaultValue() {
	if i.SignedData.GetAttribute("version") == nil {
		i.SignedData.SetAttribute("version", C.VERSION)
	}

	// Set parameter defaults
	for name, param := range i.code.GetParameters() {
		requirements := param.GetRequirements()
		if !requirements && i.SignedData.GetAttribute(name) == nil {
			defaultVal := param.GetDefault()
			i.SignedData.SetAttribute(name, defaultVal)
		}
	}

	// Contract specific defaults
	if i.mode == "transaction" {
		if i.SignedData.GetAttribute("from") == nil {
			i.SignedData.SetAttribute("from", GetAddress(i.SignedData.PublicKey))
		}

		if i.SignedData.GetAttribute("hash") == nil {
			i.SignedData.SetAttribute("hash", i.SignedData.Hash)
		}

		if i.SignedData.GetAttribute("size") == nil {
			i.SignedData.SetAttribute("size", i.SignedData.Size())
		}

		i.weight += i.SignedData.GetInt64("size")
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

func (i *Interpreter) GetResult() interface{} {
	return i.result
}
