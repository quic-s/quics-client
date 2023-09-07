package connection

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/quic-s/quics-client/pkg/badger"
	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics-protocol/pkg/utils/fileinfo"
)

var (
	QPClient *qp.QP         // quics-protocol object 여기에는 연결을 위한 메서드들이 정의되어 있음
	Conn     *qp.Connection // connection object 여기에는 연결 그 자체가 정의되어 있음
	Watcher  *fsnotify.Watcher
)

func InitWatcher() {
	// Create a new watcher.
	err := error(nil)
	Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {

	err := error(nil)
	QPClient, err = qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		panic(err)
	}

	// initialize server
	err = QPClient.RecvMessageHandleFunc("test", func(conn *qp.Connection, msgType string, data []byte) {
		log.Println("quics-client : ", " test message received ", conn.Conn.RemoteAddr().String())
		log.Println("quics-client : ", msgType, string(data))
	})
	if err != nil {
		log.Panic(err)
	}

	// initialize server
	err = QPClient.RecvMessageHandleFunc(MUSTSYNC, func(conn *qp.Connection, msgType string, data []byte) {
		log.Println("quics-client : MUSTSYNC start ---------------")
		mustSync := &types.MustSync{}
		mustSync.Decode(data)

		log.Println("quics-client :  mustsync decode message : ", mustSync)
		// Get LocalAbsPath
		splitedAfterPath := strings.Split(mustSync.AfterPath, "/")
		rootDir := splitedAfterPath[1]            // -> rootDir
		localAbsPath := utils.GetRootDir(rootDir) // -> /a/b/rootDir
		log.Println("localAbsPath 1 : ", localAbsPath)
		if localAbsPath == "" {
			localAbsPaths := utils.GetRootDirs()
			for _, l := range localAbsPaths {
				_, rootDirName := filepath.Split(l)
				if rootDirName == rootDir {
					localAbsPath = l
					break
				}
			}
		}
		localAbsPath = filepath.Join(localAbsPath, mustSync.AfterPath[len(rootDir)+1:])
		log.Println("localAbsPath 2 : ", localAbsPath)
		// Get SyncMeta by LocalAbsPath
		if _, err := os.Stat(localAbsPath); os.IsNotExist(err) && mustSync.LatestHash == "" {
			log.Println("file deleted confirm")
			badger.Delete(localAbsPath)
			return
		}
		syncMetaByte, err := badger.View(localAbsPath)
		if err != nil {

			log.Panic(err)
		}

		syncMeta := &types.SyncMetadata{}
		syncMeta.Decode(syncMetaByte)
		log.Println("Before branch ------------------- ")
		log.Println("Timestamp (sync) : ", syncMeta.LastUpdateTimestamp)
		log.Println("Timestamp (must) : ", mustSync.LatestSyncTimestamp)
		log.Println("Hash (sync) : ", syncMeta.LastUpdateHash)
		log.Println("Hash (must) : ", mustSync.LatestHash)

		// Check Condition for OverWrite
		if syncMeta.LastSyncTimestamp == syncMeta.LastUpdateTimestamp && syncMeta.LastSyncHash == syncMeta.LastUpdateHash && mustSync.LatestSyncTimestamp > syncMeta.LastUpdateTimestamp {
			log.Println("OverWrite -------------")
			before, after := utils.SplitBeforeAfterRoot(localAbsPath)
			// PLEASEFILE
			pleaseFile := &types.PleaseFile{
				Uuid:          viper.GetViperEnvVariables("UUID"),
				BeforePath:    before,
				AfterPath:     after,
				SyncTimestamp: mustSync.LatestSyncTimestamp,
			}
			log.Println("pleaseFile : ", pleaseFile)
			Conn.SendMessage(PLEASEFILE, pleaseFile.Encode())
			// when GIVEYOUFILE is come, badger.Update(localAbsPath, syncMeta.Encode()) is envoked
		} else {
			//check "isUpdated" when broadcast

			if syncMeta.LastUpdateHash == mustSync.LatestHash && syncMeta.LastUpdateTimestamp == mustSync.LatestSyncTimestamp {
				log.Println("Is Updated when broadcast ---------")
				updateMetadata := &types.SyncMetadata{
					Path:                localAbsPath,
					LastUpdateTimestamp: syncMeta.LastUpdateTimestamp,
					LastUpdateHash:      syncMeta.LastUpdateHash,
					LastSyncTimestamp:   syncMeta.LastUpdateTimestamp,
					LastSyncHash:        syncMeta.LastUpdateHash,
				}
				log.Println("updateMetadata : ", updateMetadata)
				badger.Update(localAbsPath, updateMetadata.Encode())

			} else {
				// else, PLEASESYNC
				// log.Println("Pleasesync ---------------")
				// before, after := utils.SplitBeforeAfterRoot(localAbsPath)
				// prevSyncMetadata := types.SyncMetadata{}
				// prevSyncMetaByte, err := badger.View(localAbsPath)
				// if err != nil {
				// 	log.Panic(err)
				// }
				// prevSyncMetadata.Decode(prevSyncMetaByte)

				// pleaseSync := types.PleaseSync{
				// 	Uuid:                viper.GetViperEnvVariables("UUID"),
				// 	Event:               CONFIRM,
				// 	BeforePath:          before,
				// 	AfterPath:           after,
				// 	LastUpdateTimestamp: prevSyncMetadata.LastUpdateTimestamp,
				// 	LastUpdateHash:      prevSyncMetadata.LastUpdateHash,
				// }
				// Conn.SendFileMessage(PLEASESYNC, pleaseSync.Encode(), localAbsPath)
			}
		}

	})
	if err != nil {
		log.Println("quics-client : MUSTSYNC failed : ", err)
		log.Panic(err)
	}

	err = QPClient.RecvMessageWithResponseHandleFunc(TWOOPTIONS, func(conn *qp.Connection, msgType string, data []byte) []byte {
		twoOptions := &types.TwoOptions{}
		twoOptions.Decode(data)
		newTimestamp := uint64(0)
		if twoOptions.ClientSideTimestamp > twoOptions.ServerSideSyncTimestamp {
			newTimestamp = twoOptions.ClientSideTimestamp
		} else {
			newTimestamp = twoOptions.ServerSideSyncTimestamp
		}
		// TODO Give Options to Users, but now just choose client's
		chooseOne := &types.ChooseOne{
			BeforePath:          twoOptions.BeforePath,
			AfterPath:           twoOptions.AfterPath,
			ChosenHash:          twoOptions.ClientSideHash,
			ChosenTimestamp:     twoOptions.ClientSideTimestamp,
			LastUpdateHash:      twoOptions.ClientSideHash,
			LastUpdateTimestamp: newTimestamp,
		}

		// Get LocalAbsPath
		localAbsPath := utils.GetRootDir(twoOptions.AfterPath[1:])
		if localAbsPath == "" {
			localAbsPaths := utils.GetRootDirs()
			for _, l := range localAbsPaths {
				_, rootDirName := filepath.Split(l)
				if rootDirName == twoOptions.AfterPath[1:] {
					localAbsPath = l
					break
				}
			}
		}
		// Get SyncMeta by LocalAbsPath
		prevSyncMetaByte, err := badger.View(localAbsPath)
		if err != nil {
			log.Panic(err)
		}
		prevSyncMeta := &types.SyncMetadata{}
		prevSyncMeta.Decode(prevSyncMetaByte)

		syncMetadata := &types.SyncMetadata{
			Path:                localAbsPath,
			LastUpdateTimestamp: chooseOne.LastUpdateTimestamp,
			LastSyncTimestamp:   chooseOne.LastUpdateTimestamp,
			LastSyncHash:        chooseOne.LastUpdateHash,
			LastUpdateHash:      chooseOne.LastUpdateHash,
		}
		badger.Update(localAbsPath, syncMetadata.Encode())

		return chooseOne.Encode()

	})
	if err != nil {
		log.Panic(err)
	}

	err = QPClient.RecvFileMessageHandleFunc(GIVEYOUFILE, func(conn *qp.Connection, fileMsgType string, msgData []byte, fileInfo *fileinfo.FileInfo, fileReader io.Reader) {
		mustSync := &types.MustSync{}
		mustSync.Decode(msgData)
		localAbsPath := utils.GetRootDir(mustSync.AfterPath[1:])
		if localAbsPath == "" {
			localAbsPaths := utils.GetRootDirs()
			for _, l := range localAbsPaths {
				_, rootDirName := filepath.Split(l)
				if rootDirName == mustSync.AfterPath[1:] {
					localAbsPath = l
					break
				}
			}
		}

		// Get SyncMeta by LocalAbsPath
		syncMetaByte, err := badger.View(localAbsPath)
		if err != nil {
			log.Panic(err)
		}
		syncMeta := &types.SyncMetadata{}
		if len(syncMetaByte) != 0 {
			syncMeta.Decode(syncMetaByte)
		}
		// Check Condition for OverWrite
		if len(syncMetaByte) == 0 || syncMeta.LastSyncTimestamp == syncMeta.LastUpdateTimestamp && syncMeta.LastSyncHash == syncMeta.LastUpdateHash && mustSync.LatestSyncTimestamp > syncMeta.LastUpdateTimestamp {
			file := &os.File{}
			if _, err := os.Stat(localAbsPath); os.IsNotExist(err) {
				file, err = os.Create(localAbsPath)
			} else {
				file, err = os.Open(localAbsPath)
			}

			io.Copy(file, fileReader)
			file.Chmod(fileInfo.Mode)
			os.Chtimes(localAbsPath, time.Now(), fileInfo.ModTime)
			syncMeta = &types.SyncMetadata{
				Path:                localAbsPath,
				LastUpdateTimestamp: mustSync.LatestSyncTimestamp,
				LastSyncTimestamp:   mustSync.LatestSyncTimestamp,
				LastSyncHash:        mustSync.LatestHash,
				LastUpdateHash:      mustSync.LatestHash,
			}
			badger.Update(localAbsPath, syncMeta.Encode())
		} else {
			before, after := utils.SplitBeforeAfterRoot(localAbsPath)
			syncMetaByte, err := badger.View(localAbsPath)
			if err != nil {
				log.Panic(err)
			}
			syncMeta := &types.SyncMetadata{}
			syncMeta.Decode(syncMetaByte)
			pleaseSync := types.PleaseSync{
				Uuid:                viper.GetViperEnvVariables("UUID"),
				BeforePath:          before,
				AfterPath:           after,
				LastUpdateTimestamp: syncMeta.LastUpdateTimestamp,
				LastUpdateHash:      syncMeta.LastUpdateHash,
				Event:               WRITE,
			}
			if _, err := os.Stat(localAbsPath); os.IsNotExist(err) {
				pleaseSync.Event = REMOVE
				Conn.SendMessage(PLEASESYNC, pleaseSync.Encode())
				return
			}
			Conn.SendFileMessage(PLEASESYNC, pleaseSync.Encode(), localAbsPath)
		}

	})
	if err != nil {
		log.Panic(err)
	}

	err = QPClient.RecvFileHandleFunc("test", func(conn *qp.Connection, fileType string, fileInfo *qp.FileInfo, fileReader io.Reader) {
		log.Println("quics-client :quics-protocol: ", "file received ", fileInfo.Name)
		file, err := os.Create("received.txt")
		if err != nil {
			log.Fatal(err)
		}
		n, err := io.Copy(file, fileReader)
		if err != nil {
			log.Fatal(err)
		}
		if n != fileInfo.Size {
			log.Fatalf("quics-protocol: read only %dbytes", n)
		}
		log.Println("quics-client :quics-protocol: ", "file saved with ", n, "bytes")
	})
	if err != nil {
		log.Panic(err)
	}

	err = QPClient.RecvFileMessageHandleFunc("test", func(conn *qp.Connection, fileMsgType string, data []byte, fileInfo *qp.FileInfo, fileReader io.Reader) {
		log.Println("quics-client :quics-protocol: ", "message received ", conn.Conn.RemoteAddr().String())
		log.Println("quics-client :quics-protocol: ", fileMsgType, string(data))

		log.Println("quics-client :quics-protocol: ", "file received ", fileInfo.Name)
		file, err := os.Create("received2.txt")
		if err != nil {
			log.Fatal(err)
		}
		n, err := io.Copy(file, fileReader)
		if err != nil {
			log.Fatal(err)
		}
		if n != fileInfo.Size {
			log.Println("quics-client :quics-protocol: ", "read only ", n, "bytes")
			log.Fatal(err)
		}
		log.Println("quics-client :quics-protocol: ", "file saved")
	})
	if err != nil {
		log.Panic(err)
	}
}

func CloseConnect() {
	Conn.Close()
	Conn = nil
}

func ReConnect() {
	p, err := strconv.Atoi(viper.GetViperEnvVariables("QUICS_SERVER_PORT"))
	if err != nil {
		panic(err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-s"},
	}
	Conn, err = QPClient.Dial(&net.UDPAddr{IP: net.ParseIP(viper.GetViperEnvVariables("QUICS_SERVER_IP")), Port: p}, tlsConf)
	if err != nil {
		panic(err)
	}
}
