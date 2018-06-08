package types

import "encoding/json"

type Report []*Record

type Record struct {
	Type   string `json:"type"`
	From   interface{} `json:"from"`
	To     interface{} `json:"to"`
	Result bool `json:"result"`
	Reason string `json:"reason"`
	Timestamp   string `json:"timestamp"`
}

//type Report struct {
//	Items []Record `json:"items"`
//}

type PodInfo struct {
	Name   string `json:"name"`
	PodIP  string `json:"podIP"`
	HostIP string `json:"hostIP"`
	Status string `json:"status"`
}

type NodeInfo struct {
	Name   string `json:"name"`
	HostIP string `json:"hostIP"`
	Status string `json:"status"`
}

type ServiceInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	ClusterIP string `json:"clusterIP"`
	Ports     []Port `json:"ports"`
}

type Port struct {
	Port     int32 `json:"port"`
	NodePort int32 `json:"nodePort"`
}

type Config struct {
	Apiserver string
	RedisHost string
	NodePort  int32
}

type PubSubInfo struct{
	PingInfos []*PodInfo
	Data json.RawMessage
}



