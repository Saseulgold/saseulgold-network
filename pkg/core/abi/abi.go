package abi

import (
	"fmt"
)

func DebugLog(args ...interface{}) {
	// if C.IS_TEST {
	if true {
		fmt.Println(args...)
	}
}

type ABI struct {
	Key   string
	Value interface{}
}

// Basic
func Condition(abi interface{}, errMsg interface{}) ABI {
	if errMsg == "" {
		errMsg = "Conditional error"
	}
	return ABI{
		Key:   "$condition",
		Value: []interface{}{abi, errMsg},
	}
}

func Response(abi interface{}) ABI {
	return ABI{
		Key:   "$response",
		Value: []interface{}{abi},
	}
}

func Weight() ABI {
	return ABI{
		Key:   "$weight",
		Value: []interface{}{},
	}
}

func If(condition, trueVal, falseVal interface{}) ABI {
	return ABI{
		Key:   "$if",
		Value: []interface{}{condition, trueVal, falseVal},
	}
}

func And(vars ...interface{}) ABI {
	return ABI{
		Key:   "$and",
		Value: vars,
	}
}

func Or(vars ...interface{}) ABI {
	return ABI{
		Key:   "$or",
		Value: vars,
	}
}

func Len(obj interface{}) ABI {
	return ABI{
		Key:   "$len",
		Value: []interface{}{obj},
	}
}

func Get(obj, key, defaultVal interface{}) ABI {
	return ABI{
		Key:   "$get",
		Value: []interface{}{obj, key, defaultVal},
	}
}

// Arithmetic
func Add(vars ...interface{}) ABI {
	return ABI{
		Key:   "$add",
		Value: vars,
	}
}

func Sub(vars ...interface{}) ABI {
	return ABI{
		Key:   "$sub",
		Value: vars,
	}
}

func Div(vars ...interface{}) ABI {
	return ABI{
		Key:   "$div",
		Value: vars,
	}
}

func Mul(vars ...interface{}) ABI {
	return ABI{
		Key:   "$mul",
		Value: vars,
	}
}

func PreciseAdd(a1, b interface{}, scale interface{}) ABI {
	return ABI{
		Key:   "$precise_add",
		Value: []interface{}{a1, b, scale},
	}
}

func PreciseSub(a1, b interface{}, scale interface{}) ABI {
	return ABI{
		Key:   "$precise_sub",
		Value: []interface{}{a1, b, scale},
	}
}

func PreciseDiv(a1, b interface{}, scale interface{}) ABI {
	return ABI{
		Key:   "$precise_div",
		Value: []interface{}{a1, b, scale},
	}
}

func PreciseMul(a1, b interface{}, scale interface{}) ABI {
	return ABI{
		Key:   "$precise_mul",
		Value: []interface{}{a1, b, scale},
	}
}

func PreciseSqrt(a interface{}, scale interface{}) ABI {
	return ABI{
		Key:   "$precise_sqrt",
		Value: []interface{}{a, scale},
	}
}

func Scale(value interface{}) ABI {
	return ABI{
		Key:   "$scale",
		Value: []interface{}{value},
	}
}

// Cast
func GetType(obj interface{}) ABI {
	return ABI{
		Key:   "$get_type",
		Value: []interface{}{obj},
	}
}

func IsNumeric(vars interface{}) ABI {
	return ABI{
		Key:   "$is_numeric",
		Value: []interface{}{vars},
	}
}

func IsInt(vars interface{}) ABI {
	return ABI{
		Key:   "$is_int",
		Value: []interface{}{vars},
	}
}

func IsString(vars interface{}) ABI {
	return ABI{
		Key:   "$is_string",
		Value: []interface{}{vars},
	}
}

func IsNull(vars interface{}) ABI {
	return ABI{
		Key:   "$is_null",
		Value: []interface{}{vars},
	}
}

func IsBool(vars interface{}) ABI {
	return ABI{
		Key:   "$is_bool",
		Value: []interface{}{vars},
	}
}

func IsArray(vars interface{}) ABI {
	return ABI{
		Key:   "$is_array",
		Value: []interface{}{vars},
	}
}

func IsDouble(vars interface{}) ABI {
	return ABI{
		Key:   "$is_double",
		Value: []interface{}{vars},
	}
}

// Comparison
func Eq(abi1, abi2 interface{}) ABI {
	return ABI{
		Key:   "$eq",
		Value: []interface{}{abi1, abi2},
	}
}

func Ne(abi1, abi2 interface{}) ABI {
	return ABI{
		Key:   "$ne",
		Value: []interface{}{abi1, abi2},
	}
}

func Gt(abi1, abi2 interface{}) ABI {
	return ABI{
		Key:   "$gt",
		Value: []interface{}{abi1, abi2},
	}
}

