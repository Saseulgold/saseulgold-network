package core

import (
	f "hello/pkg/util"
	"reflect"
	"fmt"
)

type Clojure 	func(interpreter *Interpreter) Ia

var appLogger = f.GetLogger()

type ABI struct {
	name			string
	items 		[]Ia
	cloj			Clojure
	Res				Ia
}

func Unwrap(item Ia) Ia {
	appLogger.Println("unwrap - item type: ", reflect.TypeOf(item))

	switch item.(type) {
	case ParamValue:
		return item.(ParamValue).v
	case *ParamValue:
		return item.(ParamValue).v
	case *ABI:
		return item.(*ABI).Res
	case ABI:
		return item.(ABI).Res
	case HFloat:
		return item
	case HInteger:
		return item
	case HString:
		return item
	case HBool:
		return item
	default: 
		fmt.Println(item)
		panic("unknown type for unwrap" )
	}
	return nil
}

func OpEq(a, b Ia) *ABI {
	abi := &ABI{ name: "eq", items: []Ia{a, b}, Res: nil}

	cloj := func(interpreter *Interpreter) Ia {
		return abi.items[0] == abi.items[1]
	}

	abi.cloj = cloj
	return abi
}

func ABIException(msg string) *ABI {
	abi := &ABI{ name: "raise_exception", items: []Ia{msg}, Res: nil }

	cloj := func(interpreter *Interpreter) Ia {
		return msg;
	}
	abi.cloj = cloj
	return abi
}

func OpCondition(a *ABI) *ABI {
	abi := &ABI{ name: "condition", items: []Ia{a}, Res: nil }

	cloj := func(interpreter *Interpreter) Ia {
		fmt.Println("processed: ", abi.items)	
		res := Unwrap(abi.items[0]).(bool)
		return res
	}
	abi.cloj = cloj
	return abi
}

func OpMul(a, b Ia) *ABI {
	abi := &ABI{ name: "add", items: []Ia{a, b}, Res: nil}

	cloj := func(interpreter *Interpreter) Ia {
		_a := abi.items[0]
		_b := abi.items[1]

		switch _a.(type) {
			case HInteger:
				if __b, ok := _b.(HInteger); ok {
					return _a.(HInteger) * __b
				}

			case HFloat:
				if __b, ok := _b.(HFloat); ok {
					return _a.(HFloat) * __b
				}

			default:
				break
		}
		RaiseTypeError(fmt.Sprintf("Can't Mul between %s, %s", reflect.TypeOf(_a), reflect.TypeOf(_b)))
		return nil
	}

	abi.cloj = cloj
	return abi
}


func OpAdd(a, b Ia) *ABI {
	abi := &ABI{ name: "add", items: []Ia{a, b}, Res: nil}

	cloj := func(interpreter *Interpreter) Ia {
		_a := abi.items[0]
		_b := abi.items[1]

		switch _a.(type) {
			case HInteger:
				if __b, ok := _b.(HInteger); ok {
					return _a.(HInteger) + __b
				}

			case HFloat:
				if __b, ok := _b.(HFloat); ok {
					return _a.(HFloat) + __b
				}

			default:
				break
		}
		RaiseTypeError(fmt.Sprintf("Can't Add between %s, %s", reflect.TypeOf(_a), reflect.TypeOf(_b)))
		return nil
	}

	abi.cloj = cloj
	return abi
}


func (abi *ABI) Eval(interpreter *Interpreter) Ia {
	nitems := []Ia{}

	if(abi.Res != nil) {
		return abi.Res
	}

	for _, item := range abi.items {
		var _res Ia = nil

		if _item, ok := item.(*ABI); ok {
			_res = _item.Eval(interpreter)
		} else if _item, ok := item.(Param); ok {
			_res = interpreter.GetParamValue(_item.GetKey())
		} else if _item, ok := item.(ParamValue); ok {
			_res = _item.GetValue()
		} else {
			_res = item
		}
		
		nitems = append(nitems, _res)
	}

	abi.SetItems(nitems)

	f.Print("processed: ", abi.items, len(abi.items))	
	res := abi.cloj(interpreter)
	f.Print("Res: ", res)
	fmt.Printf("set Res: %p", abi)	
	abi.SetRes(res)

	return abi.Res
}


func (this *ABI) SetItems(items []Ia) {
	this.items = items
}

func (this *ABI) SetRes(res Ia) {
	this.Res = res
}
