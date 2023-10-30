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

func (r *RestClient) GetRequest(path string) (*bytes.Buffer, error) {

	url := "https://localhost:" + qviper.GetViperEnvVariables("REST_SERVER_PORT") + path

	rsp, err := r.hclient.Get(url)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil

}

func (r *RestClient) PostRequest(path string, contentType string, content []byte) (*bytes.Buffer, error) {

	url := "https://localhost:" + qviper.GetViperEnvVariables("REST_SERVER_PORT") + path
	log.Println("log : " + string(content))
	contentReader := bytes.NewReader(content)
	rsp, err := r.hclient.Post(url, contentType, contentReader)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, rsp.ContentLength)
	_, err = rsp.Body.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

func (r *RestClient) Close() error {
	r.hclient.CloseIdleConnections()
	err := r.roundTripper.Close()
	if err != nil {
		return err
	}
	return nil
}
