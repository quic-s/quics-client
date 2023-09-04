package badger

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

func Update(key string, value []byte) error {
	db := Badgerdb

	return db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(key), value); err != nil {
			return err
		}
		return nil
	})
}

func View(key string) ([]byte, error) {
	//db := InitDB()
	var value []byte
	db := Badgerdb

	err := db.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(key))
		if err != nil {

			return err
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		log.Printf("The value of key is %s\n", string(value))
		return nil
	})
	return value, err
}

func Delete(key string) error {
	db := Badgerdb
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
	return err
}
