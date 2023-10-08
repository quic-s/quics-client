package sync

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/google/uuid"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-client/pkg/viper"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

// @URL /api/v1/connect/server
// ex. RegisterClientRequest("password", "S_IP", "S_PORT")
func ClientRegistration(ClientPassword string, SIP string, SPort string) error {
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

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-s"},
	}

	parsedPort, err := strconv.Atoi(viper.GetViperEnvVariables("QUICS_SERVER_PORT"))
	if err != nil {
		return fmt.Errorf("[ClientRegistration] %s", err)
	}

	NewConn, err := QPClient.DialWithTransaction(
		&net.UDPAddr{
			IP:   net.ParseIP(viper.GetViperEnvVariables("QUICS_SERVER_IP")),
			Port: parsedPort},
		tlsConf,
		qstypes.REGISTERCLIENT,
		func(stream *qp.Stream, transactionName string, transactionID []byte) error {
			clientRegisterRes, err := qclient.SendClientRegister(stream, UUID, ClientPassword)
			if err != nil {
				return err
			}
			if clientRegisterRes.UUID == "" {
				return fmt.Errorf("Server cannot register client")
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("[ClientRegistration] %s", err)
	}
	Conn = NewConn
	return nil
}

// @URL /api/v1/connect/root/local
// ex. RegistRootDir("/home/ubuntu/rootDir", "password")
func RegistRootDir(LocalRootDir string, RootDirPW string, Side string) error {
	dir, root := filepath.Split(LocalRootDir)
	log.Println("[RegisterRootDir]")

	// Add To RootDirList
	badger.AddRootDir(LocalRootDir)
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

	err := Conn.OpenTransaction(transcationName, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
		//TODO When get root folder from root then What about badger
		registerRes, err := qclient.SendRootDirRegister(stream, badger.GetUUID(), RootDirPW, dir, "/"+root)
		if err != nil {

			return err
		}
		if reflect.ValueOf(registerRes).IsZero() {
			return fmt.Errorf("server cannot register root directory")
		}

		// Update IsRegistered
		rootdir.IsRegistered = true
		badger.Update(LocalRootDir, rootdir.Encode())

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