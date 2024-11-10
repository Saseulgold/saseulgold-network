package vm

type ABI struct{}

// Machine
func (a *ABI) LegacyCondition(abi interface{}, errMsg string) []interface{} {
	if errMsg == "" {
		errMsg = "Conditional error"
	}
	return []interface{}{abi, errMsg}
}

// Basic
func (a *ABI) Condition(abi interface{}, errMsg string) map[string][]interface{} {
	if errMsg == "" {
		errMsg = "Conditional error"
	}
	return map[string][]interface{}{
		"$condition": {abi, errMsg},
	}
}

func (a *ABI) Response(abi interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$response": {abi},
	}
}

func (a *ABI) Weight() map[string][]interface{} {
	return map[string][]interface{}{
		"$weight": {},
	}
}

func (a *ABI) If(condition, trueVal, falseVal interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$if": {condition, trueVal, falseVal},
	}
}

func (a *ABI) And(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$and": vars,
	}
}

func (a *ABI) Or(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$or": vars,
	}
}

func (a *ABI) Get(abi interface{}, key interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$get": {abi, key},
	}
}

// Arithmetic
func (a *ABI) Add(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$add": vars,
	}
}

func (a *ABI) Sub(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$sub": vars,
	}
}

func (a *ABI) Div(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$div": vars,
	}
}

func (a *ABI) Mul(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$mul": vars,
	}
}

func (a *ABI) PreciseAdd(a1, b interface{}, scale interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$precise_add": {a1, b, scale},
	}
}

func (a *ABI) PreciseSub(a1, b interface{}, scale interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$precise_sub": {a1, b, scale},
	}
}

func (a *ABI) PreciseDiv(a1, b interface{}, scale interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$precise_div": {a1, b, scale},
	}
}

func (a *ABI) PreciseMul(a1, b interface{}, scale interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$precise_mul": {a1, b, scale},
	}
}

func (a *ABI) Scale(value interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$scale": {value},
	}
}

// Cast
func (a *ABI) GetType(obj interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$get_type": {obj},
	}
}

func (a *ABI) IsNumeric(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$is_numeric": vars,
	}
}

func (a *ABI) IsInt(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$is_int": vars,
	}
}

func (a *ABI) IsString(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$is_string": vars,
	}
}

func (a *ABI) IsNull(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$is_null": vars,
	}
}

func (a *ABI) IsBool(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$is_bool": vars,
	}
}

func (a *ABI) IsArray(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$is_array": vars,
	}
}

func (a *ABI) IsDouble(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$is_double": vars,
	}
}

// Comparison
func (a *ABI) Eq(abi1, abi2 interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$eq": {abi1, abi2},
	}
}

func (a *ABI) Ne(abi1, abi2 interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$ne": {abi1, abi2},
	}
}

func (a *ABI) Gt(abi1, abi2 interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$gt": {abi1, abi2},
	}
}

func (a *ABI) Lt(abi1, abi2 interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$lt": {abi1, abi2},
	}
}

func (a *ABI) Gte(abi1, abi2 interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$gte": {abi1, abi2},
	}
}

func (a *ABI) Lte(abi1, abi2 interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$lte": {abi1, abi2},
	}
}

func (a *ABI) In(target, cases interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$in": {target, cases},
	}
}

// I/O
func (a *ABI) Param(vars interface{}) map[string]interface{} {
	switch v := vars.(type) {
	case string:
		return map[string]interface{}{
			"$load_param": []interface{}{v},
		}
	case []interface{}:
		return map[string]interface{}{
			"$load_param": v,
		}
	default:
		return map[string]interface{}{
			"$load_param": vars,
		}
	}
}

func (a *ABI) ReadUniversalBypass(writer, space, attr, key, defaultVal interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$read_universal_bypass": {writer, space, attr, key, defaultVal},
	}
}

func (a *ABI) ReadUniversal(attr, key, defaultVal interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$read_universal": {attr, key, defaultVal},
	}
}

func (a *ABI) ReadLocal(attr, key, defaultVal interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$read_local": {attr, key, defaultVal},
	}
}

func (a *ABI) WriteUniversalBypass(writer, space, attr, key, value interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$write_universal_bypass": {writer, space, attr, key, value},
	}
}

func (a *ABI) WriteUniversal(attr, key, value interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$write_universal": {attr, key, value},
	}
}

func (a *ABI) WriteLocal(attr, key, value interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$write_local": {attr, key, value},
	}
}

// Util
func (a *ABI) ArrayPush(obj, key, value interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$array_push": {obj, key, value},
	}
}

func (a *ABI) Concat(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$concat": vars,
	}
}

func (a *ABI) Strlen(target interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$strlen": {target},
	}
}

func (a *ABI) RegMatch(reg, value interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$reg_match": {reg, value},
	}
}

func (a *ABI) EncodeJSON(target interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$encode_json": {target},
	}
}

func (a *ABI) DecodeJSON(target interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$decode_json": {target},
	}
}

func (a *ABI) HashLimit(target interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$hash_limit": {target},
	}
}

func (a *ABI) HashMany(vars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$hash_many": vars,
	}
}

func (a *ABI) Hash(target interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$hash": {target},
	}
}

func (a *ABI) ShortHash(target interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$short_hash": {target},
	}
}

func (a *ABI) IDHash(target interface{}) map[string][]interface{} {
	return map[string][]interface{}{
		"$id_hash": {target},
	}
}

func (a *ABI) SignVerify(obj interface{}, publicKey string, signature string) map[string][]interface{} {
	return map[string][]interface{}{
		"$sign_verify": {obj, publicKey, signature},
	}
}
