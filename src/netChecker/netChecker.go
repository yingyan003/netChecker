package main

import (
	"net/http"
	"netChecker/controller"
	"common/utils"
	"common/constant"
)

func init(){
	utils.SetLogLevel(utils.LoadEnvVarInt(constant.ENV_LOG_LEVEL, constant.LOG_LEVEL_ERROR))
	controller.Init()
}

func main() {
	go controller.GetKubeResAndPublish()
	go controller.ReceiveAndSaveData()

	fs := http.FileServer(http.Dir("frontend/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/report", controller.All)
	http.HandleFunc("/report/true", controller.True)
	http.HandleFunc("/report/false", controller.False)
	http.HandleFunc("/simple/true",controller.TrueSimple)
	http.HandleFunc("/simple/false",controller.FalseSimple)
	http.HandleFunc("/nodes",controller.NodeToNode)

	http.ListenAndServe(":8080", nil)
}



