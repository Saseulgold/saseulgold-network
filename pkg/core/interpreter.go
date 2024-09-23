package core

/*
protected $mode;

protected $signed_data;
protected $code;
protected $post_process;
protected $break;
protected $result;
public $weight;
*/

type Interpreter struct {
    status        string
		mode					string
		signed_data 	Ia

		codeException	ABI
		rgst_break		bool
		result 				Ia
		post_process 	ABI
		// code					
}

var _instance *Interpreter

func Instance() *Interpreter {
    if(_instance == nil) {
        _instance = &Interpreter{ status: "init" };
    }
    return _instance
}

func (this *Interpreter) reset() {
	this.codeException = nil
	this.result = nil
	this.post_process = nil
	this.rgst_break = nil
	this.signed_data = nil
}

func (this *Interpreter) SetMode(mode string) {
	this.mode = mode
}

func (this *Interpreter) GetMode() {
	return this.mode
}

func (this *Interpreter) SetCodeException(exception ABI) {
	this.codeException = exception
	return this.codeException
}

