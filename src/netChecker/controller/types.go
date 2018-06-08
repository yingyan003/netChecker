package controller

import (
	"encoding/json"
	"common/types"
)

type RecordJson struct {
	Type   interface{}     `json:"type"`
	From   *types.PodInfo  `json:"from"`
	To     json.RawMessage `json:"to"`
	Result interface{}     `json:"result"`
	Time   interface{}     `json:"time"`
	Reason interface{}     `json:"reason"`
}

type RecordSimp struct {
	SrcName interface{} `json:"srcName"`
	DstName interface{} `json:'dstName'`
	Result  interface{} `json:"result"`
	Time    interface{} `json:"time"`
	Reason  interface{} `json:"reason"`
}

type RecordJsons []*RecordJson

type ReportSimp []*RecordSimp
