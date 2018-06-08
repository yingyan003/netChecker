package https

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

type HttpsClient struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	pool   *x509.CertPool
	Client *http.Client
}

func NewHttpsClient(host, port, cert string) *HttpsClient {

	https := &HttpsClient{
		Host: host,
		Port: port,
		pool: x509.NewCertPool(),
	}

	caCrt, err := ioutil.ReadFile(cert)
	if err != nil {
		log.Fatalf("NewHttpsClient ReadFile err: file=%s, err=%s", cert, err)
	}

	https.pool.AppendCertsFromPEM(caCrt)

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: https.pool},
		DisableCompression: true,
	}

	https.Client = &http.Client{Transport: tr}

	return https
}
