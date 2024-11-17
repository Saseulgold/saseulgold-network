package model

import (
	S "hello/pkg/core/structure"
	"hello/pkg/crypto"
	. "hello/pkg/util"
)

type SignedRequest struct {
	Data      *S.OrderedMap `json:"data"`
	PublicKey string        `json:"public_key"`
	Signature string        `json:"signature"`
	Cid       string        `json:"-"`
	Type      string        `json:"-"`
	Timestamp int64         `json:"-"`
	Hash      string        `json:"-"`
}

func NewSignedRequest(data *S.OrderedMap) SignedRequest {
	req := SignedRequest{Data: data}

	// PHP의 null coalescing 연산자(??)를 Go에서 구현
	if v, ok := data.Get("cid"); ok {
		req.Cid = v.(string)
	}
	if v, ok := data.Get("type"); ok {
		req.Type = v.(string)
	}
	if v, ok := data.Get("timestamp"); ok {
		req.Timestamp = v.(int64)
	} else {
		req.Timestamp = Utime()
	}

	req.Hash = req.GetRequestHash()
	return req
}

func (req SignedRequest) Ser() string {
	return req.Data.Ser()
}

func (req SignedRequest) GetSize() int {
	return len(req.Ser())
}

func (req SignedRequest) GetRequestHash() string {
	return TimeHash(Hash(req.Ser()), req.Timestamp)
}

func (req *SignedRequest) Sign(privateKey string) string {
	requestHash := req.GetRequestHash()
	req.Signature = crypto.Sign(requestHash, privateKey)
	return req.Signature
}

func (req SignedRequest) IsValid() (bool, string) {
	if req.Data == nil {
		return false, "The request must contain the \"request\" parameter"
	}

	if req.Type == "" {
		return false, "The request must contain the \"request.type\" parameter"
	}

	return true, ""
}

func (req SignedRequest) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"request":    req.Data,
		"public_key": req.PublicKey,
		"signature":  req.Signature,
	}
}
