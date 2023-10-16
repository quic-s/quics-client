package badger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/quic-s/quics-client/pkg/types"
)

func GetUUID() string {
	bUUID, err := View("UUID")
	if err != nil {
		log.Println(err)
	}
	return string(bUUID)

}

/*--------------SYNC METADATA ----------------*/

// TODO NIL 반환 , 즉 아래 두개 함수를 합쳐야함
func GetSyncMetadata(path string) types.SyncMetadata {
	bsyncMetadata, err := View(path)
	if err != nil {
		return types.SyncMetadata{}
	}
	syncMetadata := types.SyncMetadata{}
	syncMetadata.Decode(bsyncMetadata)
	return syncMetadata

}

func IsSyncMetadataExisted(path string) bool {
	syncMetadata := GetSyncMetadata(path)
	if reflect.ValueOf(syncMetadata).IsZero() {
		return false
	}
	return true
}

// Get All SyncMetadata in certain rootpath
// e.g. GetAllSyncMetadataInRoot("/home/username/rootdir")
func GetAllSyncMetadataInRoot(rootpath string) ([]*types.SyncMetadata, error) {
	syncMetadataList := []*types.SyncMetadata{}

	// get all file path in rootpath
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// if file is found then add to syncMetadataList else if dir is found then keep walk and add
		if !info.IsDir() {
			syncMetadata := GetSyncMetadata(path)
			syncMetadataList = append(syncMetadataList, &syncMetadata)
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
func GetAllSyncMetadataAmongRoot() ([]*types.SyncMetadata, error) {
	syncMetadataList := []*types.SyncMetadata{}
	rootDirList := GetRootDirList()
	for _, rootDir := range rootDirList {

		syncMetadata, err := GetAllSyncMetadataInRoot(rootDir.Path)
		if err != nil {
			return nil, err
		}
		syncMetadataList = append(syncMetadataList, syncMetadata...)
	}
	return syncMetadataList, nil
}

/*------------ROOT DIRECTORY---------------- */

// key : value == string : []RootDir
func GetRootDirList() []types.RootDir {
	bRootDirList, err := View("RootDirList")
	if err != nil {
		return []types.RootDir{}
	}
	rootDirList := types.RootDirList{}
	rootDirList.Decode(bRootDirList)

	for i, rootDir := range rootDirList {
		if !rootDir.IsRegistered {
			rootDirList = append(rootDirList[:i], rootDirList[i+1:]...)
		}
	}

	return rootDirList

}

func GetRootDir(path string) types.RootDir {
	rootDirList := GetRootDirList()
	for _, rootDir := range rootDirList {
		if rootDir.Path == path {
			return rootDir
		}
	}
	return types.RootDir{}
}

func GetBeforePathWithAfterPath(afterPath string) string {
	rootDirList := GetRootDirList()
	rootbase := strings.Split(afterPath, "/")[1]
	for _, rootDir := range rootDirList {
		if rootDir.AfterPath == "/"+rootbase {
			return rootDir.BeforePath
		}
	}
	return ""
}

func SplitBeforeAfterRoot(path string) (string, string) {
	rootDirList := GetRootDirList()
	for _, rootDir := range rootDirList {
		if strings.HasPrefix(path, rootDir.BeforePath) {
			return rootDir.BeforePath, strings.TrimPrefix(path, rootDir.BeforePath)
		}
	}
	return "", ""
}

func AddRootDir(path string) error {

	//If Same Absolute Path is already exist, return
	//If Same Nickname is already taken , return
	rootDirList := GetRootDirList()
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
	err := Update("RootDirList", newRootDirList.Encode())
	if err != nil {
		return err
	}
	return nil
}

func DeleteRootDir(path string) {
	rootDirList := GetRootDirList()
	newRootDirList := types.RootDirList{}
	for i, rootDir := range rootDirList {
		if rootDir.Path == path {
			newRootDirList = append(rootDirList[:i], rootDirList[i+1:]...)
			break
		}
	}

	Update("RootDirList", newRootDirList.Encode())
}
