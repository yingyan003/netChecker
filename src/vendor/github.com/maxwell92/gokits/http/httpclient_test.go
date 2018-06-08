package http

import (
	"fmt"
	"testing"
)

func Test_HttpClient_New(t *testing.T) {
	client := NewHttpClient("http://master", "8080")
	if client == nil {
		t.Error("HttpClient_New didn't pass")

	} else {
		fmt.Printf("%v\n", client)
		t.Log("test 1 passed")
	}
}

func Test_HttpClient_Get(t *testing.T) {
	client := NewHttpClient("http://master", "8080")
	resp, err := client.Get("http://master:8080/apis/extensions/v1beta1/namespaces/default/deployments")
	if err != nil {
		t.Error("Get pods error")
	} else {
		fmt.Println(string(resp))
	}
}

/*
func Test_HttpClient_Post(t *testing.T) {

	tempDeploy := new(appdeploy.BasicDeployment)

	defaultlabels := make(map[string]string, 1)
	defaultlabels["appname"] = "nginx-test"

	tempDeploy.ApiVersion = "extensions/v1beta1"
	tempDeploy.Kind = "Deployment"

	tempDeploy.Metadata.Name = "nginx-test"
	tempDeploy.Spec.Replicas = 1
	tempDeploy.Spec.Template.Metadata.Labels = defaultlabels
	tempDeploy.Spec.Template.Spec.Containers = make([]appdeploy.ContainersSTSDL, 1)
	tempDeploy.Spec.Template.Spec.Containers[0].Name = "nginx-test"
	tempDeploy.Spec.Template.Spec.Containers[0].Image = "nginx:1.7.9"

	var result []byte
	result, _ = json.MarshalIndent(tempDeploy, "", "    ")

	client := NewHttpClient("http://master", "8080")
	resp, err := client.Post("http://master:8080/apis/extensions/v1beta1/namespaces/default/deployments", strings.NewReader(string(result)))

	if err != nil {
		t.Error("Post pods error")
	} else {
		fmt.Println(string(resp))
	}

}
*/

func Test_HttpClient_Delete(t *testing.T) {
	client := NewHttpClient("http://master", "8080")
	resp, err := client.Delete("http://master:8080/apis/extensions/v1beta1/namespaces/default/deployments/nginx-test")
	if err != nil {
		t.Error("Delete pods error")
	} else {
		fmt.Println(string(resp))
	}
}
