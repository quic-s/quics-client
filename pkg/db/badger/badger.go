package badger

import (
	"github.com/dgraph-io/badger/v3"
)

func (bg *Badger) Update(key string, value []byte) error {
	mutex := &bg.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	return bg.BadgerDB.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(key), value); err != nil {
			return err
		}
		return nil
	})
}

func (bg *Badger) View(key string) ([]byte, error) {
	mutex := &bg.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	var value []byte

	err := bg.BadgerDB.View(func(txn *badger.Txn) error {

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

func (bg *Badger) Delete(key string) error {
	mutex := &bg.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	err := bg.BadgerDB.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
	return err
}
