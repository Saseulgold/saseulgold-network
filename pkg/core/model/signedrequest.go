package model

import (
	"fmt"
	S "hello/pkg/core/structure"
	"hello/pkg/crypto"
	. "hello/pkg/util"
)

type SignedRequest struct {
	Data      *S.OrderedMap `json:"data"`
	Xpub      string        `json:"public_key"`
	Signature string        `json:"signature"`
	Cid       string        `json:"-"`
	Type      string        `json:"-"`
	Timestamp int64         `json:"-"`
	Hash      string        `json:"-"`
}

func NewSignedRequest(data *S.OrderedMap) SignedRequest {
	req := SignedRequest{Data: data}

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
		"public_key": req.Xpub,
		"signature":  req.Signature,
	}
}

func (req SignedRequest) Validate() error {
	if req.Xpub == "" {
		return fmt.Errorf("the signed transaction must contain the xpub parameter")
	}

	request, ok := req.Data.Get("request")
	if !ok || request == nil {
		return fmt.Errorf("the signed transaction must contain the transaction parameter")
	}

	txType, ok := request.(*S.OrderedMap).Get("type")
	if !ok || txType == nil {
		return fmt.Errorf("the signed transaction must contain the transaction.type parameter")
	}

	if _, ok := txType.(string); !ok {
		return fmt.Errorf("Parameter transaction.type must be of string type")
	}

	timestamp, ok := request.(*S.OrderedMap).Get("timestamp")
	if !ok || timestamp == nil {
		return fmt.Errorf("the signed transaction must contain the transaction.timestamp parameter")
	}

	if _, ok := timestamp.(int64); !ok {
		if _, ok := timestamp.(int); !ok {
			return fmt.Errorf("parameter transaction.timestamp must be of integer type")
		}
	}

	if req.Signature == "" {
		return fmt.Errorf("the signed transaction must contain the signature parameter")
	}

	// Verify signature
	hash := req.GetRequestHash()

	if !crypto.SignatureValidity(hash, req.Xpub, req.Signature) {
		return fmt.Errorf("invalid signature: %s (transaction hash: %s)", req.Signature, hash, req.Xpub)
	}

	return nil
}

func (req SignedRequest) GetRequestType() string {
	return req.Type
}

func (req SignedRequest) GetRequestTimestamp() int64 {
	return req.Timestamp
}

func (req SignedRequest) GetRequestData() *SignedData {
	return NewSignedDataFromRequest(&req)
}

func (req SignedRequest) GetRequestXpub() string {
	return req.Xpub
}

func (req SignedRequest) GetRequestSignature() string {
	return req.Signature
}

func (req SignedRequest) GetRequestCID() string {
	return req.Cid
}
