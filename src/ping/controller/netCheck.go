package controller

import (
	"common/types"
	mylog "github.com/maxwell92/gokits/log"
	"common/utils"
	"os"
	"encoding/json"
	"common/constant"
	"time"
	"github.com/garyburd/redigo/redis"
	"strings"
)

var log *mylog.Logger
var cache *utils.RedisClient

func Init() {
	log = utils.GetLog()
	utils.NewRedis(constant.NETCHECK_MAX_IDLE)
	cache = utils.Redis
}

func NetCheck() {
	waitCh := make(chan int)

	go subPod()
	go subNode()
	go subSvc()

	//防止主程序结束
	<-waitCh
}

func subPod() {
	var subConn *redis.PubSubConn
	subConn = cache.GetSubConn(constant.CHANNEL_POD)
	if subConn == nil {
		subConn=cache.RetrySubConn(constant.CHANNEL_POD)
	}
	for {
		var report types.Report
		var pongInfos []*types.PodInfo
		psInfo := new(types.PubSubInfo)

		data := cache.ReceiveSubMessage(subConn)
		if data != nil {
			log.Infof("subPod: ReceiveSubMessage success. message=%s", string(data))
			err := json.Unmarshal(data, psInfo)
			if err != nil {
				log.Errorf("subPod Unmarshal psInfo failed. err=%s, psInfo=%s", err, psInfo)
			}
			err = json.Unmarshal(psInfo.Data, &pongInfos)
			if err != nil {
				log.Errorf("subPod Unmarshal pongInfos failed. err=%s, pongInfos=%s", err, pongInfos)
			}
			self := getSelfPod(psInfo.PingInfos)
			if self == nil {
				log.Errorf("subPod getSelfPod failed.")
				continue
			}
			log.Infof("subPod getSelfPod success")
			for _, pod := range pongInfos {
				//TODO 注意：pongInfos是指针类型的切片，对指针类型的切片迭代时，注意pod的指针，会随着迭代结束，通通指向最后一次迭代的值
				r := netTest(pod.PodIP, 8080, "pod", &self, pod)
				report = append(report, r)
			}
			publist(constant.CHANNEL_REPORT_POD, report)
		}
	}
}

func subNode() {
	var subConn *redis.PubSubConn

	nodePort := int32(utils.LoadEnvVarInt(constant.ENV_NODEPORT, constant.NODEPORT))
	subConn = cache.GetSubConn(constant.CHANNEL_NODE)
	if subConn == nil {
		subConn=cache.RetrySubConn(constant.CHANNEL_NODE)
	}
	for {
		var report types.Report
		var nodeInfos []*types.NodeInfo
		psInfo := new(types.PubSubInfo)

		data := cache.ReceiveSubMessage(subConn)
		if data != nil {
			log.Infof("subNode: ReceiveSubMessage success. message=%s", string(data))
			err := json.Unmarshal(data, psInfo)
			if err != nil {
				log.Errorf("subNode Unmarshal psInfo failed. err=%s, psInfo=%s", err, psInfo)
			}
			err = json.Unmarshal(psInfo.Data, &nodeInfos)
			if err != nil {
				log.Errorf("subNode Unmarshal nodeInfos failed. err=%s, nodeInfos=%s", err, nodeInfos)
			}
			utils.CheckError("subNode JsonUnmarshal", err)
			self := getSelfPod(psInfo.PingInfos)
			if self == nil {
				log.Errorf("subNode getSelfPod failed.")
				continue
			}
			log.Infof("subNode getSelfPod success")
			for _, node := range nodeInfos {
				////TODO 注意：nodeInfos是指针类型的切片，对指针类型的切片迭代时，注意node的指针，会随着迭代结束，通通指向最后一次迭代的值
				r := netTest(node.HostIP, nodePort, "node", &self, node)
				report = append(report, r)
			}
			publist(constant.CHANNEL_REPORT_NODE, report)
		}
	}
}

func subSvc() {
	var subConn *redis.PubSubConn

	subConn = cache.GetSubConn(constant.CHANNEL_SVC)
	if subConn == nil {
		subConn=cache.RetrySubConn(constant.CHANNEL_POD)
	}
	for {
		var report types.Report
		psInfo := new(types.PubSubInfo)
		svcInfo := new(types.ServiceInfo)

		data := cache.ReceiveSubMessage(subConn)
		if data != nil {
			log.Infof("subSvc: ReceiveSubMessage success. message=%s", string(data))
			err := json.Unmarshal(data, psInfo)
			if err != nil {
				log.Errorf("subSvc Unmarshal psInfo failed. err=%s, psInfo=%s", err, psInfo)
			}
			err = json.Unmarshal(psInfo.Data, svcInfo)
			if err != nil {
				log.Errorf("subSvc Unmarshal svcInfo failed. err=%s, svcInfo=%s", err, svcInfo)
			}
			utils.CheckError("subSvc JsonUnmarshal", err)
			self := getSelfPod(psInfo.PingInfos)
			if self == nil {
				log.Errorf("subSvc getSelfPod failed.")
				continue
			}
			log.Infof("subSvc getSelfPod success. self=%s",*self)
			for _, port := range svcInfo.Ports {
				//这里实际上只测试了svc的port端口，并没有测试nodeport端口
				r := netTest(svcInfo.ClusterIP, port.Port, "service", &self, port)
				report = append(report, r)
			}
			publist(constant.CHANNEL_REPORT_SVC, report)
		}
	}
}

func publist(channel string, data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Publist marshal data failed.")
	}

	ok, err := cache.Publish(channel, message)
	if !ok {
		log.Errorf("Publist failed. err=%v, channel=%s, message=%s", err, channel, message)
		return
	}
	log.Infof("Publist success. channel=%s, message=%s", channel, message)
}

//从pog列表中获取本容器所在的pod
func getSelfPod(list []*types.PodInfo) *types.PodInfo {
	hostname, _ := os.Hostname()
	for _, pod := range list {
		if strings.EqualFold(hostname, pod.Name) {
			log.Infof("getSelfPod success. hostname=%s, selfPod=%s", hostname, pod)
			return pod
		}
	}
	return nil
}

func netTest(ip string, port int32, tp string, from, to interface{}) *types.Record {
	var podIP string
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	go goTelnet(ip, port, ch1)
	go goTelnet("www.baidu.com", 80, ch2)
	connected := <-ch1
	debugConn := <-ch2
	
	//TODO 注意：由于to是指针类型切片的迭代结果，如果直接将to赋值给Record中的To时，所有的Record中的To的值，都将变为to最后一次迭代的结果
	toTmp:=to
	r := &types.Record{
		Type:   tp,
		From:   from,
		To:     toTmp,
		Result: connected,
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
		Reason: "",
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
		log.Infof("CheckPod success.from=%s, to=%s, result=%v, reason=%s\n", podIP, pod.PodIP, r.Result, r.Reason)
	case *types.NodeInfo:
		node, _ := to.(*types.NodeInfo)
		log.Infof("CheckNode success.from=%s, to=%s, result=%v, reason=%s\n", podIP, node.HostIP, r.Result, r.Reason)
	case *types.Port:
		port, _ := to.(*types.Port)
		log.Infof("CheckSvc success.from=%s, to=%s, result=%v, reason=%s\n", podIP, port, r.Result, r.Reason)
	}

	return r
}

func goTelnet(ip string, port int32, ch chan bool) {
	connected := utils.Telnet(ip, port, 1)
	ch <- connected
}
