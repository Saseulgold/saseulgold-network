package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestPreciseDiv(t *testing.T) {
	// Initialize interpreter
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// Test cases with different precisions
	precisionTests := []struct {
		precision string
		cases     []struct {
			numerator   string
			denominator string
			expected    string
		}
	}{
		{
			precision: "0",
			cases: []struct {
				numerator   string
				denominator string
				expected    string
			}{
				{"100", "10", "10"},
				{"1000", "10", "100"},
				{"144", "12", "12"},
				{"1000000000000000000", "1000000000", "1000000000"},
				{"100000000000000000000", "10000000000", "10000000000"},
				{"0", "100", "0"},
				{"999999999999999999", "3", "333333333333333333"},
			},
		},
		{
			precision: "5",
			cases: []struct {
				numerator   string
				denominator string
				expected    string
			}{
				{"10", "3", "3.33333"},
				{"1", "3", "0.33333"},
				{"100", "6", "16.66666"},
				{"1000", "7", "142.85714"},
				{"1", "2", "0.50000"},
			},
		},
		{
			precision: "10",
			cases: []struct {
				numerator   string
				denominator string
				expected    string
			}{
				{"10", "3", "3.3333333333"},
				{"1", "3", "0.3333333333"},
				{"100", "6", "16.6666666666"},
				{"1000", "7", "142.8571428571"},
				{"1", "2", "0.5000000000"},
			},
		},
	}

	for _, precisionTest := range precisionTests {
		methodTest := &Method{
			Parameters: Parameters{
				"numerator": NewParameter(map[string]interface{}{
					"name":         "numerator",
					"requirements": true,
				}),
				"denominator": NewParameter(map[string]interface{}{
					"name":         "denominator",
					"requirements": true,
				}),
			},
			Executions: []Execution{
				abi.PreciseDiv(abi.Param("numerator"), abi.Param("denominator"), precisionTest.precision),
			},
		}

		t.Run("Precision_"+precisionTest.precision, func(t *testing.T) {
			for _, tc := range precisionTest.cases {
				method := methodTest.Copy()
				signedData := NewSignedData()
				signedData.SetAttribute("numerator", tc.numerator)
				signedData.SetAttribute("denominator", tc.denominator)

				interpreter.Reset(true)
				interpreter.SetSignedData(signedData)
				interpreter.SetCode(method)
				interpreter.SetPostProcess(post)
				_, err := interpreter.Execute()

				if err != nil {
					t.Errorf("Execution error for input %s/%s with precision %s: %v",
						tc.numerator, tc.denominator, precisionTest.precision, err)
				}

				executions := method.GetExecutions()
				if executions[0] != tc.expected {
					t.Errorf("Precise division test failed with precision %s.\nNumerator: %s\nDenominator: %s\nExpected: %s\nGot: %v",
						precisionTest.precision, tc.numerator, tc.denominator, tc.expected, executions[0])
				}
			}
		})
	}
}
