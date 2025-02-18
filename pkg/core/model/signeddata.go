package model

import (
	"encoding/json"
	. "hello/pkg/core/debug"
	S "hello/pkg/core/structure"
	"hello/pkg/util"
)

type SignedData struct {
	Data            *S.OrderedMap          `json:"data"`
	PublicKey       string                 `json:"public_key"`
	Signature       string                 `json:"signature"`
	Hash            string                 `json:"hash"`
	Cid             string                 `json:"cid"`
	Type            string                 `json:"type"`
	Timestamp       int64                  `json:"timestamp"`
	Attributes      map[string]interface{} `json:"attributes"`
	CachedUniversal map[string]interface{} `json:"cached_universal"`
	CachedLocal     map[string]interface{} `json:"cached_local"`
}

func NewSignedData() *SignedData {
	return &SignedData{
		Attributes:      make(map[string]interface{}),
		CachedUniversal: make(map[string]interface{}),
		CachedLocal:     make(map[string]interface{}),
	}
}

func NewSignedDataFromRequest(req *SignedRequest) *SignedData {
	hash := req.GetRequestHash()
	reqData, ok := req.Data.Get("request")
	if !ok {
		return nil
	}
	return &SignedData{
		Data:      reqData.(*S.OrderedMap),
		PublicKey: req.GetRequestXpub(),
		Signature: req.GetRequestSignature(),
		Hash:      hash,
		Cid:       req.GetRequestCID(),
		Type:      req.GetRequestType(),
		Timestamp: req.GetRequestTimestamp(),
	}
}

func NewSignedDataFromRequestData(req *SignedRequest) *SignedData {
	hash := req.GetRequestHash()
	reqData, ok := req.Data.Get("request")
	if !ok {
		return nil
	}
	return &SignedData{
		Data:            reqData.(*S.OrderedMap),
		PublicKey:       req.GetRequestXpub(),
		Signature:       req.GetRequestSignature(),
		Hash:            hash,
		Cid:             req.GetRequestCID(),
		Type:            req.GetRequestType(),
		Timestamp:       req.GetRequestTimestamp(),
		Attributes:      make(map[string]interface{}),
		CachedUniversal: make(map[string]interface{}),
		CachedLocal:     make(map[string]interface{}),
	}
}

func NewSignedDataFromTransaction(tx *SignedTransaction) *SignedData {

	hash := tx.GetTxHash()
	txData, ok := tx.Data.Get("transaction")
	if !ok {
		return nil
	}
	return &SignedData{
		Data:            txData.(*S.OrderedMap),
		PublicKey:       tx.GetXpub(),
		Signature:       tx.GetSignature(),
		Hash:            hash,
		Cid:             tx.GetCID(),
		Type:            tx.GetType(),
		Timestamp:       int64(tx.GetTimestamp()),
		Attributes:      make(map[string]interface{}),
		CachedUniversal: make(map[string]interface{}),
		CachedLocal:     make(map[string]interface{}),
	}
}

func (s *SignedData) GetAttribute(key string) interface{} {
	if val, ok := s.Attributes[key]; ok {
		return val
	}
	if s.Data != nil {
		v, _ := s.Data.Get(key)
		return v
	}
	return nil
}

func (s *SignedData) SetAttribute(key string, value interface{}) {
	s.Attributes[key] = value
}

func (s *SignedData) GetCachedUniversal(key string) interface{} {
	return s.CachedUniversal[key]
}

func (s *SignedData) SetCachedUniversal(key string, value interface{}) {
	DebugLog("SetCachedUniversal", "key:", key, "value:", value)
	s.CachedUniversal[key] = value
}

func (s *SignedData) GetCachedLocal(key string) interface{} {
	return s.CachedLocal[key]
}

func (s *SignedData) SetCachedLocal(key string, value interface{}) {
	s.CachedLocal[key] = value
}

func (s *SignedData) GetInt64(key string) int64 {
	if val := s.GetAttribute(key); val != nil {
		if i, ok := val.(int64); ok {
			return i
		}
	}
	return 0
}

func (s *SignedData) Obj() map[string]interface{} {
	return map[string]interface{}{
		"data":       s.Data,
		"public_key": s.PublicKey,
		"signature":  s.Signature,
	}
}

func (s *SignedData) GetHash() string {
	return util.TimeHash(util.Hash(s.Data.Ser()), s.Timestamp)
}

func (s *SignedData) JSON() string {
	data, _ := json.Marshal(s.Obj())
	return string(data)
}

func (s *SignedData) Size() int {
	return len(s.JSON())
}
