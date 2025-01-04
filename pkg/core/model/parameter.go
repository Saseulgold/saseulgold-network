package model

import "fmt"

type Parameter struct {
	Name         string      `json:"name"`
	ParamType    string      `json:"type"`
	MaxLength    int         `json:"maxlength"`
	Requirements bool        `json:"requirements"`
	DefaultVal   interface{} `json:"default"`
}

func NewParameter(initialInfo map[string]interface{}) Parameter {
	p := Parameter{}

	if name, ok := initialInfo["name"].(string); ok {
		p.Name = name
	}

	if paramType, ok := initialInfo["type"].(string); ok {
		p.ParamType = paramType
	} else {
		p.ParamType = "any"
	}

	if maxLen, ok := initialInfo["maxlength"].(int); ok {
		p.MaxLength = maxLen
	}

	if req, ok := initialInfo["requirements"].(bool); ok {
		p.Requirements = req
	}

	if def, ok := initialInfo["default"]; ok {
		p.DefaultVal = def
	}

	return p
}

func (p *Parameter) GetName() string {
	return p.Name
}

func (p *Parameter) SetName(name string) {
	if name != "" {
		p.Name = name
	}
}

func (p *Parameter) GetType() string {
	return p.ParamType
}

func (p *Parameter) SetType(paramType string) {
	if paramType != "" {
		p.ParamType = paramType
	}
}

func (p *Parameter) GetMaxLength() int {
	return p.MaxLength
}

func (p *Parameter) SetMaxLength(maxLen int) {
	if maxLen > 0 {
		p.MaxLength = maxLen
	}
}

func (p *Parameter) GetRequirements() bool {
	return p.Requirements
}

func (p *Parameter) SetRequirements(req bool) {
	p.Requirements = req
}

func (p *Parameter) GetDefault() interface{} {
	return p.DefaultVal
}

func (p *Parameter) SetDefault(def interface{}) {
	if def != nil {
		p.DefaultVal = def
	}
}

func (p *Parameter) SetDefaultNull() {
	p.DefaultVal = nil
}

func (p *Parameter) ObjValidity() bool {
	return p.Name != "" &&
		p.ParamType != "" &&
		p.MaxLength >= 0
}

func (p *Parameter) StructureValidity(value interface{}) error {
	// Requirements check
	if p.Requirements && value == nil {
		return fmt.Errorf("The data must contain the '%s' parameter.", p.Name)
	}

	// Set default if nil
	if value == nil {
		value = p.DefaultVal
	}

	// MaxLength check
	if str, ok := value.(string); ok {
		if len(str) > p.MaxLength {
			return fmt.Errorf("The length of the parameter '%s' must be less than %d characters.", p.Name, p.MaxLength)
		}
	}

	return nil
}

func (p *Parameter) TypeValidity(value interface{}) error {
	if !p.Requirements && value == nil {
		return nil
	}

	switch p.ParamType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("Parameter '%s' must be of string type.", p.Name)
		}
	case "int":
		if _, ok := value.(int); !ok {
			return fmt.Errorf("Parameter '%s' must be of integer type.", p.Name)
		}
	case "float64":
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("Parameter '%s' must be of float64 type.", p.Name)
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("Parameter '%s' must be of array type.", p.Name)
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("Parameter '%s' must be of boolean type.", p.Name)
		}
	}

	return nil
}

func (p *Parameter) Obj() map[string]interface{} {
	return map[string]interface{}{
		"name":         p.Name,
		"type":         p.ParamType,
		"maxlength":    p.MaxLength,
		"requirements": p.Requirements,
		"default":      p.DefaultVal,
	}
}
