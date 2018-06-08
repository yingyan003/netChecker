package sysinit

import (
	"common/constant"
	"common/utils"
	"common/types"
)

var conf *types.Config

func InitConfig(){
	conf=&types.Config{
		Apiserver: utils.LoadEnvVar(constant.ENV_APISERVER,constant.APISERVER),
		RedisHost: utils.LoadEnvVar(constant.ENV_REDISHOST,constant.REDISHOST),
		NodePort:utils.LoadEnvVar(constant.ENV_NODEPORT,constant.NODEPORT),
	}
}

func GetConfig()*types.Config{
	return conf
}
