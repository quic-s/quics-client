package http

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-s/quics-client/pkg/utils"
)

func RestServerStart() {
	go func() {
		// log.Println(http.ListenAndServe(":6060", nil))
	}()
	handler := setupHandler()
	qconf := quic.Config{}
	server := http3.Server{
		Handler:    handler,
		QuicConfig: &qconf,
		Addr:       ":6121",
	}
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	quicsDir := filepath.Join(tempDir, "quics")
	certDir := filepath.Join(quicsDir, "cert-quics-cli.pem")
	keyDir := filepath.Join(quicsDir, "key-quics-cli.pem")

	// load the certificate and the key from the files
	_, err = tls.LoadX509KeyPair(certDir, keyDir)
	if err != nil {
		utils.CertFile()
	}
	// fmt.Println("cert : ", cert.)

	log.Println(server.ListenAndServeTLS(certDir, keyDir))

}

func setupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%#v\n", r)
		w.Write([]byte("hello, world"))
	})

	return mux
}
