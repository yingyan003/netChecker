package main

import (
	"net/http"
	"netguarder/controller"
	"common/sysinit"
	"common/utils"
	"common/constant"
)

func init(){
	utils.SetLogLevel(constant.LOG_LEVEL_ERROR)
	sysinit.InitConfig()
	controller.Init()
}

func main() {
	go controller.GetKubeResToRedis()

	fs := http.FileServer(http.Dir("frontend/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/report", controller.All)
	http.HandleFunc("/report/true", controller.True)
	http.HandleFunc("/report/false", controller.False)
	http.HandleFunc("/simple/true",controller.TrueSimple)
	http.HandleFunc("/simple/false",controller.FalseSimple)
	http.ListenAndServe(":8080", nil)
}



