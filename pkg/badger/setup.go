package badger

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics-client/pkg/utils"
)

// Declare a global variable for the DB.
var badgerdb *badger.DB

// Define a function to open the DB.
func init() {
	// Open the Badger database located in the ./badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions(utils.GetQuicsDirPath() + "./badger")
	db, err := badger.Open(opts)
	badgerdb = db
	if err != nil {
		log.Fatal(err)
	}

	// Register a function to close the DB when the program exits.
	defer closeDB()
}

// Define a function to close the DB.
func closeDB() {
	if err := badgerdb.Close(); err != nil {
		log.Fatal(err)
	}
}
