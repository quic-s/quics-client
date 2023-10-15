package sync

import (
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/quic-s/quics-client/pkg/viper"
	qp "github.com/quic-s/quics-protocol"
)

var (
	QPClient *qp.QP
	Conn     *qp.Connection
	Watcher  *fsnotify.Watcher
	PSMut    map[byte]*sync.Mutex
)

func init() {

	err := error(nil)
	QPClient, err = qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		panic(err)
	}

	PSMut = make(map[byte]*sync.Mutex)
	for i := uint8(0); i < 16; i++ {
		PSMut[i] = &sync.Mutex{}
	}

	MustSyncMain()
	ForceSyncMain()

}
func InitWatcher() {
	// Create a new watcher.
	err := error(nil)
	Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
}

func CloseConnect() {
	Conn.Close()
	Conn = nil
}

func ReConnect() {
	p, err := strconv.Atoi(viper.GetViperEnvVariables("QUICS_SERVER_PORT"))
	if err != nil {
		panic(err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-s"},
	}
	Conn, err = QPClient.Dial(&net.UDPAddr{IP: net.ParseIP(viper.GetViperEnvVariables("QUICS_SERVER_IP")), Port: p}, tlsConf)
	if err != nil {
		panic(err)
	}
}
