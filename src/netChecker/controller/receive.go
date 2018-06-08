package controller

import (
	"common/types"
	"github.com/garyburd/redigo/redis"
	"common/utils"
	"common/constant"
	"common/k8s"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"encoding/json"
	"time"
	"strings"
)

func ReceiveAndSaveData() {
	podToPodCh := make(chan Records)
	podToNodeCh := make(chan Records)
	podToSvcCh := make(chan Records)

	//通过管道传输report
	go getDataByChannel(constant.CHANNEL_REPORT_POD, podToPodCh)
	go getDataByChannel(constant.CHANNEL_REPORT_NODE, podToNodeCh)
	go getDataByChannel(constant.CHANNEL_REPORT_SVC, podToSvcCh)

	for {
		var reportAll Records

		svcReport := <-podToSvcCh
		log.Infof("ReceiveAndSaveData podToSvcCh return")
		nodeReport := <-podToNodeCh
		log.Infof("ReceiveAndSaveData podToNodeCh return")
		podReport := <-podToPodCh
		log.Infof("ReceiveAndSaveData podToPodCh return. getDataByChannel finish")

		reportAll = append(reportAll, svcReport...)
		reportAll = append(reportAll, nodeReport...)
		reportAll = append(reportAll, podReport...)

		reportSimp := getSimpleData(&reportAll)
		reportNTN := getNodeToNodeData(&nodeReport)

		saveData(reportAll, constant.ALL)
		saveData(reportAll, constant.FALSE)
		saveData(reportAll, constant.TRUE)
		saveData(reportSimp, constant.FALSE_SIMPLE)
		saveData(reportSimp, constant.TRUE_SIMPLE)
		saveData(reportNTN, constant.NOTE_TO_NODE)

		log.Infof("ReceiveAndSaveData save all data success")
	}
}

func getNodeLength() int {
	nodeLength := 0

	//节点数在网络测试的netcheck完成后，存于redis中
	nodelen := cache.Get("nodelength")
	if !strings.EqualFold(nodelen, "") {
		nodeLength, _ = strconv.Atoi(nodelen)
		log.Infof("receive getNodeLength from redis success. nodelength=%d", nodeLength)
	}

	log.Infof("-----getNodeLength after redis. redislen=%s, nodeLength=%d", nodelen, nodeLength)

	//从redis获取失败则从k8s中获取
	if nodeLength == 0 {
		client := k8s.KubeClient
		nodelist, err := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
		if err != nil {
			log.Errorf("receive getNodeLength from k8s failed. err=%s", err)
		} else {
			nodeLength = len(nodelist.Items)
			log.Infof("receive getNodeLength from k8s success. nodelength=%d", nodeLength)
		}
	}
	return nodeLength
}

//因为不同节点，对同一种类型（如：pod to pod/node/svc）进行网络检查的结果通过同一个频道发布消息。
//而消息的订阅在这里。为了完整的收集到一种类型一次检查的所有消息，下面2个条件满足任意一个时才进行数据的组装：
//1. 消息接受的次数达到集群的节点数（因为节点都发布自己的检查结果）
//2. 定时器结束。该定时器在每轮第一次收到消息后开始启动。目的在于解决某些节点发布检查结果消息失败，
//	导致条件一无法满足。如此一来，既无法获取实际的检查结果，且导致频道中的消息堆积，无法区分最新一轮的检查结果。
func getDataByChannel(channel string, reportCh chan Records) {
	var nodeCount, nodeLength int
	var subConn *redis.PubSubConn
	receivedCh := make(chan int)
	records, recordsTmp := Records{}, Records{}

	//todo 这个定时器的时间需考究。
	timeout := time.Duration(utils.LoadEnvVarInt(constant.ENV_RECEIVE_TIMEOUT, constant.RECEIVE_TIMEOUT))
	timer := time.NewTimer(time.Second * timeout)
	subConn = cache.GetSubConn(channel)
	if subConn == nil {
		subConn = cache.RetrySubConn(channel)
	}

	go func() {
		for {
			select {
			default:
				log.Infof("wait to ReceiveSubMessage")
				//这里的定时器是防止订阅的频道始终无数据输出，而导致ReceiveAndSaveData函数阻塞在相应的channel中，无法处理结果
				timer.Reset(time.Second * timeout)
				data := cache.ReceiveSubMessage(subConn)
				if data != nil {
					//这里的定时是为解决当订阅频道无数据数据输出时，程序阻塞在receive处，导致结果无法采集并存储于redis中
					timer.Reset(time.Second * timeout)

					nodeCount++
					if nodeCount == 1 {
						//当消息成功订阅时，获取实际的节点数
						nodeLength = getNodeLength()
					}
					err := json.Unmarshal(data, &recordsTmp)
					utils.CheckError("getDataByChannel: Unmarshal failed ", err)
					records = append(records, recordsTmp...)
					if nodeCount == nodeLength {
						//接受到的订阅消息数量达到节点数时，说明该轮的检查结果获取完毕
						receivedCh <- 1
					}
					log.Infof("getDataByChannel sub message success. channel=%s, count=%d", channel, nodeCount)
				}
			}
		}
	}()

	for {
		select {
		//本轮消息接受完毕
		case <-receivedCh:
			reportCh <- records
			data, _ := json.Marshal(records)
			log.Infof("getDataByChannel success. reason is reveive all data. channel=%s, length=%d, data=%s", channel, nodeCount, string(data))
			nodeCount = 0
			//清空数据，等待下一轮接收
			records = Records{}

			//本轮消息接收超时
		case <-timer.C:
			reportCh <- records
			data, _ := json.Marshal(records)
			log.Infof("getDataByChannel success. reason is timeout. channel=%s, length=%d,  data=%s", channel, nodeCount, string(data))
			nodeCount = 0
			//清空数据，等待下一轮接收
			records = Records{}
		}
	}
}

