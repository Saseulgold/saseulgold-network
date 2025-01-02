package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestHashLimit(t *testing.T) {
	// Initialize interpreter
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// Test cases for different difficulty values
	testCases := []struct {
		difficulty string
		expected   string
	}{
		{"4000", "0000455e482248dcc5865508d5c2614021a1e41e14c9a66119959015900e515d61001e44c4a58e"},  // lowest difficulty
		{"8000", "000022af2411246e62c32a846ae130a010d0f20f0a64d3308ccac80ac80728aeb0800f226252c7"},  // medium difficulty
		{"12000", "00000f6a48eb2ca2d68fa11e6864159c79406b94e82ccfa3cccbe7213c74f5a2f91c78810f413c"}, // higher difficulty
		{"16000", "00000c9cc74c0d3f69bb55476caefa68c04bfaee3253358604a6d46102eb549c9d45d6f53b06d4"}, // even higher difficulty
		{"20000", "000008abc904491b98b0caa11ab84c2804343c83c29934cc2332b202b201ca2bac2003c89894b1"}, // highest difficulty
	}

	for _, tc := range testCases {
		methodTest := &Method{
			Parameters: Parameters{
				"difficulty": NewParameter(map[string]interface{}{
					"name":         "difficulty",
					"requirements": true,
				}),
			},
			Executions: []Execution{
				abi.HashLimit(abi.Param("difficulty")),
			},
		}

		t.Run("Difficulty_"+tc.difficulty, func(t *testing.T) {
			signedData := NewSignedData()
			signedData.SetAttribute("difficulty", tc.difficulty)

			interpreter.Reset(true)
			interpreter.SetSignedData(signedData)
			interpreter.SetCode(methodTest)
			interpreter.SetPostProcess(post)
			_, err := interpreter.Execute()

			if err != nil {
				t.Errorf("Execution error for difficulty %s: %v", tc.difficulty, err)
			}

			executions := methodTest.GetExecutions()
			if executions[0] != tc.expected {
				t.Errorf("HashLimit test failed for difficulty %s.\nExpected: %s\nGot: %v",
					tc.difficulty, tc.expected, executions[0])
			}
		})
	}
}
