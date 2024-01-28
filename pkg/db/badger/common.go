package badger

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dgraph-io/badger/v3"

	"github.com/quic-s/quics-client/pkg/types"
)

/*--------------USER INFO ----------------*/

func (db *Badger) GetUUID() string {
	bUUID, err := db.view("UUID")
	if err != nil {
		log.Println(err)
	}
	return string(bUUID)

}

/*--------------SYNC METADATA ----------------*/

// TODO NIL 반환 , 즉 아래 두개 함수를 합쳐야함
func (db *Badger) GetSyncMetadata(path string) types.SyncMetadata {
	bsyncMetadata, err := db.view(path)
	if err != nil {
		return types.SyncMetadata{}
	}
	syncMetadata := types.SyncMetadata{}
	syncMetadata.Decode(bsyncMetadata)
	return syncMetadata

}

func (db *Badger) IsSyncMetadataExisted(path string) bool {
	syncMetadata := db.GetSyncMetadata(path)
	if reflect.ValueOf(syncMetadata).IsZero() {
		return false
	}
	return true
}

// Get All SyncMetadata in certain rootpath
// e.g. GetAllSyncMetadataInRoot("/home/username/rootdir")
func (db *Badger) GetAllSyncMetadataInRoot(rootpath string) ([]*types.SyncMetadata, error) {
	syncMetadataList := []*types.SyncMetadata{}

	// get all file path in rootpath

	prefix := rootpath
	err := db.BadgerDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			syncMetaItem := &types.SyncMetadata{}
			syncMetaItem.Decode(val)

			syncMetadataList = append(syncMetadataList, syncMetaItem)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return syncMetadataList, nil
}

// Get All SyncMetadata in all rootpath
// e.g. GetAllSyncMetadataAmongRoot()
func (db *Badger) GetAllSyncMetadataAmongRoot() ([]*types.SyncMetadata, error) {
	syncMetadataList := []*types.SyncMetadata{}
	rootDirList := db.GetRootDirList()
	for _, rootDir := range rootDirList {
		syncMetadata, err := db.GetAllSyncMetadataInRoot(rootDir.Path)
		if err != nil {
			return nil, err
		}
		syncMetadataList = append(syncMetadataList, syncMetadata...)
	}
	return syncMetadataList, nil
}

/*------------ROOT DIRECTORY---------------- */

// key : value == string : []RootDir
func (db *Badger) GetRootDirList() []types.RootDir {
	bRootDirList, err := db.view("RootDirList")
	if err != nil {
		return []types.RootDir{}
	}

	rootDirList := types.RootDirList{}
	rootDirList.Decode(bRootDirList)

	return rootDirList

}

func (db *Badger) GetRootDir(path string) types.RootDir {
	rootDirList := db.GetRootDirList()
	for _, rootDir := range rootDirList {
		if rootDir.Path == path {
			return rootDir
		}
	}
	return types.RootDir{}
}

func (db *Badger) GetBeforePathWithAfterPath(afterPath string) string {
	rootDirList := db.GetRootDirList()
	rootbase := strings.Split(afterPath, "/")[1]
	for _, rootDir := range rootDirList {
		if rootDir.AfterPath == "/"+rootbase {
			return rootDir.BeforePath
		}
	}
	return ""
}

func (db *Badger) SplitBeforeAfterRoot(path string) (string, string) {
	rootDirList := db.GetRootDirList()
	for _, rootDir := range rootDirList {
		if strings.HasPrefix(path, rootDir.BeforePath) {
			return rootDir.BeforePath, strings.TrimPrefix(path, rootDir.BeforePath)
		}
	}
	return "", ""
}

func (db *Badger) AddRootDir(path string) error {

	//If Same Absolute Path is already exist, return
	//If Same Nickname is already taken , return
	rootDirList := db.GetRootDirList()
	for _, rootDir := range rootDirList {
		if rootDir.Path == path || rootDir.NickName == filepath.Base(path) {
			return fmt.Errorf("this Directory is already exist as Root")
		}
	}
	nickname := filepath.Base(path)

	BeforePath, AfterPath := filepath.Split(path)
	rootDir := types.RootDir{
		Path:         path,
		BeforePath:   BeforePath[:len(BeforePath)-1],
		AfterPath:    "/" + AfterPath,
		NickName:     nickname,
		IsRegistered: false,
	}

	newRootDirList := types.RootDirList{}
	newRootDirList = append(rootDirList, rootDir)
	err := db.update("RootDirList", newRootDirList.Encode())
	if err != nil {
		return err
	}
	return nil
}

func (db *Badger) UpdateRootdirToRegistered(path string) error {
	rootDirList := db.GetRootDirList()
	// Update IsRegistered
	newRootDirList := types.RootDirList{}
	for _, rootDir := range rootDirList {
		if rootDir.Path == path {

			registeredRootdir := types.RootDir{
				NickName:     rootDir.NickName,
				Path:         rootDir.Path,
				BeforePath:   rootDir.BeforePath,
				AfterPath:    rootDir.AfterPath,
				IsRegistered: true,
			}

			newRootDirList = append(newRootDirList, registeredRootdir)
			continue
		}
		newRootDirList = append(newRootDirList, rootDir)
	}
	err := db.update("RootDirList", newRootDirList.Encode())
	if err != nil {
		return err
	}
	return nil
}

func (db *Badger) DeleteRootDir(path string) {
	rootDirList := db.GetRootDirList()
	newRootDirList := types.RootDirList{}
	for i, rootDir := range rootDirList {
		if rootDir.Path == path {
			newRootDirList = append(rootDirList[:i], rootDirList[i+1:]...)
			break
		}
	}

	db.update("RootDirList", newRootDirList.Encode())
}
