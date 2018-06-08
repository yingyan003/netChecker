package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"common/utils"
	"common/constant"
	mylog "github.com/maxwell92/gokits/log"
	"io/ioutil"
	"strings"
)

var log = mylog.Log
var KubeClient *kubernetes.Clientset

func NewKubeClient() error {
	var err error
	//auth=0表示连接k8s不需要认证，默认为0
	auth := utils.LoadEnvVarInt(constant.ENV_AUTH, constant.AUTH)
	//非tls连接
	if auth == 0 {
		err = getKubeClinet()
	} else {
		err = getTLSKubeClinet()
	}
	return err
}

func getKubeClinet() error {
	var err error
	config := new(rest.Config)
	config.Host = utils.LoadEnvVar(constant.ENV_APISERVER, constant.APISERVER)
	KubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("getKubeClinet failed: error=%s, apiServer=%s\n", err, config.Host)
	}
	return err
}

func getTLSKubeClinet() error {
	var err error
	config := new(rest.Config)
	config.Host = utils.LoadEnvVar(constant.ENV_APISERVER, constant.APISERVER)
	caPath := utils.LoadEnvVar(constant.ENV_CAPATH, constant.CAPATH)
	certPath := utils.LoadEnvVar(constant.ENV_CERTPATH, constant.CERTPATH)
	keyPath := utils.LoadEnvVar(constant.ENV_KEYPATH, constant.KEYPATH)
	//tls连接
	if !strings.EqualFold(caPath, "") && !strings.EqualFold(certPath, "") && !strings.EqualFold(keyPath, "") {
		//从指定的路径下读文件
		config.CAData, err = ioutil.ReadFile(caPath)
		if err != nil {
			log.Errorf("getTLSKubeClinet failed: Read ca File error. err=%v", err)
			return err
		}
		config.CertData, err = ioutil.ReadFile(certPath)
		if err != nil {
			log.Errorf("getTLSKubeClinet failed: Read cert File error. err=%v", err)
			return err
		}
		config.KeyData, err = ioutil.ReadFile(keyPath)
		if err != nil {
			log.Errorf("getTLSKubeClinet failed: Read key File error. err=%v", err)
			return err
		}
	} else {
		//如果不指定文件路径，默认从当前路径下读文件
		config.CAData, err = ioutil.ReadFile(constant.CAPATH)
		if err != nil {
			log.Errorf("getTLSKubeClinet failed: Read ca File error. err=%v", err)
			return err
		}
		config.CertData, err = ioutil.ReadFile(constant.CERTPATH)
		if err != nil {
			log.Errorf("getTLSKubeClinet failed: Read cert File error. err=%v", err)
			return err
		}
		config.KeyData, err = ioutil.ReadFile(constant.KEYPATH)
		if err != nil {
			log.Errorf("getTLSKubeClinet failed: Read key File error. err=%v", err)
			return err
		}
	}

	KubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("GetTLSKubeClient failed: error=%s, apiServer=%s\n", err, config.Host)
	}
	return err
}

//func GetKubeClient() *kubernetes.Clientset {
//	Config := new(rest.Config)
//	Config.Host = utils.LoadEnvVar(constant.ENV_KUBE_HOST, constant.KUBE_HOST)
//	Config.CAData = []byte(utils.LoadEnvVar(constant.ENV_CA_DATA, constant.CA_DATA))
//	Config.CertData = []byte(utils.LoadEnvVar(constant.ENV_CERT_DATA, constant.CERT_DATA))
//	Config.KeyData = []byte(utils.LoadEnvVar(constant.ENV_KET_DATA, constant.KET_DATA))
//
//	kubeClient, err := kubernetes.NewForConfig(Config)
//	if err != nil {
//		log.Errorf("GetTLSKubeClient failed: error=%s, apiServer=%s\n", err, Config.Host)
//	}
//	return kubeClient
//}
