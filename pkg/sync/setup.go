package sync

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/google/wire"

	qp "github.com/quic-s/quics-protocol"
)

const PSMutModNum uint8 = 16

var ServiceSet = wire.NewSet(
	NewMutex,
	NewQPClient,
	NewWatcher,
	wire.Struct(new(Service), "*"),
)

type Service struct {
	Conn     *qp.Connection `wire:"-"` // ignore this field when inject
	QPClient *qp.QP
	Watcher  *fsnotify.Watcher
	PSMut    map[byte]*sync.Mutex
	db       Repository
	//TODO http
}

func NewService(conn *qp.Connection, qpClient *qp.QP, watcher *fsnotify.Watcher, pSMut map[byte]*sync.Mutex) *Service {
	return &Service{
		Conn:     conn,
		QPClient: qpClient,
		Watcher:  watcher,
		PSMut:    pSMut,
	}
}

func NewMutex() map[byte]*sync.Mutex {

	pSMut := make(map[byte]*sync.Mutex)
	for i := uint8(0); i < PSMutModNum; i++ {
		pSMut[i] = &sync.Mutex{}
	}

	// TODO when initiate
	MustSyncMain()
	ForceSyncMain()
	FullScanMain()
	NeedContentMain()

	return pSMut
}

func NewQPClient() *qp.QP {
	//newClient, err := qp.New(qp.LOG_LEVEL_INFO)
	newClient, err := qp.New(qp.LOG_LEVEL_ERROR)
	if err != nil {
		panic(err)
	}
	return newClient
}

func NewWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	return watcher
}

func (sv *Service) closeConnect() {
	sv.Conn.Close()
	sv.Conn = nil
}
