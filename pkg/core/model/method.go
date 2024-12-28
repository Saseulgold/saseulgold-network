package model

import (
	"encoding/json"
	. "hello/pkg/core/abi"
	F "hello/pkg/util"
)

type Parameters map[string]Parameter
type Execution interface{}
type Executions []Execution

type ParamValues map[string]interface{}

type Method struct {
	Type       string
	Machine    string
	Name       string
	Version    string
	Space      string
	Writer     string
	Parameters Parameters
	Executions []Execution
}

func NewMethod(initialInfo map[string]interface{}) *Method {
	m := &Method{
		Parameters: make(Parameters),
		Executions: make([]Execution, 0),
	}

	if t, ok := initialInfo["t"]; ok {
		m.Type = t.(string)
	} else if t, ok := initialInfo["type"]; ok {
		m.Type = t.(string)
	} else {
		m.Type = "request"
	}

	if v, ok := initialInfo["m"]; ok {
		m.Machine = v.(string)
	} else if v, ok := initialInfo["machine"]; ok {
		m.Machine = v.(string)
	} else {
		m.Machine = "0.2.0"
	}

	if v, ok := initialInfo["n"]; ok {
		m.Name = v.(string)
	} else if v, ok := initialInfo["name"]; ok {
		m.Name = v.(string)
	} else {
		m.Name = ""
	}

	if v, ok := initialInfo["v"]; ok {
		m.Version = v.(string)
	} else if v, ok := initialInfo["version"]; ok {
		m.Version = v.(string)
	} else {
		m.Version = "1"
	}

	if v, ok := initialInfo["s"]; ok {
		m.Space = v.(string)
	} else if v, ok := initialInfo["space"]; ok {
		m.Space = v.(string)
	} else {
		m.Space = ""
	}

	if v, ok := initialInfo["w"]; ok {
		m.Writer = v.(string)
	} else if v, ok := initialInfo["writer"]; ok {
		m.Writer = v.(string)
	} else {
		m.Writer = ""
	}

	if v, ok := initialInfo["p"]; ok {
		m.Parameters = v.(Parameters)
	} else if v, ok := initialInfo["parameters"]; ok {
		m.Parameters = v.(Parameters)
	}

	if v, ok := initialInfo["e"]; ok {
		m.Executions = v.([]Execution)
	} else if v, ok := initialInfo["executions"]; ok {
		m.Executions = v.([]Execution)
	}

	return m
}

func (m *Method) Compile() map[string]interface{} {
	parameterMap := make(map[string]interface{})
	for name, param := range m.Parameters {
		parameterMap[name] = param.Obj()
	}

	return map[string]interface{}{
		"t": m.Type,
		"m": m.Machine,
		"n": m.Name,
		"v": m.Version,
		"s": m.Space,
		"w": m.Writer,
		"p": parameterMap,
		"e": m.Executions,
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
	return m.Type
}

func (m *Method) SetType(methodType string) {
	if methodType != "" {
		m.Type = methodType
	}
}

func (m *Method) GetMachine() string {
	return m.Machine
}

func (m *Method) SetMachine(machine string) {
	if machine != "" {
		m.Machine = machine
	}
}

func (m *Method) GetName() string {
	return m.Name
}

func (m *Method) SetName(name string) {
	if name != "" {
		m.Name = name
	}
}

func (m *Method) GetVersion() string {
	return m.Version
}

func (m *Method) SetVersion(version string) {
	if version != "" {
		m.Version = version
	}
}

func (m *Method) GetSpace() string {
	return m.Space
}

func (m *Method) SetSpace(space string) {
	if space != "" {
		m.Space = space
	}
}

func (m *Method) GetWriter() string {
	return m.Writer
}

func (m *Method) SetWriter(writer string) {
	if writer != "" {
		m.Writer = writer
	}
}

func (m *Method) GetParameters() Parameters {
	return m.Parameters
}

func (m *Method) SetParameters(parameters Parameters) {
	if parameters != nil {
		m.Parameters = parameters
	}
}

func (m *Method) AddParameter(parameter Parameter) {
	if parameter.ObjValidity() {
		_, exists := m.Parameters[parameter.GetName()]
		if !exists {
			m.Parameters[parameter.GetName()] = parameter
		}
	}
}

func (m *Method) GetExecutions() []Execution {
	return m.Executions
}

func (m *Method) SetExecutions(executions []Execution) {
	if executions != nil {
		m.Executions = executions
	}
}

func (m *Method) AddExecution(execution ABI) {
	m.Executions = append(m.Executions, execution)
}

func (m *Method) GetCode() string {
	compiled := m.Compile()
	codeBytes, _ := json.Marshal(compiled)
	return string(codeBytes)
}

func (m *Method) Copy() *Method {
	newMethod := &Method{
		Type:       m.Type,
		Machine:    m.Machine,
		Name:       m.Name,
		Version:    m.Version,
		Space:      m.Space,
		Writer:     m.Writer,
		Parameters: make(Parameters),
		Executions: make([]Execution, len(m.Executions)),
	}

	for key, param := range m.Parameters {
		newMethod.Parameters[key] = param
	}

	copy(newMethod.Executions, m.Executions)
	return newMethod
}