func Lt(abi1, abi2 interface{}) ABI {
	return ABI{
		Key:   "$lt",
		Value: []interface{}{abi1, abi2},
	}
}

func Gte(abi1, abi2 interface{}) ABI {
	return ABI{
		Key:   "$gte",
		Value: []interface{}{abi1, abi2},
	}
}

func Lte(abi1, abi2 interface{}) ABI {
	return ABI{
		Key:   "$lte",
		Value: []interface{}{abi1, abi2},
	}
}

func In(target, cases interface{}) ABI {
	return ABI{
		Key:   "$in",
		Value: []interface{}{target, cases},
	}
}

func Param(key interface{}) ABI {
	return ABI{
		Key:   "$load_param",
		Value: []interface{}{key},
	}
}

func ReadUniversal(attr, key, defaultVal interface{}) ABI {
	return ABI{
		Key:   "$read_universal",
		Value: []interface{}{attr, key, defaultVal},
	}
}

func ReadLocal(attr, key, defaultVal interface{}) ABI {
	return ABI{
		Key:   "$read_local",
		Value: []interface{}{attr, key, defaultVal},
	}
}

func WriteUniversalBypass(writer, space, attr, key, value interface{}) ABI {
	return ABI{
		Key:   "$write_universal_bypass",
		Value: []interface{}{writer, space, attr, key, value},
	}
}

func WriteUniversal(attr, key, value interface{}) ABI {
	return ABI{
		Key:   "$write_universal",
		Value: []interface{}{attr, key, value},
	}
}

func WriteLocal(attr, key, value interface{}) ABI {
	return ABI{
		Key:   "$write_local",
		Value: []interface{}{attr, key, value},
	}
}

// Util
func ArrayPush(obj, key, value interface{}) ABI {
	return ABI{
		Key:   "$array_push",
		Value: []interface{}{obj, key, value},
	}
}

func Set(obj, key, value interface{}) ABI {

	return ABI{
		Key:   "$set",
		Value: []interface{}{obj, key, value},
	}

}

func Concat(vars interface{}) ABI {
	return ABI{
		Key:   "$concat",
		Value: []interface{}{vars},
	}
}

func Strlen(target interface{}) ABI {
	return ABI{
		Key:   "$strlen",
		Value: []interface{}{target},
	}
}

func RegMatch(reg, value interface{}) ABI {
	return ABI{
		Key:   "$reg_match",
		Value: []interface{}{reg, value},
	}
}

func EncodeJSON(target interface{}) ABI {
	return ABI{
		Key:   "$encode_json",
		Value: []interface{}{target},
	}
}

func DecodeJSON(target interface{}) ABI {
	return ABI{
		Key:   "$decode_json",
		Value: []interface{}{target},
	}
}

func HashLimit(target interface{}) ABI {
	return ABI{
		Key:   "$hash_limit",
		Value: []interface{}{target},
	}
}

func HashMany(vars ...interface{}) ABI {
	return ABI{
		Key:   "$hash_many",
		Value: vars,
	}
}

func Hash(vars ...interface{}) ABI {
	return ABI{
		Key:   "$hash",
		Value: vars,
	}
}

func ShortHash(target interface{}) ABI {
	return ABI{
		Key:   "$short_hash",
		Value: []interface{}{target},
	}
}

func IDHash(k0 interface{}, k1 interface{}) ABI {
	return ABI{
		Key:   "$id_hash",
		Value: []interface{}{k0, k1},
	}
}

func SignVerify(obj, publicKey, signature interface{}) ABI {
	return ABI{
		Key:   "$sign_verify",
		Value: []interface{}{obj, publicKey, signature},
	}
}

func Check(obj, key interface{}) ABI {
	return ABI{
		Key:   "$check",
		Value: []interface{}{obj, key},
	}
}

func ListBlock(page, count interface{}) ABI {
	return ABI{
		Key:   "$list_block",
		Value: []interface{}{page, count},
	}
}

func ListTransaction(count interface{}) ABI {
	return ABI{
		Key:   "$list_transaction",
		Value: []interface{}{count},
	}
}

func Min(a, b interface{}) ABI {
	return ABI{
		Key:   "$min",
		Value: []interface{}{a, b},
	}
}

func Max(a, b interface{}) ABI {
	return ABI{
		Key:   "$max",
		Value: []interface{}{a, b},
	}
}

func Era(mined interface{}) ABI {
	return ABI{
		Key:   "$era",
		Value: []interface{}{mined},
	}
}

func AsString(target interface{}) ABI {
	return ABI{
		Key:   "$as_string",
		Value: []interface{}{target},
	}
}

func SUtime() ABI {
	return ABI{
		Key:   "$sutime",
		Value: []interface{}{},
	}
}
