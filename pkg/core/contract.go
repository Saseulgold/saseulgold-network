package core

type Update struct {

}

type Contract struct {
	itype				string
	machine			string 	
	name				string
	version 		string	
	writer			string
	space				string
	parameters	[]Param
	executions	[]ABI
	updates			[]Update
}

func (this *Contract) AddUpdate(u Update) {
	this.updates = append(this.updates)
}

func (this Contract) GetUpdates() []Update {
	return this.updates
}
