package vm

import (
	. "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	. "hello/pkg/core/vm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreter_Process(t *testing.T) {
	// Create and initialize interpreter
	interpreter := NewInterpreter()
	interpreter.Init("transaction")

	// Create test data
	signedData := NewSignedData()
	signedData.SetAttribute("code", `{"name":"test","type":"contract"}`)
	signedData.SetAttribute("value", "hello")
	interpreter.SetSignedData(signedData)

	// Create test method
	method := &Method{
		Parameters: Parameters{
			"code": NewParameter(map[string]interface{}{
				"name":         "code",
				"requirements": true,
				"default":      `{"name":"test","type":"contract"}`,
			}),
			"value": NewParameter(map[string]interface{}{
				"name":         "value",
				"requirements": true,
				"default":      "hello",
			}),
		},
		Executions: []Execution{},
	}

	interpreter.SetCode(method)
	interpreter.SetPostProcess(&Method{})

	// Test cases
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "Test nested load_param, is_string and decode_json",
			input: ABI{
				Key: "$is_string",
				Value: []interface{}{
					ABI{
						Key: "$get",
						Value: []interface{}{
							ABI{
								Key: "$decode_json",
								Value: []interface{}{
									ABI{
										Key:   "$load_param",
										Value: []interface{}{"code"},
									},
								},
							},
							"name",
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Test decode_json with load_param",
			input: ABI{
				Key: "$get",
				Value: []interface{}{
					ABI{
						Key: "$decode_json",
						Value: []interface{}{
							ABI{
								Key:   "$load_param",
								Value: []interface{}{"code"},
							},
						},
					},
					"type",
				},
			},
			expected: "contract",
		},
		{
			name: "Test simple load_param",
			input: ABI{
				Key:   "$load_param",
				Value: []interface{}{"value"},
			},
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interpreter.Process(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
