package main

import (
    "fmt"
		"hello/pkg/core"
		// "reflect"
)


func execute(fn func() interface{}) {
    result := fn()
    fmt.Println(result)
}

func main() {
    //d := casti(1, "float")
    // execute(func() interface{} { return casti(1, "float") });

		var arg0 float64 = 3
		var param = core.NewParam("aa", 2)

		e := core.Add(param, arg0)
		a := core.Add(2, e)
		d := core.Mul(a, 2)

		c := core.Eq(14, d)
		d = core.Condition(c, 3, 2)
		res := d.Eval()
		fmt.Println("res: ",  res)
		/**
		r := a.Process()
		fmt.Println(r)

		c := core.Eq(1, 1, 3)
		abi_eq := core.Eq(false, c)

		
		response := abi_eq.Process()
		fmt.Println(response)


    // i := core.Instance()
		//  num := i.addf(1,2)
		// kv := core.NewPair(1, 2)

    // v := reflect.ValueOf(i)
		// fmt.Println(v)
    // method := v.MethodByName("Addi")
		// fmt.Println(method.IsValid())
		/*
		fmt.Println("hello")

    args := make([]reflect.Value, 2)

    args[0] = reflect.ValueOf(1)
    args[1] = reflect.ValueOf(2)

    response := method.Call(args)
    fmt.Println(response)
		*/

    // fmt.Println("hello")
    // fmt.Println(v.Interface())

    // fmt.Printf("method: %s", v.Type().method(0).Name)
    /**

    **/
}

