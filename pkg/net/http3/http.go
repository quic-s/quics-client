package http3

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/sync"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
)

func Http3ServerProvider() {

}
func RestServerStart(port string) {

	fmt.Println("\t-----------------------------------------\n")
	fmt.Println("\t           quics-client start\n")
	fmt.Println("\t-----------------------------------------\n")

	log.Println("quics-client : starting port " + viper.GetViperEnvVariables("REST_SERVER_PORT"))

	//TODO bagder injection
	badger.OpenDB()

	sync.InitWatcher()
	sync.DirWatchStart()
	rootdirlist := badger.GetRootDirList()
	for _, rootdir := range rootdirlist {
		sync.DirWatchAdd(rootdir.Path)
	}

	defer sync.WatchStop()
	defer badger.CloseDB()

	handler := SetupHandler()
	qconf := quic.Config{}

	if port == "" {
		port = viper.GetViperEnvVariables("REST_SERVER_PORT")
	}

	server := http3.Server{
		Handler:    handler,
		QuicConfig: &qconf,
		Addr:       "0.0.0.0:" + port,
	}

	eport := viper.GetViperEnvVariables("REST_ENTRY_SERVER_PORT")
	httpserver := http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:" + eport,
	}

	quicsDir := utils.GetQuicsDirPath()
	certDir := filepath.Join(quicsDir, viper.GetViperEnvVariables("QUICS_CLI_CERT_NAME"))
	keyDir := filepath.Join(quicsDir, viper.GetViperEnvVariables("QUICS_CLI_KEY_NAME"))

	// load the certificate and the key from the files
	_, err := tls.LoadX509KeyPair(certDir, keyDir)
	if err != nil {
		utils.CertFile()
	}

	// Check if server is already connected before
	// if yes, then reconnect to server
	if viper.GetViperEnvVariables("QUICS_SERVER_IP") != "" && viper.GetViperEnvVariables("QUICS_SERVER_PORT") != "" && viper.GetViperEnvVariables("QUICS_SERVER_PASSWORD") != "" {

		log.Println("quics-client : [INIT] Try reconnecting to server")
		go func() {
			err := sync.ClientRegistration(viper.GetViperEnvVariables("QUICS_SERVER_PASSWORD"), viper.GetViperEnvVariables("QUICS_SERVER_IP"), viper.GetViperEnvVariables("QUICS_SERVER_PORT"))
			if err != nil {
				log.Println("quics-client : [INIT] Cannot reconnect to server. !! PLEASE MAKE CONNECTION FIRST !!")
			} else {
				log.Println("quics-client : [INIT] Reconnected to server. !! LET'S GET STARTED !!")
			}

		}()
	}
	go func() {
		err := httpserver.ListenAndServeTLS(certDir, keyDir)
		if err != nil {
			log.Fatal("Client HTTP Server Error : ", err)
		}
	}()

	// start to listen
	err = server.ListenAndServeTLS(certDir, keyDir)
	if err != nil {
		log.Fatal("Client Server Error : ", err)
	}

}
