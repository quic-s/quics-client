package badger

import (
	"log"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics-client/pkg/utils"
)

// Declare a global variable for the DB.
type Badger struct {
	BadgerDB *badger.DB
	Mutex    sync.Mutex `wire:"-"` // ignore this field
}

func NewBadger() *Badger {
	// Open the Badger database located in the ./badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions(utils.GetQuicsDirPath() + "/badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return &Badger{
		BadgerDB: db,
	}

}

func (bg *Badger) CloseDB() {
	if err := bg.BadgerDB.Close(); err != nil {
		log.Fatal(err)
	}
}
