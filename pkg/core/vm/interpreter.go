package vm

import (
	"fmt"
	. "hello/pkg/core/abi"
	C "hello/pkg/core/config"
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"

	. "hello/pkg/crypto"
	F "hello/pkg/util"

	"go.uber.org/zap"
)

var logger = F.GetLogger()

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

func (p Process) String() string {
	switch p {
	case ProcessNull:
		return "ProcessNull"
	case ProcessMain:
		return "ProcessMain"
	case ProcessPost:
		return "ProcessPost"
	default:
		return fmt.Sprintf("Unknown Process(%d)", p)
	}
}

func (s State) String() string {
	switch s {
	case StateNull:
		return "StateNull"
	case StateRead:
		return "StateRead"
	case StateCondition:
		return "StateCondition"
	case StateExecution:
		return "StateExecution"
	case StateMain:
		return "StateMain"
	case StatePost:
		return "StatePost"
	default:
		return fmt.Sprintf("Unknown State(%d)", s)
	}
}

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
	universalUpdates UpdateMap
	localUpdates     UpdateMap
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		SignedData:       nil,
		code:             nil,
		postProcess:      nil,
		universals:       make(map[string]interface{}),
		locals:           make(map[string]interface{}),
		universalUpdates: &map[string]Update{},
		localUpdates:     &map[string]Update{},
	}
}

func (i *Interpreter) Reset(all bool) {
	logger.Info("reset", zap.Bool("truncate", all))

	i.SignedData = nil
	i.code = nil
	i.postProcess = nil
	i.breakFlag = false
	i.result = ""
	i.weight = 0

	if all {
		i.universals = make(map[string]interface{})
		i.locals = make(map[string]interface{})
		i.state = StateNull
		i.process = ProcessNull
	}

	i.universalUpdates = &map[string]Update{}
	i.localUpdates = &map[string]Update{}
}

func (i *Interpreter) Set(data *SignedData, code *Method, postProcess *Method) {
	i.SignedData = data
	i.code = code
	i.postProcess = postProcess
	i.breakFlag = false
	i.weight = 0
	i.result = ""
	i.setDefaultValue()
}

func (i *Interpreter) SetSignedData(signedData *SignedData) {
	i.SignedData = signedData
}

func (i *Interpreter) Init(mode string) {
	logger.Info("init", zap.String("value", mode))

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

		if mode == "transaction" {
			i.loadMethod("WriteOperator")
		} else {
			i.loadMethod("ChainOperator")
		}
	}
}

func (i *Interpreter) Read() {
	logger.Info("read")
	i.state = StateRead
	i.process = ProcessMain

	// common
	for _, execution := range i.code.GetExecutions() {
		i.Process(execution)
	}

	i.process = ProcessPost

	// post process
	for _, execution := range i.postProcess.GetExecutions() {
		i.Process(execution)
	}
}

func (i *Interpreter) loadMethod(name string) {
	for methodName, method := range OperatorFunctions[name] {
		i.methods[methodName] = method
	}
}

func (i *Interpreter) Process(abi interface{}) interface{} {
	switch op := abi.(type) {
	case ABI:
		method, ok := i.methods[op.Key]
		if !ok {
			panic("Method not found: " + op.Key)
		}

		if arr, ok := op.Value.([]interface{}); ok {
			processedArr := make([]interface{}, len(arr))
			for index, v := range arr {
				if abiVal, isABI := v.(ABI); isABI {
					processedArr[index] = i.Process(abiVal)
				} else {
					processedArr[index] = v
				}
			}
			return method(i, processedArr)
		}

		return method(i, op.Value)
	default:
		return op
	}
}

func (i *Interpreter) ParameterValidate() error {
	if i.mode == "transaction" {
		from := i.SignedData.GetAttribute("from")

		// Validate from address matches signer
		if from != GetAddress(i.SignedData.PublicKey) {
			return fmt.Errorf("Invalid from address: %v", from)
		}
	}

	// Validate all parameters
	for _, param := range i.code.GetParameters() {
		// Convert array parameter to Parameter object if needed

		if !param.ObjValidity() {
			return fmt.Errorf("%s error", i.mode)
		}

		value := i.SignedData.GetAttribute(param.GetName())
		if value == nil {
			value = param.GetDefault()
		}

		if err := param.StructureValidity(value); err != nil {
			return fmt.Errorf("Invalid parameter %s: %s", param.GetName(), err.Error())
		}

		if err := param.TypeValidity(value); err != nil {
			return fmt.Errorf("Invalid parameter %s: %s", param.GetName(), err.Error())
		}
	}

	return nil

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

		/**
		if i.SignedData.GetAttribute("size") == nil {
			i.SignedData.SetAttribute("size", i.SignedData.Size())
		}

		i.weight += i.SignedData.GetInt64("size")
		**/
	}
}

