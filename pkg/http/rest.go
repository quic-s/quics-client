package http

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
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

	quicsDir := utils.GetDirPath()
	certDir := filepath.Join(quicsDir, utils.GetViperEnvVariables("QUICS_CLI_CERT_NAME"))
	keyDir := filepath.Join(quicsDir, utils.GetViperEnvVariables("QUICS_CLI_KEY_NAME"))

	// load the certificate and the key from the files
	_, err := tls.LoadX509KeyPair(certDir, keyDir)
	if err != nil {
		utils.CertFile()
	}

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
