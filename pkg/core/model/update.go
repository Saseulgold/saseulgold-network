package model

import (
	"encoding/json"
	F "hello/pkg/util"
)

type Update struct {
	Key string `json:"status_key"`
	Old string `json:"old"`
	New string `json:"new"`
}

type UpdateLog struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func NewUpdateLog(old string, new string) UpdateLog {
	return UpdateLog{Old: old, New: new}
}

func (u Update) GetHash() string {
	return F.Concat(u.Key, F.Hash(u.SerUpdateLog()))
}

func (u Update) SerUpdateLog() string {
	j, _ := json.Marshal(NewUpdateLog(u.Old, u.New))
	return string(j)
}
