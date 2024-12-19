package model

import (
    "fmt"
    S "hello/pkg/core/structure"
    "hello/pkg/crypto"
    . "hello/pkg/util"
)

// EncodedRequest represents an encoded request with necessary cryptographic details.
type EncodedRequest struct {
    Data      *S.OrderedMap `json:"data"`
    Xpub      string        `json:"public_key"`
    Signature string        `json:"signature"`
    Cid       string        `json:"-"`
    Type      string        `json:"-"`
    Timestamp int64         `json:"-"`
    Hash      string        `json:"-"`
}

// NewEncodedRequest creates a new EncodedRequest based on provided data.
func NewEncodedRequest(data *S.OrderedMap) EncodedRequest {
    req := EncodedRequest{Data: data}

    if v, ok := data.Get("cid"); ok {
        req.Cid = v.(string)
    }
    if v, ok := data.Get("type"); ok {
        req.Type = v.(string)
    }
    req.Timestamp = getTimestamp(data)
    req.Hash = req.generateRequestHash()

    return req
}

func getTimestamp(data *S.OrderedMap) int64 {
    if v, ok := data.Get("timestamp"); ok {
        return v.(int64)
    }
    return Utime()
}

// Serialize returns the serialized string representation of the request data.
func (req EncodedRequest) Serialize() string {
    return req.Data.Ser()
}

// GetSize returns the size of the serialized request data.
func (req EncodedRequest) GetSize() int {
    return len(req.Serialize())
}

// generateRequestHash computes the request hash using its serialized data and timestamp.
func (req EncodedRequest) generateRequestHash() string {
    return TimeHash(Hash(req.Serialize()), req.Timestamp)
}

// Sign signs the request using the provided private key and updates the signature.
func (req *EncodedRequest) Sign(privateKey string) string {
    requestHash := req.generateRequestHash()
    req.Signature = crypto.Sign(requestHash, privateKey)
    return req.Signature
}

// IsValid checks if the encoded request contains necessary parameters.
func (req EncodedRequest) IsValid() (bool, string) {
    if req.Data == nil {
        return false, "The request must contain the \"data\" parameter"
    }

    if req.Type == "" {
        return false, "The request must contain the \"type\" parameter"
    }

    return true, ""
}

// ToMap converts the encoded request into a map representation.
func (req EncodedRequest) ToMap() map[string]interface{} {
    return map[string]interface{}{
        "data":       req.Data,
        "public_key": req.Xpub,
        "signature":  req.Signature,
    }
}

// Validate performs a comprehensive validation of the encoded request.
func (req EncodedRequest) Validate() error {
    if req.Xpub == "" {
        return fmt.Errorf("the request must contain the xpub parameter")
    }

    request, err := req.getValidatedRequestData()
    if err != nil {
        return err
    }

    if req.Signature == "" {
        return fmt.Errorf("the request must contain the signature parameter")
    }

    hash := req.generateRequestHash()

    if !crypto.SignatureValidity(hash, req.Xpub, req.Signature) {
        return fmt.Errorf("invalid signature: %s (transaction hash: %s)", req.Signature, hash)
    }

    return nil
}

func (req EncodedRequest) getValidatedRequestData() (*S.OrderedMap, error) {
    request, ok := req.Data.Get("request").(*S.OrderedMap)
    if !ok || request == nil {
        return nil, fmt.Errorf("the request must contain valid \"request\" data")
    }

    if _, ok := request.Get("type").(string); !ok {
        return nil, fmt.Errorf("\"type\" parameter must be of string type")
    }

    if _, ok := request.Get("timestamp").(int64); !ok {
        if _, ok := request.Get("timestamp").(int); !ok {
            return nil, fmt.Errorf("\"timestamp\" parameter must be of integer type")
        }
    }

    return request, nil
}

// Getters for EncodedRequest
func (req EncodedRequest) GetRequestType() string         { return req.Type }
func (req EncodedRequest) GetRequestTimestamp() int64     { return req.Timestamp }
func (req EncodedRequest) GetRequestData() *SignedData    { return NewSignedDataFromRequest(&req) }
func (req EncodedRequest) GetRequestXpub() string         { return req.Xpub }
func (req EncodedRequest) GetRequestSignature() string    { return req.Signature }
func (req EncodedRequest) GetRequestCID() string          { return req.Cid }