package cache

import (
	"encoding/json"
	"fmt"
)

type DeployInfo struct {
	UserName       string `json:"userName"`
	DcName         string `json:"dcName"`
	DcID           int32  `json:"dcId"`
	OrgName        string `json:"orgName"`
	DeploymentName string `json:"deploymentName"`
	UpdateTime     string `json:"updateTime"`
	// Deployment interface{} `json:"deployoment"`
}

type Deploy struct {
	Data []DeployInfo `json:"data"`
}

func (deploy *Deploy) Unmarshal(content *[]byte) error {

	err := json.Unmarshal(*content, deploy)
	if err != nil {
		fmt.Printf("Json unmarshal error: %s\n", err)
		return err
	}
	return nil
}

type Indexer struct {
	Index string
	Key   string
}
