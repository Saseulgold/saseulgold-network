package core

import (
	// f "hello/pkg/util"
	"encoding/json"
	// "reflect"
)

type Update struct{}

type Compiled map[string]Ia
type ParamMap map[string]Param
type ParamValueMap map[string]Ia

type Contract struct {
	itype       string
	machine     string
	name        string
	version     string
	writer      string
	space       string
	parameters  ParamMap
	paramValues ParamValueMap
	executions  []*ABI
}

func (this *Contract) SetParams(parameters ParamMap) {
	this.parameters = parameters
}

func (this *Contract) AddParameter(name string, itype string) Param {
	param := NewParam(name, itype)
	this.parameters[name] = param
	return param
}

func NewContract() Contract {
	c := Contract{}

	c.parameters = ParamMap{}
	c.paramValues = ParamValueMap{}
	c.executions = []*ABI{}

	return c
}

func (this *Contract) SetMachine(v string) {
	this.machine = v
}

func (this *Contract) SetName(v string) {
	this.name = v
}

func (this *Contract) SetVersion(v string) {
	this.version = v
}

func (this *Contract) SetWriter(v string) {
	this.writer = v
}

func (this *Contract) SetSpace(v string) {
	this.space = v
}

func (this *Contract) AddExecution(abi *ABI) {
	this.executions = append(this.executions, abi)
}

func (this Contract) Compile() Compiled {
	return Compiled{
		"t": this.itype,
		"m": this.machine,
		"n": this.name,
		"v": this.version,
		"s": this.space,
		"w": this.writer,
		"p": this.parameters,
		"e": this.executions,
	}
}

func (this Contract) Json() string {
	j, _ := json.Marshal(this.Compile())
	return string(j)
}

func (this Contract) GetExecutions() []*ABI {
	return this.executions
}
