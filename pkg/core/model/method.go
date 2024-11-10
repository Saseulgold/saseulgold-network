package model

import (
	"encoding/json"
	F "hello/pkg/util"
)

type Parameters map[string]interface{}
type Execution interface{}
type Excutions []Execution

type ParamValues map[string]interface{}

type Method struct {
	methodType string
	machine    string
	name       string
	version    string
	space      string
	writer     string
	parameters Parameters
	executions []Execution
}

func NewMethod(initialInfo map[string]interface{}) *Method {
	m := &Method{
		parameters: make(Parameters),
		executions: make([]Execution, 0),
	}

	if t, ok := initialInfo["t"]; ok {
		m.methodType = t.(string)
	} else if t, ok := initialInfo["type"]; ok {
		m.methodType = t.(string)
	} else {
		m.methodType = "request"
	}

	if v, ok := initialInfo["m"]; ok {
		m.machine = v.(string)
	} else if v, ok := initialInfo["machine"]; ok {
		m.machine = v.(string)
	} else {
		m.machine = "0.2.0"
	}

	if v, ok := initialInfo["n"]; ok {
		m.name = v.(string)
	} else if v, ok := initialInfo["name"]; ok {
		m.name = v.(string)
	} else {
		m.name = ""
	}

	if v, ok := initialInfo["v"]; ok {
		m.version = v.(string)
	} else if v, ok := initialInfo["version"]; ok {
		m.version = v.(string)
	} else {
		m.version = "1"
	}

	if v, ok := initialInfo["s"]; ok {
		m.space = v.(string)
	} else if v, ok := initialInfo["space"]; ok {
		m.space = v.(string)
	} else {
		m.space = ""
	}

	if v, ok := initialInfo["w"]; ok {
		m.writer = v.(string)
	} else if v, ok := initialInfo["writer"]; ok {
		m.writer = v.(string)
	} else {
		m.writer = ""
	}

	if v, ok := initialInfo["p"]; ok {
		m.parameters = v.(Parameters)
	} else if v, ok := initialInfo["parameters"]; ok {
		m.parameters = v.(Parameters)
	}

	if v, ok := initialInfo["e"]; ok {
		m.executions = v.([]Execution)
	} else if v, ok := initialInfo["executions"]; ok {
		m.executions = v.([]Execution)
	}

	return m
}

func (m *Method) Compile() map[string]interface{} {
	return map[string]interface{}{
		"t": m.methodType,
		"m": m.machine,
		"n": m.name,
		"v": m.version,
		"s": m.space,
		"w": m.writer,
		"p": m.parameters,
		"e": m.executions,
	}
}

func (m *Method) JSON() string {
	data, _ := json.Marshal(m.Compile())
	return string(data)
}

func (m *Method) CID() string {
	return F.SpaceID(m.GetWriter(), m.GetSpace())
}

func (m *Method) GetType() string {
	return m.methodType
}

func (m *Method) SetType(methodType string) {
	if methodType != "" {
		m.methodType = methodType
	}
}

func (m *Method) GetMachine() string {
	return m.machine
}

func (m *Method) SetMachine(machine string) {
	if machine != "" {
		m.machine = machine
	}
}

func (m *Method) GetName() string {
	return m.name
}

func (m *Method) SetName(name string) {
	if name != "" {
		m.name = name
	}
}

func (m *Method) GetVersion() string {
	return m.version
}

func (m *Method) SetVersion(version string) {
	if version != "" {
		m.version = version
	}
}

func (m *Method) GetSpace() string {
	return m.space
}

func (m *Method) SetSpace(space string) {
	if space != "" {
		m.space = space
	}
}

func (m *Method) GetWriter() string {
	return m.writer
}

func (m *Method) SetWriter(writer string) {
	if writer != "" {
		m.writer = writer
	}
}

func (m *Method) GetParameters() Parameters {
	return m.parameters
}

func (m *Method) SetParameters(parameters Parameters) {
	if parameters != nil {
		m.parameters = parameters
	}
}

func (m *Method) AddParameter(parameter Parameter) {
	if parameter.ObjValidity() && m.parameters[parameter.GetName()] == nil {
		m.parameters[parameter.GetName()] = parameter.Obj()
	}
}

func (m *Method) GetExecutions() []Execution {
	return m.executions
}

func (m *Method) SetExecutions(executions []Execution) {
	if executions != nil {
		m.executions = executions
	}
}

func (m *Method) AddExecution(execution interface{}) {
	m.executions = append(m.executions, execution)
}
