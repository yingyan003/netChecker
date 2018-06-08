package http

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpClient struct {
	Client *http.Client
	Host   string
	Port   string
}

var instance *HttpClient

func HttpInstance() *HttpClient {
	return instance
}

func NewHttpClient(host, port string) *HttpClient {
	tr := new(http.Transport)
	tr.TLSClientConfig = new(tls.Config)
	tr.TLSClientConfig.InsecureSkipVerify = true

	client := new(http.Client)
	client.Transport = tr

	instance = new(HttpClient)
	instance.Client = client
	instance.Host = host
	instance.Port = port

	return instance
}

func (c *HttpClient) Get(url string) (statusCode int, response []byte, err error) {
	resp, err := c.Client.Get(url)
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		return 0, nil, err
	}

	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		return 0, nil, err
	}

	return resp.StatusCode, response, nil
}
func (c *HttpClient) Post(url string, body io.Reader) (response []byte, err error) {
	resp, err := c.Client.Post(url, "application/json;charset=utf-8", body)

	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		return nil, err
	}

	return response, nil
}

//TODO: add Http Put Method
func (c *HttpClient) Put() {

}

func (c *HttpClient) Delete(url string) (response []byte, err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		return nil, err
	}

	return response, nil
}
