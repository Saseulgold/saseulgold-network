package core

type HInteger     = int
type HFloat       = float64
type HString      = string
type HBool 				= bool

type HelloTypes interface {
	HInteger | HFloat | HString
}

type Ia				 =  interface {}

type Pair	struct {
	k				interface {}
	v				interface {}	
}

func (this *Pair) SetKey(v Ia) {
	this.k = v;
}

func (this *Pair) SetValue(v Ia) {
	this.v = v;
}

func (this Pair) GetKey() Ia {
	return this.k;
}

func (this Pair) GetValue() Ia {
	return this.v;
}

func NewPair(k, v Ia) Pair {
	return Pair{k: k, v: v}
}

