package util

import (
	"fmt"
)

type Ia = interface{}

func Reduce(items []Ia, initial Ia, fn func(Ia, Ia) interface{}) Ia {
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

var memocache = make(map[string](map[Ia]Ia))

func MemoFn0(lot string, f func(var0 Ia) Ia) func(v0 Ia) Ia {
	return func(var0 Ia) Ia {
		if _, ok := memocache[lot]; !ok {
			memocache[lot] = make(map[Ia]Ia)
		}

		if _, ok := memocache[lot][var0]; !ok {
			memocache[lot][var0] = f(var0)
		}

		return memocache[lot][var0]
	}
}
