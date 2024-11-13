package vm

var OperatorFunctions = map[string]map[string]MethodFunc{
	"BasicOperator": {
		"$add": OpAdd,
		/**
		"$sub":       OpSub,
		"$mul":       OpMul,
		"$div":       OpDiv,
		"$eq":        OpEq,
		"$neq":       OpNeq,
		"$gt":        OpGt,
		"$gte":       OpGte,
		"$lt":        OpLt,
		"$lte":       OpLte,
		"$and":       OpAnd,
		"$or":        OpOr,
		"$condition": OpCondition,
		"$response":  OpResponse,
		"$if":        OpIf,
		"$weight":    OpWeight,
		**/
	},
	"UtilOperator": {
		/**
		"$array_push":  OpArrayPush,
		"$concat":      OpConcat,
		"$count":       OpCount,
		"$strlen":      OpStrlen,
		"$encode_json": OpEncodeJson,
		"$decode_json": OpDecodeJson,
		"$hash_limit":  OpHashLimit,
		"$hash_many":   OpHashMany,
		"$hash":        OpHash,
		"$short_hash":  OpShortHash,
		"$id_hash":     OpIdHash,
		"$sign_verify": OpSignVerify,
		**/
	},
	"CastOperator": {
		/**
		"$get_type":   OpGetType,
		"$is_numeric": OpIsNumeric,
		"$is_int":     OpIsInt,
		"$as_string":  OpAsString,
		"$is_string":  OpIsString,
		"$is_null":    OpIsNull,
		"$is_bool":    OpIsBool,
		"$is_array":   OpIsArray,
		"$is_double":  OpIsDouble,
		**/
	},
	"ComparisonOperator": {
		/**
		"$eq":  OpEq,
		"$neq": OpNeq,
		"$gt":  OpGt,
		"$gte": OpGte,
		"$lt":  OpLt,
		"$lte": OpLte,
		**/
	},
	"ReadOperator": {
		"$load_param": OpLoadParam,
	},
}
