package model

import (
	"encoding/json"
	F "hello/pkg/util"
)

type Update struct {
	Key string
	Old interface{}
	New interface{}
}

type UpdateLog struct {
	Old interface{} `json:"old"`
	New interface{} `json:"new"`
}

func NewUpdateLog(old string, new string) UpdateLog {
	return UpdateLog{Old: old, New: new}
}

func (u Update) GetHash() string {
	return F.Concat(u.Key, F.Hash(u.SerUpdateLog()))
}

func (u Update) SerUpdateLog() string {
	j, _ := json.Marshal(UpdateLog{Old: u.Old, New: u.New})
	return string(j)
}

func (u Update) Ser() string {
	j, _ := json.Marshal(u)
	return string(j)
}
