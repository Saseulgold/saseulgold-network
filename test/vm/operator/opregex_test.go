package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestOpRegexAndConditions(t *testing.T) {
	// Initialize interpreter
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// Define test method
	methodTest := &Method{
		Parameters: Parameters{
			"email": NewParameter(map[string]interface{}{
				"name":         "email",
				"requirements": true,
			}),
			"user": NewParameter(map[string]interface{}{
				"name":         "user",
				"requirements": true,
			}),
		},
		Executions: []Execution{
			// Email format validation
			abi.RegMatch(abi.Param("email"), "^[^@]+@[^@]+\\.[^@]+$"),
			abi.RegMatch(abi.Get(abi.Param("user"), "phone", nil), "/^010-\\d{4}-\\d{4}$/"),

			// Nested condition test
			abi.If(
				abi.And(
					abi.RegMatch(abi.Get(abi.Param("user"), "phone", nil), "/^010-\\d{4}-\\d{4}$/"),
					abi.Gt(abi.Get(abi.Param("user"), "age", nil), "16"),
				),
				"adult user",
				"minor user",
			),

			// Nested OpGet test
			abi.Get(abi.Get(abi.Param("user"), "address", nil), "city", nil),
		},
	}

	// Set test data
	signedData := NewSignedData()
	signedData.SetAttribute("email", "test@example.com")
	signedData.SetAttribute("user", map[string]interface{}{
		"name":  "John Doe",
		"age":   "20",
		"phone": "010-1234-5678",
		"address": map[string]interface{}{
			"city":   "Seoul",
			"street": "Gangnam-daero",
		},
	})

	// Execute
	interpreter.Reset(true)
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodTest)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if err != nil {
		t.Errorf("Error during execution: %v", err)
	}

	// Verify results
	executions := methodTest.GetExecutions()
	t.Log("here: ", executions[1])

	// Verify email regex validation
	if executions[0] != true {
		t.Errorf("Email regex validation failed. Expected: true, Got: %v", executions[0])
	}

	// Verify condition result
	if executions[2] != "adult user" {
		t.Errorf("Age condition verification failed. Expected: adult user, Got: %v", executions[1])
	}

	// Verify OpGet result
	if executions[3] != "Seoul" {
		t.Errorf("Address information verification failed. Expected: Seoul, Got: %v", executions[2])
	}
}
