package constant

//common
const (
	ENV_LOG_LEVEL   = "LOG_LEVEL"
	LOG_LEVEL_ERROR = 5

	//K8S
	APISERVER = "10.151.160.11:8080"
	CAPATH    = "ca.crt"
	CERTPATH  = "client.crt"
	KEYPATH   = "client.key"

	ENV_APISERVER = "APISERVER"
	ENV_CAPATH    = "CAPATH"
	ENV_CERTPATH  = "CERTPATH"
	ENV_KEYPATH   = "KEYPATH"

	ENV_AUTH = "AUTH"
	AUTH     = 0
)

//netguarder
const (
	GET_RESOURCE_TICKER     = 60
	ENV_GET_RESOURCE_TICKER = "GET_RESOURCE_TICKER"

	ENV_REDIS_EXPIRE = "REDIS_EXPIRE"
	REDIS_EXPIRE = "60"

	// 每轮数据采集的超时时间
	// （1 + ping网络检测超时时间 + n） < 建议超时时间(s) < （1 + ping网络检测超时时间 + len(k8s work node) + x）
	//说明：
	//1：每轮从prepare发布成功到receive接受完该轮数据的时间差，基本为1s
	//2：每个ping(每个节点只部署一个名字为ping的pod)网络检测一次失败时的总超时时间：1s<=超时<2s：
	// 	 因为ping->pong和ping->baidu，超时时间各设为1s(两个goroutine同时检测)
	//3：n（每个ping执行时间，pub/sub时间都会存在一些时延，最好预留个几秒中等待）
	//4：x（预留集群扩容的节点数，可选，作集群添加节点时用）
	//5：当然，超时时间可以比上面的范围大，但要小于每轮从k8s获取资源列表的时间间隔（60s）小。
	//   可如此一来，前端获取到的数据即时性就很差了
	ENV_RECEIVE_TIMEOUT = "RECEIVE_TIMEOUT"
	RECEIVE_TIMEOUT     = 10

	ALL          = iota
	TRUE
	FALSE
	TRUE_SIMPLE
	FALSE_SIMPLE
	NOTE_TO_NODE

	ALL_DATA          = "netguard-all"
	FALSE_DATA        = "netguard-false"
	TRUE_DATA         = "netguard-true"
	SIMPLE_FALSE_DATA = "netguard-simpleFalse"
	SIMPLE_TRUE_DATA  = "netguard-simpleTrue"
	NOTE_TO_NODE_DATA = "netguard-nodeToNode"
)

//ping
const (
	DEBUG_FAIL = "源pod网络出不去"

	NODEPORT     = 32089
	ENV_NODEPORT = "NODEPORT"

	TELNET_TIMER     = 60
	ENV_TELNET_TIMER = "TIMER"
)

//redis
const (
	//REDIS
	ENV_REDISHOST = "REDISHOST"
        //短域名的方式访问redis服务，kubeDNS会自动做服务发现
	REDISHOST     = "redis-master-svc.zxy:6379"

	NETGUARD_MAX_IDLE = 5
	NETCHECK_MAX_IDLE = 5
	MAX_ACTIVE        = 100
	IDLE_TIMEOUT      = 180 //180

	CHANNEL_POD  = "POD"
	CHANNEL_NODE = "NODE"
	CHANNEL_SVC  = "SVC"

	CHANNEL_REPORT_POD  = "REPORT_POD"
	CHANNEL_REPORT_NODE = "REPORT_NODE"
	CHANNEL_REPORT_SVC  = "REPORT_SVC"
)
