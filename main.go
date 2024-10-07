package main

import (
    "fmt"
		u "hello/pkg/util"
		_ "hello/pkg/core"
		// "reflect"
)

func main() {
	fmt.Println(u.Hash("hello"))
	fmt.Println(u.Time())
	fmt.Println(u.HexTime(u.Time()))
	fmt.Println(u.TimeHash("asdfasdf", u.Time()))
}

