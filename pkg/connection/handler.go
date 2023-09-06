package connection

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"time"

	"github.com/google/uuid"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
)

// message type
const (
	CREATE  string = "CREATE"
	DELETE  string = "DELETE"
	WRITE   string = "WRITE"
	RENAME  string = "RENAME"
	CONFIRM string = "CONFIRM"

	CLIENT            string = "CLIENT"
	NOTCLIENTANYMORE  string = "NOTCLIENTANYMORE"
	NOTROOTDIRANYMORE string = "NOTROOTDIRANYMORE"
	LOCALROOT         string = "LOCALROOT"
	REMOTEROOT        string = "REMOTEROOT"
	RESCAN            string = "RESCAN"
	SHARING           string = "SHARING"
	PLEASEFILE        string = "PLEASEFILE"
	PLEASESYNC        string = "PLEASESYNC"
	SHOWREMOTELIST    string = "SHOWREMOTELIST"
	CHOOSEONE         string = "CHOOSEONE"

	MUSTSYNC    string = "MUSTSYNC"
	TWOOPTIONS  string = "TWOOPTIONS"
	GIVEYOUFILE string = "GIVEYOUFILE"
)

func ClientFirstMessage(msgType string, message []byte) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-s"},
	}
	log.Println(viper.GetViperEnvVariables("QUICS_SERVER_IP"))
	log.Println(viper.GetViperEnvVariables("QUICS_SERVER_PORT"))

	// start client
	parsedPort, err := strconv.Atoi(viper.GetViperEnvVariables("QUICS_SERVER_PORT"))
	if err != nil {
		log.Println("quics-client :quics-client: ", err)
	}
	Conn, err = QPClient.DialWithMessage(&net.UDPAddr{IP: net.ParseIP(viper.GetViperEnvVariables("QUICS_SERVER_IP")), Port: parsedPort}, tlsConf, msgType, message)
	if err != nil {
		log.Println("quics-client :quics-protocol: ", err)
	}
	if err != nil {
		log.Println("quics-client :quics-protocol: ", err)
	}
}

// TODO password to crypto
func RegisterLocalRootDir(filepath string, password string) error {
	before, after := utils.SplitBeforeAfterRoot(filepath)
	body := types.RegisterRootDirRequest{
		Uuid:            viper.GetViperEnvVariables("UUID"),
		RootDirPassword: password,
		BeforePath:      before,
		AfterPath:       after,
	}
	response, err := Conn.SendMessageWithResponse(LOCALROOT, body.Encode())
	if err != nil {
		return err
	}
	if string(response) == "OK" {
		return nil
	}
	return fmt.Errorf("RegisterRootDir: %s", string(response))
}

