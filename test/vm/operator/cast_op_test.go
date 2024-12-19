package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestCastOperators(t *testing.T) {
	// 인터프리터 초기화
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// 테스트 메소드 정의
	methodTest := &Method{
		Parameters: Parameters{
			"numValue": NewParameter(map[string]interface{}{
				"name":         "numValue",
				"requirements": true,
			}),
			"strValue": NewParameter(map[string]interface{}{
				"name":         "strValue",
				"requirements": true,
			}),
			"arrayValue": NewParameter(map[string]interface{}{
				"name":         "arrayValue",
				"requirements": true,
			}),
		},
		Executions: []Execution{
			// 타입 체크 연산자 테스트
			abi.GetType(abi.Param("numValue")),

			// 숫자 타입 검증
			abi.IsNumeric(abi.Param("numValue")),
			abi.IsInt(abi.Param("numValue")),
			abi.IsDouble(abi.Param("numValue")),

			// 문자열 타입 검증
			abi.IsString(abi.Param("strValue")),

			// 배열 타입 검증
			abi.IsArray(abi.Param("arrayValue")),

			// null 체크
			abi.IsNull(abi.Param("nullValue")),
		},
	}

	// 테스트 데이터 설정
	signedData := NewSignedData()
	signedData.SetAttribute("numValue", "42")
	signedData.SetAttribute("strValue", "test string")
	signedData.SetAttribute("arrayValue", []interface{}{1, 2, 3})
	signedData.SetAttribute("nullValue", nil)

	// 실행
	interpreter.Reset()
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodTest)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if err != nil {
		t.Errorf("Error during execution: %v", err)
	}

	executions := methodTest.GetExecutions()

	// GetType validation
	expectedType := "int"
	if executions[0] != expectedType {
		t.Errorf("GetType test failed.\nExpected: %s\nGot: %v",
			expectedType, executions[0])
	}

	// IsNumeric validation
	if executions[1] != true {
		t.Errorf("IsNumeric test failed.\nExpected: true\nGot: %v",
			executions[1])
	}

	// IsInt validation
	if executions[2] != true {
		t.Errorf("IsInt test failed.\nExpected: true\nGot: %v",
			executions[2])
	}

	// IsDouble validation
	if executions[3] != false {
		t.Errorf("IsDouble test failed.\nExpected: false\nGot: %v",
			executions[3])
	}

	// IsString validation
	if executions[4] != true {
		t.Errorf("IsString test failed.\nExpected: true\nGot: %v",
			executions[4])
	}

	// IsArray validation
	if executions[5] != true {
		t.Errorf("IsArray test failed.\nExpected: true\nGot: %v",
			executions[5])
	}

	// IsNull validation
	if executions[6] != true {
		t.Errorf("IsNull test failed.\nExpected: true\nGot: %v",
			executions[6])
	}
}
