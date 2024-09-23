package core

const STATE_NULL = "null"
const STATE_READ = "read"
const STATE_CONDITION = "condition"
const STATE_EXECUTION = "execution"
const STATE_MAIN = "main"
const STATE_POST = "post"

type Interpreter struct {
		mode					string
		signed_data 	Ia

		codeException	*ABI

    rg_state 		string
		rg_break		bool
		rg_process	string

		result 				*Ia
		post_process 	[]ABI
		code 					Code
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
	this.codeException = nil
	this.result = nil
	this.post_process = []ABI{}
	this.rg_break = false 
	this.signed_data = nil
}

func (this *Interpreter) SetMode(mode string) {
	this.mode = mode
}

func (this *Interpreter) GetMode() string {
	return this.mode
}

func (this *Interpreter) SetCodeException(exception *ABI) {
	this.codeException = exception
}

func (this *Interpreter) Read() {
	this.rg_state = STATE_READ
	this.rg_process = STATE_MAIN

	exs := this.code.GetExecutions()
	
	for _, ex := range exs {
		ex.Eval()
	}

	this.rg_process = STATE_POST
}