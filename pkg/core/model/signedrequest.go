package model

import (
	"fmt"
	S "hello/pkg/core/structure"
	"hello/pkg/crypto"
	. "hello/pkg/util"
	F "hello/pkg/util"
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

	requestData, ok := req.Data.Get("request")
	if !ok {
		return req
	}

	if v, ok := requestData.(*S.OrderedMap).Get("cid"); ok {
		req.Cid = v.(string)
	}

	if v, ok := requestData.(*S.OrderedMap).Get("type"); ok {
		req.Type = v.(string)
	}

	if v, ok := requestData.(*S.OrderedMap).Get("timestamp"); ok {
		req.Timestamp = v.(int64)
	} else {
		req.Timestamp = Utime()
	}

	req.Hash = req.GetRequestHash()

	if v, ok := data.Get("public_key"); ok {
		req.Xpub = v.(string)
	}
	if v, ok := data.Get("signature"); ok {
		req.Signature = v.(string)
	}

	return req
}

func (req SignedRequest) Ser() string {
	return req.Data.Ser()
}

func (req SignedRequest) GetSize() int {
	return len(req.Ser())
}

func (req *SignedRequest) GetRequestHash() string {
	request, ok := req.Data.Get("request")
	if !ok || request == nil {
		return ""
	}

	timestamp, ok := request.(*S.OrderedMap).Get("timestamp")
	if !ok {
		return ""
	}

	ser := request.(*S.OrderedMap).Ser()

	timestampInt64, ok := timestamp.(int64)
	if !ok {
		if timestampInt, ok := timestamp.(int); ok {
			timestampInt64 = int64(timestampInt)
		} else {
			return ""
		}
	}

	return TimeHash(Hash(ser), timestampInt64)
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
		return fmt.Errorf("the signed request must contain the xpub parameter")
	}

	request, ok := req.Data.Get("request")
	if !ok || request == nil {
		return fmt.Errorf("the signed request must contain the transaction parameter")
	}

	txType, ok := request.(*S.OrderedMap).Get("type")
	if !ok || txType == nil {
		return fmt.Errorf("the signed request must contain the request.type parameter")
	}

	if _, ok := txType.(string); !ok {
		return fmt.Errorf("Parameter request.type must be of string type")
	}

	timestamp, ok := request.(*S.OrderedMap).Get("timestamp")
	if !ok || timestamp == nil {
		return fmt.Errorf("the signed request must contain the request.timestamp parameter")
	}

	if _, ok := timestamp.(int64); !ok {
		if _, ok := timestamp.(int); !ok {
			return fmt.Errorf("parameter request.timestamp must be of integer type")
		}
	}

	if req.Signature == "" {
		return fmt.Errorf("the signed request must contain the signature parameter")
	}

	// Verify signature
	hash := req.GetRequestHash()

	if !crypto.SignatureValidity(hash, req.Xpub, req.Signature) {
		return fmt.Errorf("invalid signature: %s (request hash: %s)", req.Signature, hash)
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
	request, ok := req.Data.Get("request")
	if !ok || request == nil {
		return nil
	}

	return NewSignedDataFromRequestData(&req)
}

func (req SignedRequest) GetRequestXpub() string {
	return req.Xpub
}

func (req SignedRequest) GetRequestSignature() string {
	return req.Signature
}

func (req SignedRequest) GetRequestCID() string {
	if req.Cid == "" {
		return F.RootSpaceId()
	}
	return req.Cid
}
