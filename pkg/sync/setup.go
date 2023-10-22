package sync

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"

	qp "github.com/quic-s/quics-protocol"
)

var (
	QPClient    *qp.QP
	Conn        *qp.Connection
	Watcher     *fsnotify.Watcher
	PSMut       map[byte]*sync.Mutex
	PSMutModNum uint8 = 64
)

func init() {

	InitQPClient()

	PSMut = make(map[byte]*sync.Mutex)
	for i := uint8(0); i < PSMutModNum; i++ {
		PSMut[i] = &sync.Mutex{}
	}

	MustSyncMain()
	ForceSyncMain()
	FullScanMain()
	NeedContentMain()

}
func InitQPClient() {
	//newClient, err := qp.New(qp.LOG_LEVEL_INFO)
	newClient, err := qp.New(qp.LOG_LEVEL_ERROR)
	if err != nil {
		panic(err)
	}
	QPClient = newClient
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
