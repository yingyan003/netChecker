package main

import (
	"common/utils"
	"net/http"
	"common/sysinit"
	"common/constant"
	"ping/controller"
)

func init(){
	utils.SetLogLevel(constant.LOG_LEVEL_ERROR)
	sysinit.InitConfig()
	controller.Init()
}

func main() {
	utils.SetLogLevel(constant.LOG_LEVEL_ERROR)
	sysinit.InitConfig()

	go controller.NetCheck()

	http.ListenAndServe(":8080", nil)
}