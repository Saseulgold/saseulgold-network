package rpc

type RawRequest struct {
	Payload string
	Peer    string
}

func CreateRawRequest(payload string, peer string) *RawRequest {
	req := RawRequest{
		Payload: payload,
		Peer:    peer,
	}

	return &req
}
