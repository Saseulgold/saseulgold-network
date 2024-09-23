package core

type Payload map[string]interface{}

type Request struct {
	reqType		string
	payload 	Payload
}