func Ping() bool {
	timeout := time.Duration(1 * time.Second)
	conn, err := net.DialTimeout("ip4:icmp", viper.GetViperEnvVariables("QUICS_SERVER_IP"), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// @URL /api/v1/connect/root/local
// ex. RegisterLocalRootDirRequest("/home/ubuntu/rootDir", "password")
func RegisterLocalRootDirRequest(LocalRootDir string, RootDirPW string) error {

	_, file := filepath.Split(LocalRootDir)
	before, after := utils.SplitBeforeAfterRoot(LocalRootDir)

	body := types.RegisterRootDirRequest{
		Uuid:            viper.GetViperEnvVariables("UUID"),
		RootDirPassword: RootDirPW,
		BeforePath:      before,
		AfterPath:       after,
	}
	response, err := Conn.SendMessageWithResponse(LOCALROOT, body.Encode())
	if err != nil {
		return err
	}
	if string(response) == "OK" {
		viper.WriteViperEnvVariables("ROOT."+file, LocalRootDir)
		DirWatchAdd(LocalRootDir)
		return nil
	}
	return fmt.Errorf("RegisterRootDir: %s", string(response))
}

// @URL /api/v1/connect/root/remote
// ex. RegisterRemoteRootDirRequest("/home/ubuntu/rootDir", "password")
func RegisterRemoteRootDirRequest(LocalRootAbsPath string, RootDirPW string) error {

	//DISCUSS Remote Root Dir의 리스트를 받는다고하면,
	//remote rootdir의 전체 경로가 아닌 폴더명만 오기 때문에
	// 그 폴더명을 이용해서 저장할 LocalRootAbsPath를 입력받게 한다.
	_, rootdir := filepath.Split(LocalRootAbsPath)
	before, after := utils.SplitBeforeAfterRoot(LocalRootAbsPath)

	body := types.RegisterRootDirRequest{
		Uuid:            viper.GetViperEnvVariables("UUID"),
		RootDirPassword: RootDirPW,
		BeforePath:      before,
		AfterPath:       after, // /rootDir
	}
	response, err := Conn.SendMessageWithResponse(REMOTEROOT, body.Encode())
	if err != nil {
		return err
	}
	if string(response) == "OK" {
		viper.WriteViperEnvVariables("ROOT."+rootdir, LocalRootAbsPath)
		_, err := os.Stat(LocalRootAbsPath)
		if os.IsNotExist(err) {
			log.Println("quics-client :quics-client: create new root directory coming from remote")
			os.Create(LocalRootAbsPath)
		}

		Conn.SendMessage(RESCAN, []byte(""))
		return nil
	}
	return fmt.Errorf("Register Remote Root Dir: %s", string(response))

}

// @URL /api/v1/disconnect/root
// ex. UnRegisterRootDirRequest("/home/ubuntu/rootDir", "password")
func UnRegisterRootDirRequest(DisconnectRootDir string, RootDirPW string) error {

	_, file := filepath.Split(DisconnectRootDir)
	before, after := utils.SplitBeforeAfterRoot(DisconnectRootDir)

	body := types.RegisterRootDirRequest{
		Uuid:            viper.GetViperEnvVariables("UUID"),
		RootDirPassword: RootDirPW,
		BeforePath:      before,
		AfterPath:       after,
	}
	response, err := Conn.SendMessageWithResponse(NOTROOTDIRANYMORE, body.Encode())
	if err != nil {
		return err
	}
	if string(response) == "OK" {
		viper.DeleteViperVariablesByKey("ROOT." + file)
		return nil
		//TODO make clear from fsnotify

	}
	return fmt.Errorf("RegisterRootDir: %s", string(response))
}

// @URL /api/v1/connect/list/remote
// ex. ListRemoteRootDirRequest()
func ShowListRemoteRootDirRequest() {
	resp, err := Conn.SendMessageWithResponse(SHOWREMOTELIST, []byte(""))
	if err != nil {
		log.Println(err)
	}
	log.Println("quics-client :response : " + string(resp))
}

// @URL /api/v1/connect/server
// ex. RegisterClientRequest("password", "S_IP", "S_PORT")
func RegisterClient(ClientPW string, SIp string, SPort string) {
	if SIp != "" {
		viper.WriteViperEnvVariables("QUICS_SERVER_IP", SIp)
	}
	if SPort != "" {
		viper.WriteViperEnvVariables("QUICS_SERVER_PORT", SPort)
	}

	UUID := uuid.New().String()
	viper.WriteViperEnvVariables("UUID", UUID)

	body := types.RegisterClientRequest{
		Uuid:           UUID,
		ClientPassword: ClientPW,
	}
	ClientFirstMessage(CLIENT, body.Encode())
}

// @URL /api/v1/disconnect/server
// ex. DisconnectClientRequest("password")
func DisconnectClientRequest(password string) {
	NotClientAnymoreRequest := types.NotClientAnymoreRequest{
		Uuid:           viper.GetViperEnvVariables("UUID"),
		ClientPassword: password,
	}
	resp, err := Conn.SendMessageWithResponse(NOTCLIENTANYMORE, NotClientAnymoreRequest.Encode())
	if err != nil {
		log.Println(err)
	}
	log.Println("quics-client :response : " + string(resp))
	CloseConnect()
}

// @URL /api/v1/disconnect/root
// ex. DisConnectClientRequest( "/home/ubuntu/rootDir", "password")
func DisconnectRootDirRequest(rootpath string, password string) {
	_, after := utils.SplitBeforeAfterRoot(rootpath)
	notRootDirAnymorRequest := types.NotRootDirAnymorRequest{
		Uuid:            viper.GetViperEnvVariables("UUID"),
		RootPath:        after, // /rootDir
		RootDirPassword: password,
	}
	resp, err := Conn.SendMessageWithResponse(NOTROOTDIRANYMORE, notRootDirAnymorRequest.Encode())
	if err != nil {
		log.Println(err)
	}
	log.Println("quics-client :response : " + string(resp))

}
