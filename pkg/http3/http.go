package http3

import (
	"crypto/tls"
	"log"
	"path/filepath"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-s/quics-client/pkg/badger"
	"github.com/quic-s/quics-client/pkg/connection"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
)

func RestServerStart() {

	log.Println("quics-client : starting port " + viper.GetViperEnvVariables("REST_SERVER_PORT"))
	badger.OpenDB()
	connection.InitWatcher()

	rootdirlist := utils.GetRootDirs()
	for _, value := range rootdirlist {
		if value != "" {
			connection.DirWatchAdd(value)

		}

	}
	connection.DirWatchStart()
	defer connection.WatchStop()
	defer badger.CloseDB()

	handler := SetupHandler()
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
