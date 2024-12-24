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
