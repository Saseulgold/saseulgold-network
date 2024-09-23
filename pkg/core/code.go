package core

import (
	"encoding/json"
)

type Param Pair
type ParamValue Pair

type Compiled = map[string]Ia

type ICode interface {
	Compile() Compiled
	Json()		
}

type Code struct {
	itype				string
	machine			string 	
	name				string
	version 		string	
	writer			string
	space				string
	parameters	[]Param
	executions	[]ABI
}

func (this *Code) AddExecution(a ABI) {
	this.executions = append(this.executions, a)
}

func (this Code) Compile() Compiled{
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

func (this Code) Json() string {
	j, _ := json.Marshal(this.Compile())
	return string(j)
}

func (this Code) GetExecutions() []ABI {
	return this.executions
}