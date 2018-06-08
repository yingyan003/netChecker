package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

//该版本只能连接不需要证书访问的k8s集群
func GetKubeClient(apiServer string) *kubernetes.Clientset {
	cfg := &rest.Config{
		Host: "http://" + apiServer,
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Printf("GetKubeClient Error: error=%s, apiServer=%s", err, apiServer)
		return nil
	}
	return clientset
}
