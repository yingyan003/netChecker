package main

import (
	"common/utils"
	"net/http"
	"common/constant"
	"ping/controller"
)

func init(){
	utils.SetLogLevel(utils.LoadEnvVarInt(constant.ENV_LOG_LEVEL, constant.LOG_LEVEL_ERROR))
	controller.Init()
}

func main() {
	go controller.NetCheck()
	http.ListenAndServe(":8080", nil)
}