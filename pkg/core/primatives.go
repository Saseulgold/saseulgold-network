package core

type HInteger = int
type HFloat = float64
type HString = string
type HBool = bool
type HObject = Ia
type HArray = []Ia

const IntegerFlag = "int"
const FloatFlag = "float"
const StringFlag = "string"
const BoolFlag = "bool"
const ObjectFlag = "object"
const ArrayFlag = "array"

type HelloTypes interface {
	HInteger | HFloat | HString
}

func NewHInteger(a Ia) HInteger {
	return a.(HInteger)
}

type Ia = interface{}

type Param struct {
	K string
	V interface{}
}

func (this *Param) SetKey(v string) {
	this.K = v
}

func (this *Param) SetValue(v Ia) {
	this.V = v
}

func (this *Param) GetKey() string {
	return this.K
}

func (this *Param) GetValue() Ia {
	return this.V
}

type ParamValue struct {
	k string
	v interface{}
}

func (this *ParamValue) SetKey(v string) {
	this.k = v
}

func (this *ParamValue) SetValue(v Ia) {
	this.v = v
}

func (this *ParamValue) GetKey() string {
	return this.k
}

func (this *ParamValue) GetValue() Ia {
	return this.v
}

func NewParam(key string, value Ia) Param {
	return Param{K: key, V: value}
}

func NewParamValue(key string, value Ia) ParamValue {
	return ParamValue{k: key, v: value}
}

type Pair struct {
	k interface{}
	v interface{}
}

func (this *Pair) SetKey(v Ia) {
	this.k = v
}

func (this *Pair) SetValue(v Ia) {
	this.v = v
}

func (this Pair) GetKey() Ia {
	return this.k
}

func (this Pair) GetValue() Ia {
	return this.v
}

func NewPair(k, v Ia) Pair {
	return Pair{k: k, v: v}
}
