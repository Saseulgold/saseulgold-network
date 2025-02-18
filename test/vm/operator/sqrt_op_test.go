package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestLargeNumberSqrt(t *testing.T) {
	// Initialize interpreter
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	methodTest := &Method{
		Parameters: Parameters{
			"num": NewParameter(map[string]interface{}{
				"name":         "num",
				"requirements": true,
			}),
		},
		Executions: []Execution{
			abi.PreciseSqrt(abi.Param("num"), "0"),
		},
	}

	// Test cases for large numbers (10^18 ~ 10^32)
	testCases := []struct {
		input    string
		expected string
	}{
		{"36", "6"},
		{"144", "12"},
		{"1440000000000", "1200000"},
		{"1000000000000000000", "1000000000"},
		{"100000000000000000000", "10000000000"},
		{"1000000000000000000000000", "1000000000000"},
		{"10000000000000000000000000000", "100000000000000"},
		{"100000000000000000000000000000000", "10000000000000000"},
	}

	for _, tc := range testCases {
		mtehod := methodTest.Copy()
		signedData := NewSignedData()
		signedData.SetAttribute("num", tc.input)

		interpreter.Reset(true)
		interpreter.SetSignedData(signedData)
		interpreter.SetCode(mtehod)
		interpreter.SetPostProcess(post)
		_, err := interpreter.Execute()

		if err != nil {
			t.Errorf("Execution error for input %s: %v", tc.input, err)
		}

		executions := mtehod.GetExecutions()
		if executions[0] != tc.expected {
			t.Errorf("Large number sqrt test failed.\nInput: %s\nExpected: %s\nGot: %v",
				tc.input, tc.expected, executions[0])
		}
	}
}