func saveData(report interface{}, flag int) {
	var err error
	var key string
	var result []byte
	var records Records
	var repSimp ReportSimp

	switch flag {
	case constant.ALL:
		key = constant.ALL_DATA
		result, err = json.Marshal(report)
	case constant.FALSE:
		reportAll := report.(Records)
		key = constant.FALSE_DATA
		for _, item := range reportAll {
			if !item.Result {
				records = append(records, item)
			}
		}
		result, err = json.Marshal(records)
	case constant.TRUE:
		reportAll := report.(Records)
		key = constant.TRUE_DATA
		for _, item := range reportAll {
			if item.Result {
				records = append(records, item)
			}
		}
		result, err = json.Marshal(records)
	case constant.FALSE_SIMPLE:
		reportSimp := report.(ReportSimp)
		key = constant.SIMPLE_FALSE_DATA
		for _, item := range reportSimp {
			if !item.Result {
				repSimp = append(repSimp, item)
			}
		}
		result, err = json.Marshal(repSimp)
	case constant.TRUE_SIMPLE:
		reportSimp := report.(ReportSimp)
		key = constant.SIMPLE_TRUE_DATA
		for _, item := range reportSimp {
			if item.Result {
				repSimp = append(repSimp, item)
			}
		}
		result, err = json.Marshal(repSimp)
	case constant.NOTE_TO_NODE:
		reportNTN := report.(ReportNodeToNode)
		key = constant.NOTE_TO_NODE_DATA
		result, err = json.Marshal(reportNTN)
	}

	utils.CheckError("saveData json marshal failed. key="+key, err)
	cache.SetWithExpire(key, string(result), strconv.Itoa(constant.REDIS_EXPIRE))
	log.Infof("saveData success. key=%s, data=%s", key, string(result))
}

func getSimpleData(reportAll *Records) ReportSimp {
	recordSimp := new(RecordSimp)
	records := Records{}
	reportSimp := ReportSimp{}
	podInfo := types.PodInfo{}
	nodeInfo := types.NodeInfo{}
	//svcInfo := types.ServiceInfo{}

	data, err := json.Marshal(reportAll)
	if err != nil {
		log.Errorf("getSimpleData json marshal failed. err=%s", err)
		return reportSimp
	}
	json.Unmarshal(data, &records)
	for _, item := range records {
		recordSimp.SrcName = item.From.Name
		recordSimp.Result = item.Result
		recordSimp.Timestamp = item.Timestamp
		recordSimp.Reason = item.Reason
		if strings.EqualFold(item.Type, "pod") {
			json.Unmarshal(item.To, &podInfo)
			recordSimp.DstName = podInfo.Name
		} else if strings.EqualFold(item.Type, "node") {
			json.Unmarshal(item.To, &nodeInfo)
			recordSimp.DstName = nodeInfo.Name
		} else if strings.EqualFold(item.Type, "service") {
			//json.Unmarshal(item.To, &svcInfo)
			//因为ping checknet时，检测service类型的item.To只封装了port和nodeport
			//由于检测的service只有一个，命名为pong-svc，映射到pond，故此处写死了
			recordSimp.DstName = "pong-svc"
		}
		reportSimp = append(reportSimp, *recordSimp)
	}
	log.Infof("getSimpleData success")
	return reportSimp
}

func getNodeToNodeData(nodeReport *Records) ReportNodeToNode {
	//组装node to node数据，方便前端展示
	if len(*nodeReport) != 0 {
		ntns := make(ReportNodeToNode, 0)

		ntn := new(RecordNodeToNode)
		ntn.To = make([]NodeResult, 0)

		srcIp := (*nodeReport)[0].From.HostIP
		for _, node := range *nodeReport {
			if strings.EqualFold(srcIp, node.From.HostIP) {
				ntn.SrcIp = node.From.HostIP
				result := getNTNResult(&node)
				ntn.To = append(ntn.To, *result)
			} else {
				ntns = append(ntns, *ntn)

				ntn = new(RecordNodeToNode)
				ntn.To = make([]NodeResult, 0)
				srcIp = node.From.HostIP

				ntn.SrcIp = node.From.HostIP
				result := getNTNResult(&node)
				ntn.To = append(ntn.To, *result)
			}
		}
		log.Infof("getNodeToNodeData success。 data=%v", ntns)
		return ntns
	}
	return nil
}

func getNTNResult(node *Record) *NodeResult {
	nodeInfo := types.NodeInfo{}
	err := json.Unmarshal(node.To, &nodeInfo)
	if err != nil {
		log.Errorf("getNodeToNodeData unmarshal failed. err=%s", err)
		return nil
	}
	result := &NodeResult{}
	result.DstIP = nodeInfo.HostIP
	result.Result = node.Result
	result.Reason = node.Reason
	result.Timestamp = node.Timestamp
	return result
}
