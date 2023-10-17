package sync

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

var ClientPW string

// @URL /api/v1/connect/server
// ex. RegisterClientRequest("password", "S_IP", "S_PORT")
func ClientRegistration(ClientPassword string, SIP string, SPort string) error {
	ClientPW = ClientPassword
	if SIP != "" {
		viper.WriteViperEnvVariables("QUICS_SERVER_IP", SIP)
	}
	if SPort != "" {
		viper.WriteViperEnvVariables("QUICS_SERVER_PORT", SPort)
	}

	// Check UUID is existed
	UUID := ""
	if badger.GetUUID() != "" {
		UUID = badger.GetUUID()
	} else {
		UUID = uuid.New().String()
		badger.Update("UUID", []byte(UUID))
	}

	var err error
	Conn, err = ConnectServer(UUID)
	if err != nil {
		return err
	}

	go Rescan()
	go Reconnect()

	return nil
}

func Reconnect() {
	prevStatus := true
	for {

		isOnline := utils.CheckInternetConnection()
		if !prevStatus && isOnline {
			ip := viper.GetViperEnvVariables("QUICS_SERVER_IP")
			if ip == "" {
				continue
			}

			Conn.Close()
			Conn = nil
			QPClient.Close()
			QPClient = nil
			InitQPClient()

			port := viper.GetViperEnvVariables("QUICS_SERVER_PORT")
			err := ClientRegistration(ClientPW, ip, port)
			if err != nil {
				log.Println(err)
			}

		}
		prevStatus = isOnline
		time.Sleep(5 * time.Second)
	}
}

func ConnectServer(uuid string) (*qp.Connection, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-s"},
	}

	parsedPort, err := strconv.Atoi(viper.GetViperEnvVariables("QUICS_SERVER_PORT"))
	if err != nil {
		return nil, fmt.Errorf("[ClientRegistration] %s", err)
	}

	NewConn, err := QPClient.DialWithTransaction(
		&net.UDPAddr{
			IP:   net.ParseIP(viper.GetViperEnvVariables("QUICS_SERVER_IP")),
			Port: parsedPort},
		tlsConf,
		qstypes.REGISTERCLIENT,
		func(stream *qp.Stream, transactionName string, transactionID []byte) error {
			clientRegisterRes, err := qclient.SendClientRegister(stream, uuid, ClientPW)
			if err != nil {
				return err
			}
			if clientRegisterRes.UUID == "" {
				return fmt.Errorf("Server cannot register client")
			}
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("[ClientRegistration] %s", err)
	}
	return NewConn, nil
}

// @URL /api/v1/connect/root/local
// ex. RegistRootDir("/home/ubuntu/rootDir", "password")
func RegistRootDir(LocalRootDir string, RootDirPW string, Side string) error {
	dir, root := filepath.Split(LocalRootDir)
	log.Println("[RegisterRootDir]")

	//In Remote Register case
	_, err := os.Stat(LocalRootDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(LocalRootDir, 0755)
		if err != nil {
			return err
		}

	}

	// Add To RootDirList
	err = badger.AddRootDir(LocalRootDir)
	if err != nil {
		return fmt.Errorf("[RegisterRootDir] %s", err)
	}

	rootdir := badger.GetRootDir(LocalRootDir)

	if reflect.ValueOf(rootdir).FieldByName("IsRegistered").Bool() == true {
		return fmt.Errorf("[RegisterRootDir] Root Directory named ", rootdir.NickName, " is already registered")
	}

	// Named Transaction
	transcationName := ""
	if Side == "LOCAL" {
		transcationName = qstypes.REGISTERROOTDIR
	} else if Side == "REMOTE" {
		transcationName = qstypes.SYNCROOTDIR
	}

	err = Conn.OpenTransaction(transcationName, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
		//TODO When get root folder from root then What about badger
		registerRes, err := qclient.SendRootDirRegister(stream, badger.GetUUID(), RootDirPW, dir, "/"+root)
		if err != nil {

			return err
		}
		if reflect.ValueOf(registerRes).IsZero() {
			return fmt.Errorf("server cannot register root directory")
		}

		err = badger.UpdateRootdirToRegistered(LocalRootDir)
		if err != nil {
			return err
		}

		DirWatchAdd(LocalRootDir)
		return nil

	})
	if err != nil {
		return fmt.Errorf("[RegisterRootDir] %s", err)
	}
	return nil
}

// @URL /api/v1/connect/list/remote
// ex. GetRemoteRootList()
func GetRemoteRootList() (qstypes.AskRootDirRes, error) {
	rootList := &qstypes.AskRootDirRes{}
	err := Conn.OpenTransaction(qstypes.GETROOTDIRS, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
		askRootDirRes, err := qclient.SendAskRootList(stream, badger.GetUUID())
		if err != nil {
			return err
		}
		rootList = askRootDirRes
		return nil
	})
	if err != nil {
		return *rootList, fmt.Errorf("[GetRemoteRootList] %s", err)
	}
	log.Println("rootList : ", rootList)
	return *rootList, nil
}

// // @URL /api/v1/disconnect/root
// // ex. UnRegisterRootDirRequest("/home/ubuntu/rootDir", "password")
// func UnRegisterRootDirRequest(DisconnectRootDir string, RootDirPW string) error {

// 	_, file := filepath.Split(DisconnectRootDir)
// 	before, after := utils.SplitBeforeAfterRoot(DisconnectRootDir)

// 	body := types.RegisterRootDirRequest{
// 		UUID:            badger.GetUUID(),
// 		RootDirPassword: RootDirPW,
// 		BeforePath:      before,
// 		AfterPath:       after,
// 	}
// 	response, err := Conn.SendMessageWithResponse("NOTROOTDIRANYMORE", body.Encode())
// 	if err != nil {
// 		return err
// 	}
// 	if string(response) == "OK" {
// 		badger.DeleteRootDir(DisconnectRootDir)
// 		return nil
// 		//TODO make clear from fsnotify

// 	}
// 	return fmt.Errorf("RegisterRootDir: %s", string(response))
// }

// // @URL /api/v1/disconnect/server
// // ex. DisconnectClientRequest("password")
// func DisconnectClientRequest(password string) {
// 	NotClientAnymoreRequest := qstypes.NotClientAnymoreRequest{
// 		UUID:           badger.GetUUID(),
// 		ClientPassword: password,
// 	}
// 	resp, err := Conn.SendMessageWithResponse(NOTCLIENTANYMORE, NotClientAnymoreRequest.Encode())
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println("quics-client :response : " + string(resp))
// 	CloseConnect()
// }

// // @URL /api/v1/disconnect/root
// // ex. DisConnectClientRequest( "/home/ubuntu/rootDir", "password")
// func DisconnectRootDirRequest(rootpath string, password string) {
// 	_, after := utils.SplitBeforeAfterRoot(rootpath)
// 	notRootDirAnymorRequest := types.NotRootDirAnymorRequest{
// 		UUID:            badger.GetUUID(),
// 		RootPath:        after, // /rootDir
// 		RootDirPassword: password,
// 	}
// 	resp, err := Conn.SendMessageWithResponse(NOTROOTDIRANYMORE, notRootDirAnymorRequest.Encode())
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println("quics-client :response : " + string(resp))

// }
