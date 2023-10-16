package badger

import (
	"github.com/dgraph-io/badger/v3"
)

type CRUD interface {
	Update(key string, value []byte) error
	View(key string) ([]byte, error)
	Delete(key string) error
}

func Update(key string, value []byte) error {
	mutex.Lock()
	defer mutex.Unlock()
	return badgerdb.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(key), value); err != nil {
			return err
		}
		return nil
	})
}

func View(key string) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()
	var value []byte

	err := badgerdb.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(key))
		if err != nil {

			return err
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		//log.Printf("The value of key is %s\n", string(value))
		return nil
	})
	return value, err
}

func Delete(key string) error {
	mutex.Lock()
	defer mutex.Unlock()
	err := badgerdb.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
	return err
}
