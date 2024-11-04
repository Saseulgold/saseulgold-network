package model

import (
	"encoding/json"
	F "hello/pkg/util"
)

type Update struct {
	Key string `json:"status_key"`
	Old Ia     `json:"old"`
	New Ia     `json:"new"`
}

type UpdateLog struct {
	Old Ia `json:"old"`
	New Ia `json:"new"`
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
