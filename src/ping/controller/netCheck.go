package controller

import (
	cach "github.com/maxwell92/gokits/cache"
	"common/sysinit"
	"common/types"
	mylog "github.com/maxwell92/gokits/log"
	"common/utils"
	"strconv"
	"os"
	"encoding/json"
	"common/constant"
	"time"
)

var log *mylog.Logger
var conf *types.Config
var cache *cach.RedisCache

func Init(){
	log = utils.GetLog()
	conf= sysinit.GetConfig()
	cache = cach.NewRedisCache()
}

func NetCheck() {
	timer := time.NewTicker(time.Duration(utils.LoadEnvVarInt(constant.ENV_TELNET_TIMER, constant.TELNET_TIMER)) * time.Second)
	for _ = range timer.C {
		CheckPod()
		CheckNode(conf.NodePort)
		CheckService()
	}
}

func CheckPod() {
	var report types.Report
	self := types.PodInfo{}
	pingInfoList := make([]types.PodInfo, 0)
	pongInfoList := make([]types.PodInfo, 0)

	pingdata := cache.Get("pinglist")
	if pingdata == "" {
		return
	}
	pongdata := cache.Get("ponglist")
	if pongdata == "" {
		return
	}

	err := json.Unmarshal([]byte(pingdata), &pingInfoList)
	utils.CheckError("TestPod JsonUnmarshal", err)

	getSelfPod(&self, &pingInfoList)
	log.Infof("CheckPod PodName=%s, HostIP=%s, PodIP=%s, Status=%s", self.Name, self.HostIP, self.PodIP, self.Status)

	err = json.Unmarshal([]byte(pongdata), &pongInfoList)
	utils.CheckError("CheckPod JsonUnmarshal", err)

	for _, pod := range pongInfoList {
		r := netTest(pod.PodIP, 8080, "pod", &self, &pod)
		report = append(report, r)
	}

	saveToRedis(report, self.HostIP, "pod")
	log.Infoln("CheckPod Success")
}

func CheckNode(nodeport string) {
	var report types.Report
	var pingdata, data string
	self := types.PodInfo{}
	pingInfoList := make([]types.PodInfo, 0)
	nodeInfoList := make([]types.NodeInfo, 0)

	if pingdata = cache.Get("pinglist"); pingdata == "" {
		return
	}
	err := json.Unmarshal([]byte(pingdata), &pingInfoList)
	utils.CheckError("CheckNode JsonUnmarshal", err)

	getSelfPod(&self, &pingInfoList)
	log.Infof("CheckNode PodName=%s, HostIP=%s, PodIP=%s, Status=%s", self.Name, self.HostIP, self.PodIP, self.Status)

	if data = cache.Get("nodelist"); data == "" {
		return
	}
	json.Unmarshal([]byte(data), &nodeInfoList)

	np, _ := strconv.Atoi(nodeport)
	for _, node := range nodeInfoList {
		r := netTest(node.HostIP, int32(np), "node", &self, &node)
		report = append(report, r)
	}

	saveToRedis(report, self.HostIP, "node")
	log.Infoln("CheckNode Success")
}

func CheckService() {
	var report types.Report
	var pingdata, data string
	self := types.PodInfo{}
	serviceInfo := types.ServiceInfo{}
	pingInfoList := make([]types.PodInfo, 0)

	if pingdata = cache.Get("pinglist"); pingdata == "" {
		return
	}

	err := json.Unmarshal([]byte(pingdata), &pingInfoList)
	utils.CheckError("CheckService JsonUnmarshal", err)

	getSelfPod(&self, &pingInfoList)
	log.Infof("CheckService PodName=%s, HostIP=%s, PodIP=%s, Status=%s", self.Name, self.HostIP, self.PodIP, self.Status)

	if data = cache.Get("service"); data == "" {
		return
	}
	json.Unmarshal([]byte(data), &serviceInfo)

	for _, port := range serviceInfo.Ports {
		r := netTest(serviceInfo.ClusterIP, port.Port, "service", &self, &port)
		report = append(report, r)
	}

	saveToRedis(report, self.HostIP, "service")
	log.Infoln("CheckService Success")
}

//从pog列表中获取本容器所在的pod
func getSelfPod(self *types.PodInfo, list *[]types.PodInfo) {
	hostname, _ := os.Hostname()
	for _, pod := range *list {
		if hostname == pod.Name {
			*self = pod
			break
		}
	}
}

func netTest(ip string, port int32, tp string, from, to interface{}) *types.Record {
	var podIP string
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	go goTelnet(ip, port, ch1)
	go goTelnet("www.baidu.com", 80, ch2)
	connected := <-ch1
	debugConn := <-ch2

	r := &types.Record{
		Type:   tp,
		From:   from,
		To:     to,
		Result: connected,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
	if !connected && !debugConn {
		r.Reason = constant.DEBUG_FAIL
	}

	if f, ok := from.(*types.PodInfo); ok {
		podIP = f.PodIP
	}

	switch to.(type) {
	case *types.PodInfo:
		pod, _ := to.(*types.PodInfo)
		log.Errorf("CheckPod failed.from=%s, to=%s, result=%v, reason=%s\n", podIP, pod.PodIP, r.Result, r.Reason)
	case *types.NodeInfo:
		node, _ := to.(*types.NodeInfo)
		log.Errorf("CheckNode failed.from=%s, to=%s, result=%v, reason=%s\n", podIP, node.HostIP, r.Result, r.Reason)
	case *types.Port:
		port, _ := to.(*types.Port)
		log.Errorf("CheckNode failed.from=%s, to=%s, result=%v, reason=%s\n", podIP, port, r.Result, r.Reason)
	}

	return r
}

func goTelnet(ip string, port int32, ch chan bool) {
	connected := utils.Telnet(ip, port, 1)
	ch <- connected
}

func saveToRedis(report []*types.Record, hostIP, types string) {
	var key string
	result, _ := json.Marshal(report)
	switch types {
	case "pod":
		key = "pod:" + hostIP + ":pod"
	case "node":
		key = "pod:" + hostIP + ":noce"
	case "service":
		key = "pod:" + hostIP + ":service"
	}
	cache.Set(key, string(result))
	log.Infof("saveToRedis key=%s, value=%s\n", key, string(result))
}
