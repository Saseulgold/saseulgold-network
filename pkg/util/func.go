package util

import (
	"fmt"
)

type Ia = interface {}

// Reduce 함수 정의: 리스트의 각 요소를 순차적으로 처리
func Reduce(items []Ia, initial Ia, fn func(Ia, Ia) interface {}) Ia {
    result := initial
    for _, item := range items {
        result = fn(result, item)
    }
    return result
}

func AllEqual(args ...Ia) Ia {
	first := args[0]

	return Reduce(args[1:], true, func(acc Ia, current Ia) Ia {
			return acc.(bool) && (first == current)
	})
}

func Print(args ...Ia) {
	fmt.Println(args...)	
}