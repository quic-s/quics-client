package quic

//qp "github.com/quic-s/quics-protocol"

// func ClientFile(filepath string) {

// 	// initialize client
// 	quicClient, err := qp.New()
// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	// start client
// 	err = quicClient.Dial(host + ":" + port)
// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	file, err := os.ReadFile(filepath)
// 	if err != nil {
// 		log.Println(err)
// 		return

// 	}
// 	quicClient.SendFile(utils.LocalAbsToRoot(filepath, getAccelerDir()), file)

// 	// delay for waiting message sent to server
// 	time.Sleep(3 * time.Second)
// 	quicClient.Close()
// }
