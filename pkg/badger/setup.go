package badger

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics-client/pkg/utils"
)

// Declare a global variable for the DB.
var Badgerdb *badger.DB

const (
	META string = "META"
)

// Define a function to open the DB.
func init() {
	// Open the Badger database located in the ./badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions(utils.GetQuicsDirPath() + "./badger")
	db, err := badger.Open(opts)
	Badgerdb = db
	if err != nil {
		log.Fatal(err)
	}

}

func CloseDB() {
	if err := Badgerdb.Close(); err != nil {
		log.Fatal(err)
	}
}
