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

//从redis订阅频道中，采集网络检测结果的数据，然后组装存入redis中
func ReceiveAndSaveData() {
	termFinish := 0
	podToPodCh := make(chan int, 1)
	podToNodeCh := make(chan int, 1)
	podToSvcCh := make(chan int, 1)

	//通过管道传输report
	go getDataByChannel(constant.CHANNEL_REPORT_POD, podToPodCh)
	go getDataByChannel(constant.CHANNEL_REPORT_NODE, podToNodeCh)
	go getDataByChannel(constant.CHANNEL_REPORT_SVC, podToSvcCh)

	for {
		select {
		case <-podToSvcCh:
			termFinish++
		case <-podToNodeCh:
			termFinish++
		case <-podToPodCh:
			termFinish++
		default:
			//防止程序阻塞
		}

		//3种类型的数据采集完毕：pod->pod,pod->node,pod->svc
		if termFinish >= 3 {
			//重置新一轮的计时
			termFinish = 0
			//组装数据，存入redis
			getAndSave()
		}
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
func getDataByChannel(channel string, reportCh chan int) {
	var nodeCount, nodeLength int
	var subConn *redis.PubSubConn
	var newTermTime time.Time
	newTerm := false
	records, recordsTmp := Records{}, Records{}

	timeout := utils.LoadEnvVarInt(constant.ENV_RECEIVE_TIMEOUT, constant.RECEIVE_TIMEOUT)
	//todo 这里只在当前goroutine执行时初始化连接redis。当下面goroutine中的for循环在执行时，redis连接断了，
	//todo 那么就在也连不上redis了，除非重启netgaurd的pod。后续考虑在下面的for循环中处理连接情况。
	subConn = cache.GetSubConn(channel)
	if subConn == nil {
		subConn = cache.RetrySubConn(channel)
	}

	go func() {
		for {
			//本轮数据采集结束。条件：1.是新的一轮 2.所有数据成功采集或数据采集超时
			if newTerm && (nodeCount >= nodeLength || time.Now().Sub(newTermTime).Seconds() > float64(timeout) ) {
				result, err := json.Marshal(records)
				utils.CheckError("getDataByChannel: json marshal failed.", err)
				//将采集结果存入redis，该结果只是某种类型的网络检测的全部结果，如svc/node/pod
				cache.SetWithExpire(channel, string(result), redisExpire)
				log.Infof("getDataByChannel: get and save all data success. channel=%s, datalen=%d, nodelen=%d", channel, nodeCount, nodeLength)
				newTerm = false
				nodeCount = 0
				records = Records{}
				//通知函数ReceiveAndSaveData()，本类型本轮数据采集完毕
				reportCh <- 1
			}
			log.Infof("wait to ReceiveSubMessage")
			data := cache.ReceiveSubMessage(subConn)
			if data != nil {
				nodeCount++
				if nodeCount == 1 {
					//当消息成功订阅时，获取实际的节点数
					nodeLength = getNodeLength()
					//标记新一轮的数据采集
					newTerm = true
					newTermTime = time.Now()
					log.Infof("getDataByChannel: new term come. newTermTime=%s", newTermTime.String())
				}
				err := json.Unmarshal(data, &recordsTmp)
				utils.CheckError("getDataByChannel: Unmarshal failed ", err)
				records = append(records, recordsTmp...)
				log.Infof("getDataByChannel: sub message success. channel=%s, count=%d", channel, nodeCount)
			}
		}
	}()
}

func getAndSave() {
	var podReport, nodeReport, svcReport, reportAll Records

	dataPod := cache.Get(constant.CHANNEL_REPORT_POD)
	dataNode := cache.Get(constant.CHANNEL_REPORT_NODE)
	dataSvc := cache.Get(constant.CHANNEL_REPORT_SVC)

	if !strings.EqualFold(dataPod, "") {
		err := json.Unmarshal([]byte(dataPod), &podReport)
		utils.CheckError("ReceiveAndSaveData: dataPod json Unmarshal failed.", err)
	}
	if !strings.EqualFold(dataNode, "") {
		err := json.Unmarshal([]byte(dataNode), &nodeReport)
		utils.CheckError("ReceiveAndSaveData: dataNode json Unmarshal failed.", err)
	}
	if !strings.EqualFold(dataSvc, "") {
		err := json.Unmarshal([]byte(dataSvc), &svcReport)
		utils.CheckError("ReceiveAndSaveData: dataSvc json Unmarshal failed.", err)
	}

	log.Infof("ReceiveAndSaveData: all data get success. podslen=%d, nodeslen=%d, svcslen=%d", len(podReport), len(nodeReport), len(svcReport))

	reportAll = append(reportAll, svcReport...)
	reportAll = append(reportAll, nodeReport...)
	reportAll = append(reportAll, podReport...)

	reportSimp := getSimpleData(&reportAll)
	reportNTN := getNTNDataFromPodReport(&nodeReport)

	saveData(reportAll, constant.ALL)
	saveData(reportAll, constant.FALSE)
	saveData(reportAll, constant.TRUE)
	saveData(reportSimp, constant.FALSE_SIMPLE)
	saveData(reportSimp, constant.TRUE_SIMPLE)
	saveData(reportNTN, constant.NOTE_TO_NODE)

	log.Infof("ReceiveAndSaveData save all data success")
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
	cache.SetWithExpire(key, string(result), redisExpire)
	if string(result) != "null" {
		log.Infof("saveData success. key=%s, data!=null", key)
	} else {
		log.Infof("saveData success. key=%s, data=%s", key, string(result))
	}
}

func getSimpleData(reportAll *Records) ReportSimp {
	recordSimp := new(RecordSimp)
	records := Records{}
	reportSimp := ReportSimp{}
	podInfo := types.PodInfo{}
	nodeInfo := types.NodeInfo{}

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
			recordSimp.DstName = "pong-svc"
		}
		reportSimp = append(reportSimp, *recordSimp)
	}
	log.Infof("getSimpleData success")
	return reportSimp
}

