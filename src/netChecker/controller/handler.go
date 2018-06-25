package controller

import (
	"fmt"
	"net/http"
	"common/constant"
)

func True(w http.ResponseWriter, r *http.Request) {
	data := cache.Get(constant.TRUE_DATA)
	log.Infof("handler True http requst success")
	fmt.Fprintf(w, "%s\n", data)
}

func False(w http.ResponseWriter, r *http.Request) {
	data := cache.Get(constant.FALSE_DATA)
	log.Infof("handler False http requst success")
	fmt.Fprintf(w, "%s\n", data)
}

func All(w http.ResponseWriter, r *http.Request) {
	data := cache.Get(constant.ALL_DATA)
	log.Infof("handler All http requst success")
	fmt.Fprintf(w, "%s\n", data)
}

func TrueSimple(w http.ResponseWriter, r *http.Request) {
	data := cache.Get(constant.SIMPLE_TRUE_DATA)
	log.Infof("handler TrueSimple http requst success")
	fmt.Fprintf(w, "%s\n", data)
}

func FalseSimple(w http.ResponseWriter, r *http.Request) {
	data := cache.Get(constant.SIMPLE_FALSE_DATA)
	log.Infof("handler FalseSimple http requst success")
	fmt.Fprintf(w, "%s\n", data)
}

func NodeToNode(w http.ResponseWriter, r *http.Request) {
	data := cache.Get(constant.NOTE_TO_NODE_DATA)
	log.Infof("handler NodeToNode http requst success")
	fmt.Fprintf(w, "%s\n", data)
}

