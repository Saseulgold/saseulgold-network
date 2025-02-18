package model

import (
	"encoding/json"
	. "hello/pkg/core/abi"
	"hello/pkg/core/structure"
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

func ParseMethod(code string) *Method {
	omap, err := structure.ParseOrderedMap(code)
	if err != nil {
		return nil
	}

	newMethod := &Method{
		Parameters: make(Parameters),
	}

	executions, _ := omap.Get("e")
	writer, _ := omap.Get("w")
	space, _ := omap.Get("s")
	version, _ := omap.Get("v")
	name, _ := omap.Get("n")
	machine, _ := omap.Get("m")
	methodtype, _ := omap.Get("t")
	parameters, _ := omap.Get("p")

	newMethod.SetMachine(machine.(string))
	newMethod.SetName(name.(string))
	newMethod.SetVersion(version.(string))
	newMethod.SetSpace(space.(string))
	newMethod.SetWriter(writer.(string))
	newMethod.SetType(methodtype.(string))

	for _, key := range parameters.(*structure.OrderedMap).Keys() {
		paramMap, _ := parameters.(*structure.OrderedMap)
		param, _ := paramMap.Get(key)

		newMethod.AddParameter(ParseParameter(param.(*structure.OrderedMap)))
	}

	for _, execution := range executions.([]interface{}) {
		executionMap := execution.(*structure.OrderedMap)
		key, _ := executionMap.Get("Key")
		value, _ := executionMap.Get("Value")

		// Convert to ABI
		abiExecution := ABI{
			Key:   key.(string),
			Value: parseABIValue(value),
		}

		newMethod.AddExecution(abiExecution)
	}

	return newMethod
}

// Helper function to parse ABI values recursively
func parseABIValue(value interface{}) []interface{} {
	switch v := value.(type) {
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			// If item is a map, it might be another ABI
			if mapItem, ok := item.(*structure.OrderedMap); ok {
				key, _ := mapItem.Get("Key")
				nestedValue, _ := mapItem.Get("Value")
				result[i] = ABI{
					Key:   key.(string),
					Value: parseABIValue(nestedValue),
				}
			} else {
				// For primitive types (string, number, etc)
				result[i] = item
			}
		}
		return result
	default:
		// Single value should still be wrapped in array
		return []interface{}{value}
	}
}

func ParseParameter(paramMap *structure.OrderedMap) Parameter {
	param := Parameter{
		Name:         "",
		ParamType:    "any", // default value
		MaxLength:    0,
		Requirements: false,
		DefaultVal:   nil,
	}

	name, _ := paramMap.Get("name")
	param.Name = name.(string)

	paramType, _ := paramMap.Get("type")
	param.ParamType = paramType.(string)

	maxLength, _ := paramMap.Get("maxlength")
	param.MaxLength = int(maxLength.(int64))

	requirements, _ := paramMap.Get("requirements")
	param.Requirements = requirements.(bool)

	defaultVal, _ := paramMap.Get("default")
	param.DefaultVal = defaultVal

	return param
}