//从pod->pod的网络检测结果中解析node->node的网络信息
func getNTNDataFromPodReport(podReport *Records) ReportNodeToNode {
	if len(*podReport) != 0 {
		ntns := make(ReportNodeToNode, 0)

		ntn := new(RecordNodeToNode)
		ntn.To = make([]NodeResult, 0)

		srcIp := (*podReport)[0].From.HostIP
		for i, pod := range *podReport {
			if strings.EqualFold(srcIp, pod.From.HostIP) {
				ntn.SrcIp = pod.From.HostIP
				result := getNTNResultFromPod(&pod)
				ntn.To = append(ntn.To, *result)
				//最后一组数据
				if i == len(*podReport)-1 {
					ntns = append(ntns, *ntn)
				}
			} else {
				ntns = append(ntns, *ntn)

				ntn = new(RecordNodeToNode)
				ntn.To = make([]NodeResult, 0)
				srcIp = pod.From.HostIP

				ntn.SrcIp = pod.From.HostIP
				result := getNTNResultFromPod(&pod)
				ntn.To = append(ntn.To, *result)
			}
		}
		log.Infof("getNTNDataFromPodReport success。 podReport len=%d, nodeToNodeReport len=%d", len(*podReport), len(ntns))
		return ntns
	}
	log.Errorf("getNTNDataFromPodReport failed。 podReport=%v", *podReport)
	return nil
}

func getNTNResultFromPod(pod *Record) *NodeResult {
	podInfo := types.PodInfo{}
	err := json.Unmarshal(pod.To, &podInfo)
	if err != nil {
		log.Errorf("getNTNResultFromPod unmarshal failed. err=%s", err)
		return nil
	}
	result := &NodeResult{}
	result.DstIP = podInfo.HostIP
	result.Result = pod.Result
	result.Reason = pod.Reason
	result.Timestamp = pod.Timestamp
	return result
}

//从pod->node的网络检测结果中解析node->node的网络信息
func getNTNDataFromNodeReport(nodeReport *Records) ReportNodeToNode {
	//组装node to node数据，方便前端展示
	if len(*nodeReport) != 0 {
		ntns := make(ReportNodeToNode, 0)

		ntn := new(RecordNodeToNode)
		ntn.To = make([]NodeResult, 0)

		srcIp := (*nodeReport)[0].From.HostIP
		for i, node := range *nodeReport {
			if strings.EqualFold(srcIp, node.From.HostIP) {
				ntn.SrcIp = node.From.HostIP
				result := getNTNResultFromNode(&node)
				ntn.To = append(ntn.To, *result)
				//最后一组数据
				if i == len(*nodeReport) {
					ntns = append(ntns, *ntn)
				}
			} else {
				ntns = append(ntns, *ntn)

				ntn = new(RecordNodeToNode)
				ntn.To = make([]NodeResult, 0)
				srcIp = node.From.HostIP

				ntn.SrcIp = node.From.HostIP
				result := getNTNResultFromNode(&node)
				ntn.To = append(ntn.To, *result)
			}
		}
		log.Infof("getNTNDataFromNodeReport success。 data=%v", ntns)
		return ntns
	}
	return nil
}

func getNTNResultFromNode(node *Record) *NodeResult {
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
