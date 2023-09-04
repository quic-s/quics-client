package sync

import (
	"github.com/quic-s/quics-client/pkg/badger"
	"github.com/quic-s/quics-client/pkg/viper"
)

//TODO
// 1. 폴더를 connection을 맺는다, badger에 메타데이터를 저장시작한다
// 2. 폴더의 변경사항을 감지한다
// 3. 변경사항을 서버에 전송한다
// 4. 몇분에 한번씩 서버에서 변경사항을 받는다
// 5. 변경사항을 받아서 파일을 수정/생성/삭제 한다
// 6. disable할 경우 파일의 연결을 끊고 badger에 저장된 메타데이터를 삭제한다

func MakeLocalSync(localrootdir string) {

	if certaindir := viper.GetViperEnvVariables(localrootdir); certaindir != "" {
		//make sync with centain directory

	} else {
		//make sync with LocalRootDir
	}
}

func MakeRemoteSync(remoterootdir string) {
	if certaindir := viper.GetViperEnvVariables(remoterootdir); certaindir != "" {
		//make sync with centain directory

	} else {
		//make sync with RemoteRootDir
	}
}

func MakeDisableSync(disablerootdir string) {
	if certaindir := viper.GetViperEnvVariables(disablerootdir); certaindir != "" {
		//make sync with centain directory

	} else {
		//make sync with DisableRootDir
	}
}

func CanOverWrite(lastUpdate uint64, lastSync uint64, lastestSync uint64) bool {
	if lastUpdate != lastSync {
		return false
	}
	if lastestSync <= lastSync {
		return false
	}
	return true
}

func UpdateSyncMetadata(syncMetadata SyncMetadata) error {
	error := badger.Update(syncMetadata.Path, syncMetadata.Encode())
	if error != nil {
		return error
	}
	return nil
}
