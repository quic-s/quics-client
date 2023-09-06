package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	qviper "github.com/quic-s/quics-client/pkg/viper"
)

type RestClient struct {
	qconf        *quic.Config
	roundTripper *http3.RoundTripper
	hclient      *http.Client
}

func NewRestClient() *RestClient {
	restClient := &RestClient{
		qconf: &quic.Config{
			KeepAlivePeriod: 60,
		},
	}

	restClient.roundTripper = &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		QuicConfig: restClient.qconf,
	}

	restClient.hclient = &http.Client{
		Transport: restClient.roundTripper,
	}
	return restClient
}

func (r *RestClient) GetRequest(path string) *bytes.Buffer {

	url := "https://localhost:" + qviper.GetViperEnvVariables("REST_SERVER_PORT") + path

	// hclient := NewRestClient()
	//getReqest, err := http.NewRequest("GET", urls, nil)

	rsp, err := r.hclient.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("quics-client :rsp : ", rsp)

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("quics-client :Response body : ", body.String())
	return body

}

func (r *RestClient) PostRequest(path string, contentType string, content []byte) *bytes.Buffer {

	url := "https://localhost:" + qviper.GetViperEnvVariables("REST_SERVER_PORT") + path

	contentReader := bytes.NewReader(content)
	rsp, err := r.hclient.Post(url, contentType, contentReader)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("quics-client :rsp : ", rsp)

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("quics-client :Response body : ", body.String())
	return body
}

func (r *RestClient) Close() error {
	r.hclient.CloseIdleConnections()
	err := r.roundTripper.Close()
	if err != nil {
		return err
	}
	return nil
}
