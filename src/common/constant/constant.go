package constant

//common
const (
	ENV_LOG_LEVEL   = "LOG_LEVEL"
	LOG_LEVEL_ERROR = 5

	//K8S
	APISERVER = "k8s任意一个master节点的apiserver所在的节点ip:apiserver开放的端口"
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

	REDIS_EXPIRE = 60

	ENV_RECEIVE_TIMEOUT = "RECEIVE_TIMEOUT"
	//建议超时时间= 1 + 2 * len(k8s work node) + x
	//说明：
	//1：每轮从prepare发布成功到receive接受完该轮数据的时间差，基本为1s
	//2：每个ping(每个节点只部署一个名字为ping的pod)网络检测一次失败时的总超时时间：ping->pong和ping->baidu，超时时间各为1s
	//3：x(x可选，作集群添加节点时用)
	//此处是： 1 + 2*8 +3（集群可扩展3个节点）
	RECEIVE_TIMEOUT = 20

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

	NODEPORT     = 32079
	ENV_NODEPORT = "NODEPORT"

	TELNET_TIMER     = 60
	ENV_TELNET_TIMER = "TIMER"
)

//redis
const (
	//REDIS
	ENV_REDISHOST = "REDISHOST"
	//1. 这里的redis用的是deployment的形式部署的。并建了一个对应的service
	//2. k8s集群中的pod通过kubedns的形式访问该redis服务。
	//只需将redis的host设置为"serviceName.namespace"，
	//kubedns自动会解析该短域名，并访问Redis服务。
	REDISHOST = "redis-master-svc.zxy:6379"

	NETGUARD_MAX_IDLE = 5
	NETCHECK_MAX_IDLE = 5
	MAX_ACTIVE        = 100
	IDLE_TIMEOUT      = 0 //180

	CHANNEL_POD  = "POD"
	CHANNEL_NODE = "NODE"
	CHANNEL_SVC  = "SVC"

	CHANNEL_REPORT_POD  = "REPORT_POD"
	CHANNEL_REPORT_NODE = "REPORT_NODE"
	CHANNEL_REPORT_SVC  = "REPORT_SVC"
)
