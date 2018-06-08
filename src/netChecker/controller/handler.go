package controller

import (
	"fmt"
	"net/http"
	"encoding/json"
	"common/constant"
	"common/types"
	"common/k8s"
	"k8s.io/client-go/pkg/api/v1"
	"common/utils"
)

func True(w http.ResponseWriter, r *http.Request) {
	result := getdata(constant.TRUE)
	fmt.Fprintf(w, "%s\n", result)
}

func False(w http.ResponseWriter, r *http.Request) {
	result := getdata(constant.FALSE)
	fmt.Fprintf(w, "%s\n", result)
}

func All(w http.ResponseWriter, r *http.Request) {
	result := getdata(constant.ALL)
	fmt.Fprintf(w, "%s\n", result)
}

func TrueSimple(w http.ResponseWriter, r *http.Request) {
	result := getSimpData(constant.TRUE_SIMPLE)
	fmt.Fprintf(w, "%s\n", result)
}

func FalseSimple(w http.ResponseWriter, r *http.Request) {
	result := getSimpData(constant.FALSE_SIMPLE)
	fmt.Fprintf(w, "%s\n", result)
}

func getdata(flag int) string {
	var report types.Report
	nodeList := make([]types.NodeInfo, 0)
	data := cache.Get("nodelist")

	//todo 考虑保留与否
	if data == "" {
		client := k8s.GetKubeClient(conf.Apiserver)
		list, _ := client.Nodes().List(v1.ListOptions{})
		for _, node := range list.Items {
			n := new(types.NodeInfo)
			n.Name = node.Name
			for _, addr := range node.Status.Addresses {
				if addr.Type == "LegacyHostIP" {
					n.HostIP = addr.Address
				}
			}
			n.Status = string(node.Status.Phase)
			nodeList = append(nodeList, *n)
		}
	} else {
		json.Unmarshal([]byte(data), &nodeList)
	}

	for _, node := range nodeList {
		podRecords := getRecordByKey("pod:" + node.HostIP + ":pod")
		nodeRecords := getRecordByKey("pod:" + node.HostIP + ":node")
		svcRecords := getRecordByKey("pod:" + node.HostIP + ":service")
		report = append(report, podRecords...)
		report = append(report, nodeRecords...)
		report = append(report, svcRecords...)
	}

	return retData(flag, report)
}

func getSimpData(flag int) string {
	var rs ReportSimp
	nodeList := make([]types.NodeInfo, 0)
	data := cache.Get("nodelist")
	json.Unmarshal([]byte(data), &nodeList)

	for _, node := range nodeList {
		podRecords := getSimpDataByKey("pod:" + node.HostIP + ":pod")
		nodeRecords := getSimpDataByKey("pod:" + node.HostIP + ":node")
		svcRecords := getSimpDataByKey("pod:" + node.HostIP + ":service")
		rs = append(rs, podRecords...)
		rs = append(rs, nodeRecords...)
		rs = append(rs, svcRecords...)
	}

	return retData(flag, rs)
}

func getRecordByKey(key string) types.Report {
	var r types.Report
	data := cache.Get(key)
	err := json.Unmarshal([]byte(data), &r)
	utils.CheckError("getRedisDataByKey: "+key, err)
	return r
}

func getSimpDataByKey(key string) ReportSimp {
	var rjs RecordJsons
	var rps ReportSimp
	rcs := new(RecordSimp)
	podInfo := types.PodInfo{}
	nodeInfo := types.NodeInfo{}
	svcInfo := types.ServiceInfo{}

	data := cache.Get(key)
	json.Unmarshal([]byte(data), &rjs)
	for _, item := range rjs {
		rcs.SrcName = item.From.Name
		rcs.Result = item.Result
		rcs.Time = item.Time
		rcs.Reason = item.Reason
		if item.Type == "pod" {
			json.Unmarshal(item.To, &podInfo)
			rcs.DstName = podInfo.Name
		} else if item.Type == "node" {
			json.Unmarshal(item.To, &nodeInfo)
			rcs.DstName = nodeInfo.Name
		} else if item.Type == "service" {
			json.Unmarshal(item.To, &svcInfo)
			rcs.DstName = svcInfo.Name
		}
		rps = append(rps, rcs)
	}
	return rps
}

func retData(flag int, report interface{}) string {
	var result []byte
	rp := types.Report{}
	rs := ReportSimp{}

	switch flag {
	case constant.ALL:
		{
			r, _ := report.(types.Report)
			result, _ = json.Marshal(r)
		}
	case constant.TRUE:
		{
			r, _ := report.(types.Report)
			for _, item := range r {
				if item.Result.(bool) {
					rp = append(rp, item)
				}
			}
			result, _ = json.Marshal(rp)
		}
	case constant.FALSE:
		{
			r, _ := report.(types.Report)
			for _, item := range r {
				if !(item.Result.(bool) ) {
					rp = append(rp, item)
				}
			}
			result, _ = json.Marshal(rp)
		}
	case constant.TRUE_SIMPLE:
		{
			r, _ := report.(ReportSimp)
			for _, item := range r {
				if item.Result.(bool) {
					rs = append(rs, item)
				}
			}
			result, _ = json.Marshal(rs)
		}
	case constant.FALSE_SIMPLE:
		{
			r, _ := report.(ReportSimp)
			for _, item := range r {
				if !(item.Result.(bool)) {
					rs = append(rs, item)
				}
			}
			result, _ = json.Marshal(rs)
		}
	}

	log.Infof("%s", string(result))
	return string(result)
}
