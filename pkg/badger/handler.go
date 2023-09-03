package badger

import (
	"github.com/dgraph-io/badger/v3"
)

func Update(db *badger.DB, key []byte, value []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(key, value); err != nil {
			return err
		}
		return nil
	})
}

func View(db *badger.DB, key []byte) ([]byte, error) {
	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	return value, err
}

func Delete(db *badger.DB, key string) error {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
	return err
}
