package abi

type ABI map[string][]interface{}

// Basic
func Condition(abi interface{}, errMsg string) ABI {
	if errMsg == "" {
		errMsg = "Conditional error"
	}
	return ABI{
		"$condition": {abi, errMsg},
	}
}

func Response(abi interface{}) ABI {
	return ABI{
		"$response": {abi},
	}
}

func Weight() ABI {
	return ABI{
		"$weight": {},
	}
}

func If(condition, trueVal, falseVal interface{}) ABI {
	return ABI{
		"$if": {condition, trueVal, falseVal},
	}
}

func And(vars interface{}) ABI {
	return ABI{
		"$and": []interface{}{vars},
	}
}

func Or(vars interface{}) ABI {
	return ABI{
		"$or": []interface{}{vars},
	}
}

func Get(abi interface{}, key interface{}) ABI {
	return ABI{
		"$get": {abi, key},
	}
}

// Arithmetic
func Add(vars interface{}) ABI {
	return ABI{
		"$add": []interface{}{vars},
	}
}

func Sub(vars interface{}) ABI {
	return ABI{
		"$sub": []interface{}{vars},
	}
}

func Div(vars interface{}) ABI {
	return ABI{
		"$div": []interface{}{vars},
	}
}

func Mul(vars interface{}) ABI {
	return ABI{
		"$mul": []interface{}{vars},
	}
}

func PreciseAdd(a1, b interface{}, scale interface{}) ABI {
	return ABI{
		"$precise_add": {a1, b, scale},
	}
}

func PreciseSub(a1, b interface{}, scale interface{}) ABI {
	return ABI{
		"$precise_sub": {a1, b, scale},
	}
}

func PreciseDiv(a1, b interface{}, scale interface{}) ABI {
	return ABI{
		"$precise_div": {a1, b, scale},
	}
}

func PreciseMul(a1, b interface{}, scale interface{}) ABI {
	return ABI{
		"$precise_mul": {a1, b, scale},
	}
}

func Scale(value interface{}) ABI {
	return ABI{
		"$scale": {value},
	}
}

// Cast
func GetType(obj interface{}) ABI {
	return ABI{
		"$get_type": {obj},
	}
}

func IsNumeric(vars interface{}) ABI {
	return ABI{
		"$is_numeric": []interface{}{vars},
	}
}

func IsInt(vars interface{}) ABI {
	return ABI{
		"$is_int": []interface{}{vars},
	}
}

func IsString(vars interface{}) ABI {
	return ABI{
		"$is_string": []interface{}{vars},
	}
}

func IsNull(vars interface{}) ABI {
	return ABI{
		"$is_null": []interface{}{vars},
	}
}

func IsBool(vars interface{}) ABI {
	return ABI{
		"$is_bool": []interface{}{vars},
	}
}

func IsArray(vars interface{}) ABI {
	return ABI{
		"$is_array": []interface{}{vars},
	}
}

func IsDouble(vars interface{}) ABI {
	return ABI{
		"$is_double": []interface{}{vars},
	}
}

// Comparison
func Eq(abi1, abi2 interface{}) ABI {
	return ABI{
		"$eq": {abi1, abi2},
	}
}

func Ne(abi1, abi2 interface{}) ABI {
	return ABI{
		"$ne": {abi1, abi2},
	}
}

func Gt(abi1, abi2 interface{}) ABI {
	return ABI{
		"$gt": {abi1, abi2},
	}
}

func Lt(abi1, abi2 interface{}) ABI {
	return ABI{
		"$lt": {abi1, abi2},
	}
}

func Gte(abi1, abi2 interface{}) ABI {
	return ABI{
		"$gte": {abi1, abi2},
	}
}

func Lte(abi1, abi2 interface{}) ABI {
	return ABI{
		"$lte": {abi1, abi2},
	}
}

func In(target, cases interface{}) ABI {
	return ABI{
		"$in": {target, cases},
	}
}

// I/O
func Param(vars interface{}) ABI {
	switch v := vars.(type) {
	case string:
		return ABI{
			"$load_param": []interface{}{v},
		}
	case []interface{}:
		return ABI{
			"$load_param": v,
		}
	default:
		return ABI{
			"$load_param": []interface{}{vars},
		}
	}
}

func ReadUniversalBypass(writer, space, attr, key, defaultVal interface{}) ABI {
	return ABI{
		"$read_universal_bypass": {writer, space, attr, key, defaultVal},
	}
}

func ReadUniversal(attr, key, defaultVal interface{}) ABI {
	return ABI{
		"$read_universal": {attr, key, defaultVal},
	}
}

func ReadLocal(attr, key, defaultVal interface{}) ABI {
	return ABI{
		"$read_local": {attr, key, defaultVal},
	}
}

func WriteUniversalBypass(writer, space, attr, key, value interface{}) ABI {
	return ABI{
		"$write_universal_bypass": {writer, space, attr, key, value},
	}
}

func WriteUniversal(attr, key, value interface{}) ABI {
	return ABI{
		"$write_universal": {attr, key, value},
	}
}

func WriteLocal(attr, key, value interface{}) ABI {
	return ABI{
		"$write_local": {attr, key, value},
	}
}

// Util
func ArrayPush(obj, key, value interface{}) ABI {
	return ABI{
		"$array_push": {obj, key, value},
	}
}

func Concat(vars interface{}) ABI {
	return ABI{
		"$concat": []interface{}{vars},
	}
}

func Strlen(target interface{}) ABI {
	return ABI{
		"$strlen": {target},
	}
}

func RegMatch(reg, value interface{}) ABI {
	return ABI{
		"$reg_match": {reg, value},
	}
}

func EncodeJSON(target interface{}) ABI {
	return ABI{
		"$encode_json": {target},
	}
}

func DecodeJSON(target interface{}) ABI {
	return ABI{
		"$decode_json": {target},
	}
}

func HashLimit(target interface{}) ABI {
	return ABI{
		"$hash_limit": {target},
	}
}

func HashMany(vars interface{}) ABI {
	return ABI{
		"$hash_many": []interface{}{vars},
	}
}

func Hash(target interface{}) ABI {
	return ABI{
		"$hash": {target},
	}
}

func ShortHash(target interface{}) ABI {
	return ABI{
		"$short_hash": {target},
	}
}

func IDHash(target interface{}) ABI {
	return ABI{
		"$id_hash": {target},
	}
}

func SignVerify(obj interface{}, publicKey string, signature string) ABI {
	return ABI{
		"$sign_verify": {obj, publicKey, signature},
	}
}
