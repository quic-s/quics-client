package sync

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func DirWatchStart() {

	go func() {
		for {
			select {
			case event, ok := <-Watcher.Events:
				if !ok {
					continue
				}

				//event file
				path := event.Name

				// lock mutex by hash value of file path
				// using hash value is to reduce the number of mutex

				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("quics-client : REMOVE event ")
					go PleaseSync(path)
					continue
				}
				info, err := os.Stat(path)
				if err != nil {
					log.Println("quics-client : ", err)
					continue
				}
				// continue this case, when PS happened by MS

				if event.Op&fsnotify.Create == fsnotify.Create && info.IsDir() { // IsDirectory
					DirWatchAdd(path)
					continue
				}

				if event.Op&fsnotify.Create == fsnotify.Create && !info.IsDir() { // IsFile
					log.Println("quics-client : CREATE event ")
					go PleaseSync(path)
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("quics-client : WRITE event ")
					go PleaseSync(path)
					continue
				}

			case err, ok := <-Watcher.Errors:
				if !ok {
					continue
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
