package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	qviper "github.com/quic-s/quics-client/pkg/viper"
)

func GetRestClient(path string) *bytes.Buffer {

	flag.Parse()
	url := "https://localhost:" + qviper.GetViperEnvVariables("REST_SERVER_PORT") + path

	var qconf quic.Config

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{

			InsecureSkipVerify: true,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	//getReqest, err := http.NewRequest("GET", urls, nil)

	rsp, err := hclient.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("rsp : ", rsp)

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response body : ", body.String())
	return body

}

func PostRestClient(path string, contentType string, content io.Reader) *bytes.Buffer {

	flag.Parse()
	url := "https://localhost:" + qviper.GetViperEnvVariables("REST_SERVER_PORT") + path

	var qconf quic.Config

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{

			InsecureSkipVerify: true,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	rsp, err := hclient.Post(url, contentType, content)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("rsp : ", rsp)

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response body : ", body.String())
	return body

}
