package core

const STATE_NULL = "null"
const STATE_READ = "read"
const STATE_CONDITION = "condition"
const STATE_EXECUTION = "execution"
const STATE_MAIN = "main"
const STATE_POST = "post"

func RaiseTypeError(msg string) {
	emsg := "Exception: TypeError - " + msg
	panic(emsg)
}

type Interpreter struct {
		mode					string
		signed_data 	Ia

		rg_exception 	*ABI
    rg_state 			string
		rg_break			bool
		rg_process		string

		result 				*Ia
		post_process 	[]*ABI

		executions 		[]*ABI
		contract			*Contract

		paramValues 	ParamValueMap
}

var _instance *Interpreter

func Instance() *Interpreter {
    if(_instance == nil) {
        _instance = &Interpreter{ rg_state: "init" };
				_instance.reset()
    }
    return _instance
}

func (this *Interpreter) reset() {
	/** TODO CHECK WHEN THIS SHOULD BE EXECUTED **/
	this.signed_data = nil

	this.rg_exception = nil
	this.rg_state = "init"
	this.rg_break = false 
	this.rg_process = ""

	this.result = nil
	this.post_process = []*ABI{}

	this.executions = []*ABI{}
}

func (this *Interpreter) SetMode(mode string) {
	this.mode = mode
}

func (this *Interpreter) setState(state string) {
	// READ OR ... //
	this.rg_state = state
}

func (this *Interpreter) setProcessState(process string) {
	// MAIN OR POST //
	this.rg_process = process
}

func (this *Interpreter) GetMode() string {
	return this.mode
}

func (this *Interpreter) EvalExecution(entry *ABI) {
	nitems := []Ia{}

	for _, item := range entry.items {
		nitems = append(nitems, item)
	}

	entry.SetItems(nitems)
	entry.Eval(this)
}

func (this *Interpreter) setParamValues(paramValues ParamValueMap) {
	this.paramValues = paramValues
}

func (this *Interpreter) GetParamValue(key string) Ia {
	return this.paramValues[key]
}

func (this *Interpreter) ExecuteContract(contract *Contract, paramValues ParamValueMap) {
	this.contract = contract
	this.executions = contract.GetExecutions()
	this.setParamValues(paramValues)

	this.Read()
	this.reset()
}

func (this *Interpreter) Read() {
	this.setState(STATE_READ)
	this.setProcessState(STATE_MAIN)

	for _, ex := range this.executions {
		this.EvalExecution(ex)
	}

	this.setProcessState(STATE_POST)

	for _, ex := range this.post_process {
		this.EvalExecution(ex)
	}
}

func (this *Interpreter) GetProcessState() string {
	return this.rg_process
}
