package core
/**
protected $type = 'code';
protected $machine;
protected $name;
protected $version;
protected $space;
protected $writer;
protected $parameters = [];
protected $executions = [];
*/


type Param = Pair
type ParamValue = Pair

type Code struct {
	itype				string
	machine			string 	
	writer			string
	space				string
	parameters	[]Param
	executions	[]ABI
}

func (this *Code) AddExecution(a ABI) {
	this.executions = append(this.executions, a)
}


