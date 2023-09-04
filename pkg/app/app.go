package app

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-s/quics-client/pkg/badger"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
)

func Reboot() {
	log.Println("\n\trebooting ...")

	str, err := os.Executable()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	newProcess := exec.Command(str)
	newProcess.Stdout = os.Stdout
	newProcess.Stderr = os.Stderr
	newProcess.Stdin = os.Stdin
	newProcess.Env = os.Environ()

	err = newProcess.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)

}

func RestServerStart() {

	defer badger.CloseDB()

	handler := setupHandler()
	qconf := quic.Config{}
	server := http3.Server{
		Handler:    handler,
		QuicConfig: &qconf,

		Addr: "0.0.0.0:" + viper.GetViperEnvVariables("REST_SERVER_PORT"),
	}

	quicsDir := utils.GetQuicsDirPath()
	certDir := filepath.Join(quicsDir, viper.GetViperEnvVariables("QUICS_CLI_CERT_NAME"))
	keyDir := filepath.Join(quicsDir, viper.GetViperEnvVariables("QUICS_CLI_KEY_NAME"))

	// load the certificate and the key from the files
	_, err := tls.LoadX509KeyPair(certDir, keyDir)
	if err != nil {
		utils.CertFile()
	}

	err = server.ListenAndServeTLS(certDir, keyDir)
	if err != nil {
		log.Fatal("Client Server Error : ", err)
	}

}

func setupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%#v\n", r)
		w.Write([]byte("hello, Quics Client here"))
	})

	// TODO : add handler for each api
	// mux.HandleFunc("/api/v1/settings/server", func(w http.ResponseWriter, r *http.Request) {

	// }

	return mux
}
