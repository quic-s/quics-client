package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {

	insecure := true

	flag.Parse()
	urls := []string{"https://localhost:6121"}

	var qconf quic.Config

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{

			InsecureSkipVerify: insecure,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, addr := range urls {
		log.Println("addr : ", addr)

		go func(addr string) {
			rsp, err := hclient.Get(addr)
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

			wg.Done()
		}(addr)
	}
	wg.Wait()
}
