package sync

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-client/pkg/viper"
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics-protocol/pkg/stream"
	qstypes "github.com/quic-s/quics/pkg/types"
)

var ClientPW string

//var PingCnt = 0

// @URL /api/v1/connect/server
// ex. RegisterClientRequest("password", "S_IP", "S_PORT")
func ClientRegistration(ClientPassword string, SIP string, SPort string) error {
	if SIP != "" {
		viper.WriteViperEnvVariables("QUICS_SERVER_IP", SIP)
	}
	if SPort != "" {
		viper.WriteViperEnvVariables("QUICS_SERVER_PORT", SPort)
	}
	ClientPW = ClientPassword
	// Check UUID is existed
	UUID := ""
	if badger.GetUUID() != "" {
		UUID = badger.GetUUID()
	} else {
		UUID = uuid.New().String()
		badger.Update("UUID", []byte(UUID))
	}

	InitQPClient()

	var err error
	Conn, err = ConnectServer(UUID)
	if err != nil {
		return err
	}

	viper.WriteViperEnvVariables("QUICS_SERVER_PASSWORD", ClientPassword)

	go Rescan()
	go Reconnect()

	return nil
}

func CheckInternetConnection() bool {
	err := Conn.OpenTransaction("PING", func(stream *stream.Stream, transactionName string, transactionID []byte) error {
		res, err := qclient.SendPing(stream, badger.GetUUID())
		if err != nil {
			log.Println("error in transaction: ", err)
			return err
		}
		if reflect.ValueOf(res).IsZero() {
			return fmt.Errorf("ping do not success")
		}
		return nil
	})
	if err != nil {
		log.Println("quics-client : [PING] ERROR ", err)
		return false
	}
	//log.Println("PingCnt >> ", PingCnt)
	//PingCnt++
	return true
}

func Reconnect() {

	for {
		isOnline := CheckInternetConnection()
		//if !isOnline || PingCnt == 80 {
		if !isOnline {
			//PingCnt = 0
			ip := viper.GetViperEnvVariables("QUICS_SERVER_IP")
			if ip == "" {
				log.Println("quics-client : [RECONNECT] ERROR ", "server IP is not set")
				continue
			}

			err := Conn.Close()
			if err != nil {
				log.Println("quics-client : [RECONNECT] ERROR ", err)
			}
			Conn = nil
			err = QPClient.Close()
			if err != nil {
				log.Println("quics-client : [RECONNECT] ERROR ", err)
			}
			QPClient = nil

			port := viper.GetViperEnvVariables("QUICS_SERVER_PORT")
			err = ClientRegistration(ClientPW, ip, port)
			if err != nil {
				log.Println(err)
				continue
			}
			return
		}

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
		viper.GetViperEnvVariables("QUICS_SERVER_IP"),
		parsedPort,
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

// @URL /api/v1/disconnect/server
// ex. DisconnectClient()
func DisconnectClient() error {

	err := Conn.OpenTransaction("DISCONNECTCLIENT", func(stream *stream.Stream, transactionName string, transactionID []byte) error {
		res, err := qclient.SendDisconnectClient(stream, badger.GetUUID())
		if err != nil {

			return err
		}
		if reflect.ValueOf(res).IsZero() {

			return fmt.Errorf("[DISCONNECTCLIENT] still registerd")
		}

		return nil
	})
	if err != nil {
		return err
	}
	viper.DeleteViperVariablesByKey("QUICS_SERVER_IP")
	viper.DeleteViperVariablesByKey("QUICS_SERVER_PORT")
	viper.DeleteViperVariablesByKey("QUICS_SERVER_PASSWORD")

	Conn.Close()
	Conn = nil
	QPClient.Close()
	QPClient = nil

	return nil
}

// @URL /api/v1/disconnect/root
// ex. DisConnectRootDir( "/home/ubuntu/rootDir")
func DisconnectRootDir(path string) (qstypes.DisconnectRootDirRes, error) {
	_, after := badger.SplitBeforeAfterRoot(path)

	result := qstypes.DisconnectRootDirRes{}
	err := Conn.OpenTransaction(qstypes.DISCONNECTROOTDIR, func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		res, err := qclient.SendDisconnectRootDir(stream, badger.GetUUID(), after)
		if err != nil {

			return err
		}
		if reflect.ValueOf(res).IsZero() {
			return fmt.Errorf("[DISCONNECTROOTDIR] still registerd")
		}

		result = res

		DirWatchStop(path)
		syncList, err := badger.GetAllSyncMetadataInRoot(path)
		if err != nil {
			log.Println(err)
		}
		for _, item := range syncList {
			badger.Delete(filepath.Join(item.BeforePath, item.AfterPath))
		}
		badger.DeleteRootDir(path)
		return nil

	})
	if err != nil {
		log.Println("log")
		return result, err
	}
	return result, nil

}
