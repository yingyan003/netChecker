package utils

import (
	"os"
	"strconv"
	"net/http"
)

func LoadEnvVar(key, value string)string{
	var v string
	if v=os.Getenv(key);v==""{
		v=value
	}
	return v
}

func LoadEnvVarInt(key string, value int)int{
	i,_:=strconv.Atoi(LoadEnvVar(key,strconv.Itoa(value)))
	return i
}

func CheckError(msg string, err error) {
	if err != nil {
		log.Errorf("Msg=%s, Error=%s", msg, err)
	}
}

//解决浏览器跨域
func CorsHandler(w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Set("Access-Control-Allow-Headers", "POST, GET, OPTIONS, PUT, DELETE") //header的类型
	w.Header().Set("Access-Control-Allow-Headers",
		"Origin, Authorization, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
}
