package util

import (
	"fmt"
	"sort"
)

type Ia = interface{}

func Map[T any, V any](ts []T, f func(T) V) []V {
	us := make([]V, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Keys[V any](m map[string]V) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

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

func SortStrings(arr []string) []string {
	sortedArr := make([]string, len(arr))
	copy(sortedArr, arr)
	sort.Strings(sortedArr)
	return sortedArr
}

func SortedValueK[T any](m map[string]T) []T {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make([]T, len(keys))

	for i, k := range keys {
		sorted[i] = m[k]
	}
	return sorted
}