func (i *Interpreter) Execute() (interface{}, error) {

	executions := i.code.GetExecutions()
	postExecutions := i.postProcess.GetExecutions()
	logger.Info("execute", zap.Bool("breakFlag", i.breakFlag), zap.String("executions", fmt.Sprintf("%v", executions)))

	i.state = StateCondition
	i.process = ProcessMain

	// main, condition
	for key, execution := range executions {
		executions[key] = i.Process(execution)

		if i.breakFlag {
			return nil, fmt.Errorf("%v", i.result)
		}
	}

	// TODO: Hash(executions)
	processLength := len(executions)

	switch i.mode {
	case "transaction":
		if processLength > C.TX_SIZE_LIMIT {
			msg := "Too long processing."
			return msg, fmt.Errorf(msg)
		}
	default:
		if processLength > C.BLOCK_TX_SIZE_LIMIT {
			msg := "Too long processing."
			return msg, fmt.Errorf(msg)
		}
	}

	// post, condition
	i.process = ProcessPost

	for key, execution := range postExecutions {
		postExecutions[key] = i.Process(execution)

		if i.breakFlag {
			return postExecutions[key], fmt.Errorf("%v", i.result)
		}
	}

	// main, execution
	i.state = StateExecution
	i.process = ProcessMain

	for key, execution := range executions {
		executions[key] = i.Process(execution)

		if i.breakFlag {
			return executions[key], fmt.Errorf("%v", i.result)
		}
	}

	// post, execution
	i.process = ProcessPost

	for key, execution := range postExecutions {
		postExecutions[key] = i.Process(execution)

		if i.breakFlag {
			return postExecutions[key], fmt.Errorf("%v", i.result)
		}
	}

	return i.result, nil
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

func (i *Interpreter) AddUniversalLoads(statusHash string) {
	statusHash = F.FillHash(statusHash)
	logger.Info("add universal loads", zap.String("hash", statusHash))

	if _, ok := i.universals[statusHash]; !ok {
		i.universals[statusHash] = nil
	}
}

func (i *Interpreter) AddLocalLoads(statusHash string) {
	statusHash = F.FillHash(statusHash)

	if _, ok := i.locals[statusHash]; !ok {
		i.locals[statusHash] = nil
	}
}

func (i *Interpreter) SetLocalLoads(statusHash string, value interface{}) {
	statusHash = F.FillHash(statusHash)
	i.locals[statusHash] = value
}

func (i *Interpreter) SetUniversalLoads(statusHash string, value interface{}) {
	logger.Info("set universal loads", zap.String("hash", statusHash), zap.String("value", fmt.Sprintf("%s", value)))
	statusHash = F.FillHash(statusHash)
	i.universals[statusHash] = value
}

func (i *Interpreter) GetLocalStatus(statusHash string, defaultVal interface{}) interface{} {
	statusHash = F.FillHash(statusHash)
	if val, ok := i.locals[statusHash]; ok {
		return val
	}
	return defaultVal
}

func (i *Interpreter) SetLocalStatus(statusHash string, value interface{}) bool {
	statusHash = F.FillHash(statusHash)

	if _, ok := (*i.localUpdates)[statusHash]; ok {
		oldUpdate := (*i.localUpdates)[statusHash]
		(*i.localUpdates)[statusHash] = Update{
			Old: oldUpdate.Old,
			New: value,
		}
	} else {
		(*i.localUpdates)[statusHash] = Update{
			Old: i.GetLocalStatus(statusHash, nil),
			New: value,
		}
	}

	i.locals[statusHash] = value

	return true
}

func (i *Interpreter) SetUniversalStatus(statusHash string, value interface{}) bool {
	statusHash = F.FillHash(statusHash)

	if _, ok := (*i.universalUpdates)[statusHash]; ok {
		oldUpdate := (*i.universalUpdates)[statusHash]
		(*i.universalUpdates)[statusHash] = Update{
			Old: oldUpdate.Old,
			New: value,
		}
	} else {
		(*i.universalUpdates)[statusHash] = Update{
			Old: i.GetUniversalStatus(statusHash, nil),
			New: value,
		}
	}

	i.universals[statusHash] = value
	logger.Info("Interpreter set univ", zap.String("value", fmt.Sprintf("%v", value)), zap.String("statusHash", statusHash))

	return true
}

func (i *Interpreter) GetUniversals() map[string]interface{} {
	return i.universals
}

func (i *Interpreter) GetUniversalStatus(statusHash string, defaultVal interface{}) interface{} {
	statusHash = F.FillHash(statusHash)

	if val, ok := i.universals[statusHash]; ok {
		if val == nil {
			logger.Info("Interpreter get univ default", zap.String("value", fmt.Sprintf("%v", defaultVal)), zap.String("statusHash", statusHash))
			return defaultVal
		}
		logger.Info("Interpreter get univ ", zap.String("value", fmt.Sprintf("%v", val)), zap.String("statusHash", statusHash))
		return val
	}
	return defaultVal
}

func (i *Interpreter) GetUniversalUpdates() UpdateMap {
	return i.universalUpdates
}

func (i *Interpreter) GetLocalUpdates() UpdateMap {
	return i.localUpdates
}

func (i *Interpreter) LoadUniversalStatus() {
	logger.Info("load univ status", zap.Int("len", len(i.universals)))

	if len(i.universals) == 0 {
		return
	}
	statusFile := storage.GetStatusFileInstance()

	for k := range i.universals {
		value := statusFile.GetUniversalStatus(k)
		i.universals[k] = value
		logger.Info("Loading from status file",
			zap.String("key", k),
			zap.Any("value", value))
	}
}

/**

func (i *Interpreter) LoadLocalStatus() {
	if len(i.locals) > 0 {
		statusFile := S.GetStatusFileInstance()

		for k := range i.locals {
			value := statusFile.GetLocalStatus(k)
			i.locals[k] = value
		}
	}
}
*/
