package rpc

type RawRequest struct {
	Payload string
	Peer    string
}

func CreateRequest(payload string, peer string) *RawRequest {
	req := RawRequest{
		Payload: payload,
		Peer:    peer,
	}

	return &req
}

type TransactionRequest struct {
	Payload string
	Peer    string
}

func CreateTransactionRequest(payload string, peer string) *TransactionRequest {
	req := TransactionRequest{
		Payload: payload,
		Peer:    peer,
	}

	return &req
}
