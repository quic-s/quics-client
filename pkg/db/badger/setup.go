package badger

import (
	"log"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics-client/pkg/utils"
)

// Declare a global variable for the DB.
var badgerdb *badger.DB
var mutex sync.Mutex

// Define a function to open the DB.
func OpenDB() {
	// Open the Badger database located in the ./badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions(utils.GetQuicsDirPath() + "/badger")
	db, err := badger.Open(opts)
	badgerdb = db
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	if err := badgerdb.Close(); err != nil {
		log.Fatal(err)
	}
}
