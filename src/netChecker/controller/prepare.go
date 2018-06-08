package controller

import (
	"common/k8s"
	"time"
	cach "github.com/maxwell92/gokits/cache"
	"strings"
	"k8s.io/client-go/pkg/api/v1"
	kube "k8s.io/client-go/kubernetes"
	"common/sysinit"
	"encoding/json"
	"common/constant"
	"common/utils"
	"strconv"
	"common/types"
	mylog "github.com/maxwell92/gokits/log"
)

var log *mylog.Logger
var conf *types.Config
var cache *cach.RedisCache

func Init(){
	log = utils.GetLog()
	conf= sysinit.GetConfig()
	cache = cach.NewRedisCache()
}

//定期获取k8s资源并写入redis
func GetKubeResToRedis() {
	client := k8s.GetKubeClient(conf.Apiserver)
	timer := time.NewTicker(time.Duration(utils.LoadEnvVarInt(constant.ENV_RESOURCE_TIMER, constant.RESOURCE_TIMER)) * time.Second)
	for _ = range timer.C {
		checkDeamonsetCover(client)
		GetPodList(client)
		GetNodeList(client)
		GetSvc(client)
	}
}

func checkDeamonsetCover(client *kube.Clientset) {
	podList, err := client.Pods("yce").List(v1.ListOptions{})
	utils.CheckError("Get PodList", err)
	pingList := new(v1.PodList)
	pongList := new(v1.PodList)
	for _, pod := range podList.Items {
		if strings.Contains(pod.Name, "ping-") {
			pingList.Items = append(pingList.Items, pod)
		}

		if strings.Contains(pod.Name, "pong-") {
			pongList.Items = append(pongList.Items, pod)
		}
	}

	nodeList, err := client.Nodes().List(v1.ListOptions{})
	utils.CheckError("Get NodeList", err)

	if len(pingList.Items) != len(pongList.Items) {
		log.Errorf("PingList doesn't match PongList")
	}

	if len(pingList.Items) != len(nodeList.Items) || len(pongList.Items) != len(nodeList.Items) {
		log.Errorf("PingList doesn't match NodeList or PongList doesn't match NodeList")

		/*
			pingDaemonSet, err := client.DaemonSets("yce").Get("ping")
			checkError("Get Ping DaemonSet", err)
			pongDaemonSet, err := client.DaemonSets("yce").Get("pong")
			checkError("Get Pong DaemonSet", err)

			client.DaemonSets("yce").Delete("ping", &v1.DeleteOptions{})
			client.DaemonSets("yce").Delete("pong", &v1.DeleteOptions{})

			time.Sleep(5 * time.Second)

			client.DaemonSets("yce").Create(pingDaemonSet)
			client.DaemonSets("yce").Create(pongDaemonSet)

			for _, node := range nodeList.Items {
				for _, addr := range node.Status.Addresses {
					if addr.Type == "LegacyHostIP" {
						c.Delete("pod:" + addr.Address + ":pod")
						c.Delete("pod:" + addr.Address + ":node")
						c.Delete("pod:" + addr.Address + ":service")
					}
				}
			}
		*/
	}
}

func GetPodList(client *kube.Clientset) {
	pingInfoList := make([]types.PodInfo, 0)
	pongInfoList := make([]types.PodInfo, 0)
	podList, err := client.Pods("yce").List(v1.ListOptions{})
	pingList := new(v1.PodList)
	pongList := new(v1.PodList)
	for _, pod := range podList.Items {
		if strings.Contains(pod.Name, "ping-") {
			pingList.Items = append(pingList.Items, pod)
		}

		if strings.Contains(pod.Name, "pong-") {
			pongList.Items = append(pongList.Items, pod)
		}
	}
	log.Infof("GetPodList PongPodList %d", len(pongList.Items))
	log.Infof("GetPodList PingPoddList %d", len(pingList.Items))
	utils.CheckError("GetPodList List pods", err)

	for _, pod := range pingList.Items {
		p := new(types.PodInfo)
		p.Name = pod.Name
		p.PodIP = pod.Status.PodIP
		p.HostIP = pod.Status.HostIP
		p.Status = string(pod.Status.Phase)
		pingInfoList = append(pingInfoList, *p)
	}

	for _, pod := range pongList.Items {
		p := new(types.PodInfo)
		p.Name = pod.Name
		p.PodIP = pod.Status.PodIP
		p.HostIP = pod.Status.HostIP
		p.Status = string(pod.Status.Phase)
		pongInfoList = append(pongInfoList, *p)
	}

	rawPingData, err := json.Marshal(pingInfoList)
	rawPongData, err := json.Marshal(pongInfoList)
	utils.CheckError("GetPodList JsonMarshal", err)
	pingdata := string(rawPingData)
	pongdata := string(rawPongData)
	cache.SetWithExpire("pinglist", pingdata, strconv.Itoa(constant.RESOURCE_TIMER))
	cache.SetWithExpire("ponglist", pongdata, strconv.Itoa(constant.RESOURCE_TIMER))
}

func GetNodeList(client *kube.Clientset) {
	nodeInfoList := make([]types.NodeInfo, 0)
	list, err := client.Nodes().List(v1.ListOptions{})
	log.Infof("GetNodeList NodeList %d", len(list.Items))
	utils.CheckError("GetNodeList Nodes List", err)

	for _, node := range list.Items {
		n := new(types.NodeInfo)
		n.Name = node.Name
		for _, addr := range node.Status.Addresses {
			if addr.Type == "LegacyHostIP" {
				n.HostIP = addr.Address
			}
		}
		n.Status = string(node.Status.Phase)
		nodeInfoList = append(nodeInfoList, *n)
	}

	rawData, _ := json.Marshal(nodeInfoList)
	data := string(rawData)
	cache.SetWithExpire("nodelist", data, strconv.Itoa(constant.RESOURCE_TIMER))
}

func GetSvc(client *kube.Clientset) {
	serviceInfo := types.ServiceInfo{}
	service, err := client.Services("yce").Get("ping-svc")
	utils.CheckError("GetSvcList Service Get", err)

	serviceInfo.Name = service.Name
	serviceInfo.Type = string(service.Spec.Type)
	serviceInfo.ClusterIP = service.Spec.ClusterIP
	for _, p := range service.Spec.Ports {
		port := types.Port{
			Port:     p.Port,
			NodePort: p.NodePort,
		}
		serviceInfo.Ports = append(serviceInfo.Ports, port)
	}

	data, _ := json.Marshal(serviceInfo)
	cache.SetWithExpire("service", string(data), strconv.Itoa(constant.RESOURCE_TIMER))
}
