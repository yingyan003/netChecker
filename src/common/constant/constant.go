package constant

//common
const(
	LOG_LEVEL_ERROR = 5

	APISERVER = "k8s apiserver所在的节点IP:apiserver暴露的端口"
	ENV_APISERVER = "APISERVER"
	REDISHOST = "redis-master-svc.zxy:6379"
	ENV_REDISHOST = "REDISHOST"

	NODEPORT="32079"
	ENV_NODEPORT="NODEPORT"
)

//netguarder
const (
	RESOURCE_TIMER   = 59
	ENV_RESOURCE_TIMER   = "TIMER"

	ALL = iota
	TRUE
	FALSE
	TRUE_SIMPLE
	FALSE_SIMPLE
)

//ping
const (
	DEBUG_FAIL="源pod网络故障"

	TELNET_TIMER   = 60
	ENV_TELNET_TIMER   = "TIMER"
)
