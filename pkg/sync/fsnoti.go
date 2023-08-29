package sync

// func DirWatchStart() {
// 	// Create a new watcher.
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer watcher.Close()

// 	// Define a channel to receive events.
// 	done := make(chan bool)

// 	go func() {
// 		for {
// 			select {
// 			case event, ok := <-watcher.Events:
// 				if !ok {
// 					return
// 				}
// 				log.Println("event:", event)

// 				//이벤트가 일어난 파일만 전송
// 				filepath := event.Name

// 				if info, err := os.Stat(filepath); event.Op&fsnotify.Create == fsnotify.Create && err == nil && info.IsDir() { // 디렉토리 인것
// 					watcher.Add(filepath)
// 				}
// 				if info, err := os.Stat(filepath); event.Op&fsnotify.Create == fsnotify.Create && err == nil && !info.IsDir() { //파일 인 것
// 					//RestClientFile(filepath)
// 				}

// 				if event.Op&fsnotify.Remove == fsnotify.Remove {
// 					filepath = utils.LocalAbsToRoot(filepath, getAccelerDir())
// 					qc.ClientMessage("removed", []byte(filepath))
// 				}

// 				if event.Op&fsnotify.Write == fsnotify.Write {

// 					RestClientFile(filepath)
// 				}
// 			case err, ok := <-watcher.Errors:
// 				if !ok {
// 					return
// 				}
// 				log.Println("error:", err)
// 			}
// 		}
// 	}()

// 	filepath.Walk(getAccelerDir(), func(path string, info os.FileInfo, err error) error {
// 		if info.IsDir() {
// 			watcher.Add(path)
// 		}
// 		return nil
// 	})

// 	// Wait until the channel is closed.
// 	<-done

// }
