package connection

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/quic-s/quics-client/pkg/badger"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
)

func DirWatchStart() {
	go func() {
		for {
			select {
			case event, ok := <-Watcher.Events:
				if !ok {
					return
				}
				//event file
				filepath := event.Name

				if info, err := os.Stat(filepath); event.Op&fsnotify.Create == fsnotify.Create && err == nil && info.IsDir() { // IsDirectory
					Watcher.Add(filepath)
				}
				if info, err := os.Stat(filepath); event.Op&fsnotify.Create == fsnotify.Create && err == nil && !info.IsDir() { // IsFile
					log.Println("quics-client : CREATE event ")
					before, after := utils.SplitBeforeAfterRoot(filepath)
					//get hash
					h := sha512.New()
					h.Write([]byte(after)) // /root/*
					h.Write([]byte(info.ModTime().String()))
					h.Write([]byte(info.Mode().String()))
					h.Write([]byte(fmt.Sprint(info.Size())))

					syncMetadata := types.SyncMetadata{
						Path:                filepath,
						LastUpdateTimestamp: 1,
						LastUpdateHash:      hex.EncodeToString(h.Sum(nil)),
						LastSyncTimestamp:   0,
						LastSyncHash:        "",
					}
					badger.Update(filepath, syncMetadata.Encode())
					body := types.PleaseSync{
						Uuid:                viper.GetViperEnvVariables("UUID"),
						Event:               CREATE,
						BeforePath:          before,
						AfterPath:           after,
						LastUpdateTimestamp: syncMetadata.LastUpdateTimestamp,
						LastUpdateHash:      syncMetadata.LastUpdateHash,
					}
					Conn.SendFileMessage(PLEASESYNC, body.Encode(), filepath)
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("quics-client : REMOVE event ")
					Conn.SendMessage(REMOVE, []byte(filepath))
					before, after := utils.SplitBeforeAfterRoot(filepath)

					prevSyncMetaByte, err := badger.View(filepath)
					if err != nil {
						log.Println(err)
						continue
					}
					PrevSyncMetadata := types.SyncMetadata{}
					PrevSyncMetadata.Decode(prevSyncMetaByte)

					pleaseSync := types.PleaseSync{
						Uuid:                viper.GetViperEnvVariables("UUID"),
						Event:               REMOVE,
						BeforePath:          before,
						AfterPath:           after,
						LastUpdateTimestamp: PrevSyncMetadata.LastUpdateTimestamp + 1,
						LastUpdateHash:      "",
					}
					Conn.SendFileMessage(PLEASESYNC, pleaseSync.Encode(), filepath)

				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("quics-client : WRITE event ")
					info, err := os.Stat(filepath)
					if err != nil {
						log.Panicln(err)
					}
					before, after := utils.SplitBeforeAfterRoot(filepath)
					//get hash
					h := sha512.New()
					h.Write([]byte(after)) // /root/*
					h.Write([]byte(info.ModTime().String()))
					h.Write([]byte(info.Mode().String()))
					h.Write([]byte(fmt.Sprint(info.Size())))

					prevSyncMetaByte, err := badger.View(filepath)
					if err != nil {
						log.Println(err)
					}
					prevSyncMetadata := types.SyncMetadata{}
					prevSyncMetadata.Decode(prevSyncMetaByte)

					SyncMetadata := types.SyncMetadata{
						Path:                prevSyncMetadata.Path,
						LastUpdateTimestamp: prevSyncMetadata.LastUpdateTimestamp + 1,
						LastUpdateHash:      hex.EncodeToString(h.Sum(nil)), // make new hash
						LastSyncTimestamp:   prevSyncMetadata.LastSyncTimestamp,
						LastSyncHash:        prevSyncMetadata.LastSyncHash,
					}
					badger.Update(filepath, SyncMetadata.Encode())

					body := types.PleaseSync{
						Uuid:                viper.GetViperEnvVariables("UUID"),
						Event:               WRITE,
						BeforePath:          before,
						AfterPath:           after,
						LastUpdateTimestamp: SyncMetadata.LastUpdateTimestamp,
						LastUpdateHash:      SyncMetadata.LastUpdateHash,
					}
					Conn.SendFileMessage(PLEASESYNC, body.Encode(), filepath)
				}
			case err, ok := <-Watcher.Errors:
				if !ok {
					return
				}
				log.Println("quics-cleint : ", err)
			}
		}
	}()
}

func DirWatchAdd(rootpath string) {
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err := Watcher.Add(path)
			if err != nil {
				log.Println("quics-client : Watcher Adding unsuccessful ", err)
			}

		}

		return nil
	})
	if err != nil {
		log.Println("quics-client : ", err)
	}

}

func DirWatchStop(rootpath string) {
	filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			Watcher.Remove(path)
		}
		return nil
	})

}

func WatchStop() {
	Watcher.Close()
}
