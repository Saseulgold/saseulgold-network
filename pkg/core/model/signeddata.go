package model

import (
	"encoding/json"
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

func NewSignedDataFromTransaction(data *SignedTransaction) *SignedData {
	hash, err := data.GetTxHash()
	if err != nil {
		return nil
	}
	return &SignedData{
		Data:            data.Data,
		PublicKey:       data.GetXpub(),
		Signature:       data.GetSignature(),
		Hash:            hash,
		Cid:             data.GetCid(),
		Type:            data.GetType(),
		Timestamp:       int64(data.GetTimestamp()),
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
