package model

import (
	. "hello/pkg/core/model"
	"testing"
)

func TestNewParameter(t *testing.T) {
	// Set up test cases
	testCases := []struct {
		name     string
		input    map[string]interface{}
		expected Parameter
	}{
		{
			name: "Case with all fields",
			input: map[string]interface{}{
				"name":         "username",
				"type":         "string",
				"maxlength":    50,
				"requirements": true,
				"default":      "guest",
			},
			expected: Parameter{
				Name:         "username",
				ParamType:    "string",
				MaxLength:    50,
				Requirements: true,
				DefaultVal:   "guest",
			},
		},
		{
			name:  "Case with empty input",
			input: map[string]interface{}{},
			expected: Parameter{
				ParamType: "any",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewParameter(tc.input)
			if result.Name != tc.expected.Name ||
				result.ParamType != tc.expected.ParamType ||
				result.MaxLength != tc.expected.MaxLength ||
				result.Requirements != tc.expected.Requirements {
				t.Errorf("NewParameter() = %v, expected %v", result, tc.expected)
			}
		})
	}
}

func TestParameterValidations(t *testing.T) {
	p := Parameter{
		Name:         "test",
		ParamType:    "string",
		MaxLength:    5,
		Requirements: true,
	}

	t.Run("Structure validation", func(t *testing.T) {
		// Missing required parameter
		err := p.StructureValidity(nil)
		if err == nil {
			t.Error("Should fail when required parameter is missing")
		}

		// Exceeds max length
		err = p.StructureValidity("123456")
		if err == nil {
			t.Error("Should fail when exceeding max length")
		}

		// Valid case
		err = p.StructureValidity("123")
		if err != nil {
			t.Error("Should succeed with valid input")
		}
	})

	t.Run("Type validation", func(t *testing.T) {
		// Invalid type
		err := p.TypeValidity(123)
		if err == nil {
			t.Error("Should fail with invalid type input")
		}

		// Valid case
		err = p.TypeValidity("test")
		if err != nil {
			t.Error("Should succeed with valid type input")
		}
	})
}

func TestParameterGettersAndSetters(t *testing.T) {
	p := Parameter{}

	t.Run("Name setting/getting", func(t *testing.T) {
		p.SetName("testName")
		if p.GetName() != "testName" {
			t.Error("Name is not set correctly")
		}
	})

	t.Run("Type setting/getting", func(t *testing.T) {
		p.SetType("string")
		if p.GetType() != "string" {
			t.Error("Type is not set correctly")
		}
	})

	t.Run("Max length setting/getting", func(t *testing.T) {
		p.SetMaxLength(10)
		if p.GetMaxLength() != 10 {
			t.Error("Max length is not set correctly")
		}
	})

	t.Run("Requirements setting/getting", func(t *testing.T) {
		p.SetRequirements(true)
		if !p.GetRequirements() {
			t.Error("Requirements is not set correctly")
		}
	})

	t.Run("Default value setting/getting", func(t *testing.T) {
		p.SetDefault("default")
		if p.GetDefault() != "default" {
			t.Error("Default value is not set correctly")
		}

		p.SetDefaultNull()
		if p.GetDefault() != nil {
			t.Error("Default value null setting is not working")
		}
	})
}
