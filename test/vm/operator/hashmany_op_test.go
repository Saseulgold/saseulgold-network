package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestHashMany(t *testing.T) {
	// Initialize interpreter
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// Test cases with input strings and expected hashes
	testCases := []struct {
		inputs   []string
		expected string
	}{
		{
			[]string{"hello", "world", "test"},
			"6b9889152c357a76aed282caea6d3e8dbc850ab2d1acf88c7762c54458b46669",
		},
	}

	for _, tc := range testCases {
		methodTest := &Method{
			Parameters: Parameters{
				"inputs": NewParameter(map[string]interface{}{
					"name":         "inputs",
					"requirements": true,
				}),
			},
			Executions: []Execution{
				abi.HashMany("hello", "world", "test"),
			},
		}

		t.Run("HashMany_Test", func(t *testing.T) {
			signedData := NewSignedData()
			signedData.SetAttribute("inputs", tc.inputs)

			interpreter.Reset(true)
			interpreter.SetSignedData(signedData)
			interpreter.SetCode(methodTest)
			interpreter.SetPostProcess(post)
			_, err := interpreter.Execute()

			if err != nil {
				t.Errorf("Execution error: %v", err)
			}

			executions := methodTest.GetExecutions()
			result := executions[0].(string)

			// Verify results
			if result != tc.expected {
				t.Errorf("Hash mismatch\nExpected: %s\nGot: %s",
					tc.expected, result)
			}

		})
	}
}
