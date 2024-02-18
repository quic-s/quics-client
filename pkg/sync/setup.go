package sync

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/google/wire"
)

const PSMutModNum uint8 = 16

var ServiceSet = wire.NewSet(
// TODO
)

type Service struct {
	Watcher *fsnotify.Watcher
	PSMut   map[byte]*sync.Mutex
	db      Repository
	qclient QPPort
	//TODO http
}

//TODO ServiceProvider

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

func NewWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	return watcher
}
