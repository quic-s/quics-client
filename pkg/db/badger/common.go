package badger

import (
	"fmt"
	"log"
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

func GetSyncMetadata(path string) types.SyncMetadata {
	bsyncMetadata, err := View(path)
	if err != nil {
		log.Println(err)
	}
	syncMetadata := types.SyncMetadata{}
	syncMetadata.Decode(bsyncMetadata)
	return syncMetadata

}

// Check if SyncMetadata.Conflict is nil
func IsConflictExisted(path string) bool {
	syncMetadata := GetSyncMetadata(path)
	conflictField := reflect.ValueOf(syncMetadata).FieldByName("Conflict")
	if conflictField.IsZero() {
		return false
	}
	return true
}

func IsSyncMetadataExisted(path string) bool {
	syncMetadata := GetSyncMetadata(path)
	if syncMetadata == (types.SyncMetadata{}) {
		return false
	}
	return true
}

/*---------- CONFLICT FILE LIST ------------*/

func AddConflictAndConflictFileList(path string, conflictMetadata types.ConflictMetadata) error {
	syncMetadata := GetSyncMetadata(path)
	syncMetadata.Conflict = conflictMetadata
	err := Update(path, syncMetadata.Encode())
	if err != nil {
		return err
	}

	conflictFileList := GetConflictFileList()
	newConflictFileList := types.ConflictFileList{}
	newConflictFileList = append(conflictFileList, syncMetadata)
	err = Update("ConflictFileList", newConflictFileList.Encode())
	if err != nil {
		return err
	}
	return nil

}

func RemoveConflictAndFromConflictFileList(path string) error {
	syncMetadata := GetSyncMetadata(path)
	syncMetadata.Conflict = types.ConflictMetadata{}
	err := Update(path, syncMetadata.Encode())
	if err != nil {
		return err
	}

	conflictFileList := GetConflictFileList()
	newConflictFileList := types.ConflictFileList{}
	for i, conflictFile := range conflictFileList {
		if filepath.Join(conflictFile.BeforePath, conflictFile.AfterPath) == path {
			newConflictFileList = append(conflictFileList[:i], conflictFileList[i+1:]...)
			err := Update("ConflictFileList", newConflictFileList.Encode())
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("ConflictFileList does not have %s", path)
}

func GetConflictFileList() types.ConflictFileList {
	bConflictFileList, err := View("ConflictFileList")
	if err != nil {
		log.Println(err)
	}
	conflictFileList := types.ConflictFileList{}
	conflictFileList.Decode(bConflictFileList)
	return conflictFileList
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
	rootbase := strings.Split(afterPath, "/")[0]
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

func AddRootDir(path string) {

	//If Same Absolute Path is already exist, return
	//If Same Nickname is already take , return
	rootDirList := GetRootDirList()
	for _, rootDir := range rootDirList {
		if rootDir.Path == path || rootDir.NickName == filepath.Base(path) {
			return
		}
	}
	nickname := filepath.Base(path)

	BeforePath, AfterPath := filepath.Split(path)
	rootDir := types.RootDir{
		Path:         path,
		BeforePath:   BeforePath,
		AfterPath:    "/" + AfterPath,
		NickName:     nickname,
		IsRegistered: false,
	}

	newRootDirList := types.RootDirList{}
	newRootDirList = append(rootDirList, rootDir)
	Update("RootDirList", newRootDirList.Encode())
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
