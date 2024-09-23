package core

import (
	f "hello/pkg/util"
	"fmt"
)

type Clojure 	func(abi *ABI) Ia

type ABI struct {
	name			string
	items 		[]Ia
	cloj			Clojure
	Res				Ia
}

func NewParam(key string, value Ia) Param {
	return Param{ k: key, v: value }
}

func Unwrap(item Ia) Ia {
	switch item.(type) {
	case Param:
		return item.(Param).v
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
		panic("unknown type for unwrap")
	}
	return nil
}

func Eq(a, b Ia) ABI {

	cloj := func(abi *ABI) Ia {
		f.Print("eq: ", abi.items[0], abi.items[1])
		return abi.items[0] == abi.items[1]
	}

	return ABI{ name: "eq", items: []Ia{a, b}, cloj: cloj, Res: nil}
}

func ABIException(msg string) ABI {
	cloj := func(abi *ABI) Ia {
		return msg;
	}

	return ABI{ name: "raise_exception", items: []Ia{msg}, cloj: cloj, Res: nil }
}

func Condition(a ABI, t Ia, f Ia) ABI {

	cloj := func(abi *ABI) Ia {
		fmt.Println("processed: ", a.items)	
		res := Unwrap(abi.items[0]).(bool)

		if(res) {
			return Unwrap(abi.items[1])
		} else {
			return Unwrap(abi.items[2])	
		}
	}

	return ABI{ name: "condition", items: []Ia{a, t, f}, cloj: cloj, Res: nil }
}

func Mul(a, b Ia) ABI {

	cloj := func(abi *ABI) Ia {
		_a := abi.items[0]
		_b := abi.items[1]

		switch _a.(type) {
			case HInteger:
				if __b, ok := _b.(HInteger); ok {
					return _a.(HInteger) * __b
				}

				if __b, ok := _b.(HFloat); ok {
					return _a.(HInteger) * HInteger(__b)
				}

			case HFloat:
				if __b, ok := _b.(HInteger); ok {
					return _a.(HFloat) * HFloat(__b)
				}

				if __b, ok := _b.(HFloat); ok {
					return _a.(HFloat) * __b
				}

			default:
				return nil
		}
		return nil;
	}

	return ABI{ name: "add", items: []Ia{a, b}, cloj: cloj, Res: nil}
}


func Add(a, b Ia) ABI {

	cloj := func(abi *ABI) Ia {
		_a := abi.items[0]
		_b := abi.items[1]

		switch _a.(type) {
			case HInteger:
				if __b, ok := _b.(HInteger); ok {
					return _a.(HInteger) + __b
				}

				if __b, ok := _b.(HFloat); ok {
					return _a.(HInteger) + HInteger(__b)
				}

			case HFloat:
				if __b, ok := _b.(HInteger); ok {
					return _a.(HFloat) + HFloat(__b)
				}

				if __b, ok := _b.(HFloat); ok {
					return _a.(HFloat) + __b
				}

			default:
				return nil
		}
		return nil;
	}

	return ABI{ name: "add", items: []Ia{a, b}, cloj: cloj, Res: nil}
}

func isABI(abi Ia) bool {
	_, ok := abi.(ABI)
	return ok
}

func (abi *ABI) Eval() Ia {
	nitems := []Ia{}
	for _, item := range abi.items {
		var _res Ia = nil

		if _item, ok := item.(ABI); ok {
			_res = _item.Eval() 
		} else if _item, ok := item.(Param); ok {
			_res = Unwrap(_item)
		} else {
			_res = item
		}
		
		nitems = append(nitems, _res)
	}

	abi.items = nitems
	f.Print("processed: ", abi.items, len(abi.items))	
	abi.Res = abi.cloj(abi)
	return abi.Res

}