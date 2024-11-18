package model

type Parameter struct {
	name         string      `json:"name"`
	paramType    string      `json:"type"`
	maxLength    int         `json:"maxlength"`
	requirements bool        `json:"requirements"`
	defaultVal   interface{} `json:"default"`
}

func NewParameter(initialInfo map[string]interface{}) Parameter {
	p := Parameter{}

	if name, ok := initialInfo["name"].(string); ok {
		p.name = name
	}

	if paramType, ok := initialInfo["type"].(string); ok {
		p.paramType = paramType
	} else {
		p.paramType = "any"
	}

	if maxLen, ok := initialInfo["maxlength"].(int); ok {
		p.maxLength = maxLen
	}

	if req, ok := initialInfo["requirements"].(bool); ok {
		p.requirements = req
	}

	if def, ok := initialInfo["default"]; ok {
		p.defaultVal = def
	}

	return p
}

func (p *Parameter) GetName() string {
	return p.name
}

func (p *Parameter) SetName(name string) {
	if name != "" {
		p.name = name
	}
}

func (p *Parameter) GetType() string {
	return p.paramType
}

func (p *Parameter) SetType(paramType string) {
	if paramType != "" {
		p.paramType = paramType
	}
}

func (p *Parameter) GetMaxLength() int {
	return p.maxLength
}

func (p *Parameter) SetMaxLength(maxLen int) {
	if maxLen > 0 {
		p.maxLength = maxLen
	}
}

func (p *Parameter) GetRequirements() bool {
	return p.requirements
}

func (p *Parameter) SetRequirements(req bool) {
	p.requirements = req
}

func (p *Parameter) GetDefault() interface{} {
	return p.defaultVal
}

func (p *Parameter) SetDefault(def interface{}) {
	if def != nil {
		p.defaultVal = def
	}
}

func (p *Parameter) SetDefaultNull() {
	p.defaultVal = nil
}

func (p *Parameter) ObjValidity() bool {
	return p.name != "" &&
		p.paramType != "" &&
		p.maxLength >= 0
}

func (p *Parameter) StructureValidity(value interface{}) (bool, string) {
	// Requirements check
	if p.requirements && value == nil {
		return false, "The data must contain the '" + p.name + "' parameter."
	}

	// Set default if nil
	if value == nil {
		value = p.defaultVal
	}

	// MaxLength check
	if str, ok := value.(string); ok {
		if len(str) > p.maxLength {
			return false, "The length of the parameter '" + p.name + "' must be less than " + string(p.maxLength) + " characters."
		}
	}

	return true, ""
}

func (p *Parameter) TypeValidity(value interface{}) (bool, string) {
	if !p.requirements {
		return true, ""
	}

	switch p.paramType {
	case "string":
		if _, ok := value.(string); !ok {
			return false, "Parameter '" + p.name + "' must be of string type."
		}
	case "int":
		if _, ok := value.(int); !ok {
			return false, "Parameter '" + p.name + "' must be of integer type."
		}
	case "float64":
		if _, ok := value.(float64); !ok {
			return false, "Parameter '" + p.name + "' must be of float64 type."
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return false, "Parameter '" + p.name + "' must be of array type."
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return false, "Parameter '" + p.name + "' must be of boolean type."
		}
	}

	return true, ""
}

func (p *Parameter) Obj() map[string]interface{} {
	return map[string]interface{}{
		"name":         p.name,
		"type":         p.paramType,
		"maxlength":    p.maxLength,
		"requirements": p.requirements,
		"default":      p.defaultVal,
	}
}
