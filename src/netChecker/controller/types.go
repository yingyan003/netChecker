package controller

import (
	"encoding/json"
	"common/types"
)

type Record struct {
	Type      string          `json:"type"`
	From      types.PodInfo  `json:"from"`
	To        json.RawMessage `json:"to"`
	Result    bool            `json:"result"`
	Timestamp string          `json:"timestamp"`
	Reason    string          `json:"reason"`
}

type RecordSimp struct {
	SrcName   string `json:"srcName"`
	DstName   string `json:'dstName'`
	Result    bool   `json:"result"`
	Timestamp string `json:"timestamp"`
	Reason    string `json:"reason"`
}

type RecordNodeToNode struct {
	SrcIp string        `json:"srcIp"`
	To    []NodeResult `json:"to"`
}

type NodeResult struct {
	DstIP     string `json:"dstIp"`
	Result    bool   `json:result`
	Timestamp string `json:"timestamp"`
	Reason    string `json:"reason"`
}

type Records []Record

type ReportSimp []RecordSimp

type ReportNodeToNode []RecordNodeToNode
