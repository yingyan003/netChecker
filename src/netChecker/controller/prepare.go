package controller

import (
	"common/k8s"
	"time"
	"strings"
	kube "k8s.io/client-go/kubernetes"
	"encoding/json"
	"common/constant"
	"common/utils"
	"common/types"
	mylog "github.com/maxwell92/gokits/log"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

var log *mylog.Logger
var cache *utils.RedisClient
var redisExpire string

func Init() {
	log = utils.GetLog()
	k8s.NewKubeClient()
	utils.NewRedis(constant.NETGUARD_MAX_IDLE)
	cache = utils.Redis
	redisExpire = utils.LoadEnvVar(constant.ENV_REDIS_EXPIRE, constant.REDIS_EXPIRE)
}

//定期获取k8s资源并写入redis
func GetKubeResAndPublish() {
	client := k8s.KubeClient
	ticker := time.NewTicker(time.Duration(utils.LoadEnvVarInt(constant.ENV_GET_RESOURCE_TICKER, constant.GET_RESOURCE_TICKER)) * time.Second)
	for _ = range ticker.C {
		pings := getPings(client)
		go getPonds(client, pings)
		go getNodes(client, pings)
		go getSvc(client, pings)
	}
}

func getPings(client *kube.Clientset) []*types.PodInfo {
	pingInfos := make([]*types.PodInfo, 0)
	pods, err := client.CoreV1().Pods("yce").List(meta_v1.ListOptions{LabelSelector: "name=ping"})
	if err != nil {
		log.Errorf("getPings: list pods failed. err=%s, pods=%s", err, pods)
		return pingInfos
	}

	for _, pod := range pods.Items {
		p := new(types.PodInfo)
		p.Name = pod.Name
		p.PodIP = pod.Status.PodIP
		p.HostIP = pod.Status.HostIP
		p.Status = string(pod.Status.Phase)
		pingInfos = append(pingInfos, p)
	}
	log.Infof("getPings success")
	return pingInfos
}

func getPonds(client *kube.Clientset, pingInfos []*types.PodInfo) {
	pongInfos := make([]*types.PodInfo, 0)
	pods, err := client.CoreV1().Pods("yce").List(meta_v1.ListOptions{LabelSelector: "name=pong"})
	if err != nil {
		log.Errorf("getPonds: list pods failed. err=%s. pod=%s", err, pods)
	}
	for _, pod := range pods.Items {
		p := new(types.PodInfo)
		p.Name = pod.Name
		p.PodIP = pod.Status.PodIP
		p.HostIP = pod.Status.HostIP
		p.Status = string(pod.Status.Phase)
		pongInfos = append(pongInfos, p)
	}

	checkDeamonsetCover(client, len(pingInfos), len(pongInfos))
	Publist(constant.CHANNEL_POD, pingInfos, pongInfos)
}

func getNodes(client *kube.Clientset, pingInfos []*types.PodInfo) {
	nodeInfos := make([]*types.NodeInfo, 0)
	list, err := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		log.Errorf("prepare getNodes: list nodes failed. err=%s. nodelist=%s", err, list)
	}
	if list != nil {
		log.Infof("prepare getNodes: Node count %d", len(list.Items))
		for _, node := range list.Items {
			n := new(types.NodeInfo)
			n.Name = node.Name
			for _, addr := range node.Status.Addresses {
				if strings.EqualFold(string(addr.Type), "LegacyHostIP") {
					n.HostIP = addr.Address
				}
			}
			n.Status = string(node.Status.Phase)
			nodeInfos = append(nodeInfos, n)
		}
	}

	//将node数通过一般的key-value方式存于Redis中，目的在于netgaurd订阅获取netcheck结果后，组装数据
	if ok, err := cache.SetWithExpire("nodelength", strconv.Itoa(len(nodeInfos)),redisExpire); !ok {
		log.Infof("prepare getNodes save nodelength to redis failed. nodelength=%d, err=%s", len(nodeInfos), err)
	}

	Publist(constant.CHANNEL_NODE, pingInfos, nodeInfos)
}

func getSvc(client *kube.Clientset, pingInfos []*types.PodInfo) {
	svcInfo := types.ServiceInfo{}
	svc, err := client.CoreV1().Services("yce").Get("pong-svc", meta_v1.GetOptions{})
	if err != nil {
		log.Errorf("prepare getSvc: get svc failed. err=%s. svc=%s", err, svc)
	}
	if svc != nil {
		svcInfo.Name = svc.Name
		svcInfo.Type = string(svc.Spec.Type)
		svcInfo.ClusterIP = svc.Spec.ClusterIP
		for _, p := range svc.Spec.Ports {
			port := types.Port{
				Port:     p.Port,
				NodePort: p.NodePort,
			}
			svcInfo.Ports = append(svcInfo.Ports, port)
		}
	}
	Publist(constant.CHANNEL_SVC, pingInfos, svcInfo)
}

func checkDeamonsetCover(client *kube.Clientset, pinglen, ponglen int) {
	nodeList, err := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	utils.CheckError("Get NodeList", err)

	if pinglen != ponglen {
		log.Errorf("PingList doesn't match PongList. len(ping)=%d, len(pong)=%d", pinglen, ponglen)
	}

	if pinglen != len(nodeList.Items) || ponglen != len(nodeList.Items) {
		log.Errorf("PingList doesn't match NodeList or PongList doesn't match NodeList. pinglist=%d, ponglist=%d, nodelist=%d", pinglen, ponglen, len(nodeList.Items))
	}
}

func Publist(channel string, pingInfos []*types.PodInfo, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Publist marshal data failed.")
	}
	pb := new(types.PubSubInfo)
	pb.PingInfos = pingInfos
	pb.Data = bytes
	message, err := json.Marshal(pb)
	utils.CheckError("Publist JsonMarshal", err)
	ok, err := cache.Publish(channel, message)
	if !ok {
		log.Errorf("Publist failed. err=%v, channel=%s, message=%s", err, channel, message)
		return
	}
	log.Infof("Publist success. channel=%s", channel)
}
