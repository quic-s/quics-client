package sync

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/quic-s/quics-client/pkg/db/badger"
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
				path := event.Name

				if badger.IsConflictExisted(path) {
					continue
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("quics-client : REMOVE event ")
					go PSwhenRemove(path)
					continue

				}
				info, err := os.Stat(path)
				if err != nil {
					log.Println("quics-client : ", err)
					continue
				}

				if event.Op&fsnotify.Create == fsnotify.Create && info.IsDir() { // IsDirectory
					Watcher.Add(path)
					continue
				}
				if event.Op&fsnotify.Create == fsnotify.Create && !info.IsDir() { // IsFile
					log.Println("quics-client : CREATE event ")
					go PSwhenCreate(path, info)
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("quics-client : WRITE event ")
					go PSwhenWrite(path, info)
					continue

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
