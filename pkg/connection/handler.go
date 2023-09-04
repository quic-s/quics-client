package connection

import (
	qc "github.com/quic-s/quics-client/pkg/quic"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
)

func RegisterQuicsClient() {

	body := RegisterClientRequest{Ip: utils.GetIp()}
	qc.ClientMessage(qc.CLIENT, body.Encode())
}

// TODO password to crypto
func RegisterRootDir(filepath string, password string) {
	before, after := utils.SplitBeforeAfterRoot(filepath)
	body := RegisterRootDirRequest{
		Uuid:       viper.GetViperEnvVariables("UUID"),
		Password:   password,
		BeforePath: before,
		AfterPath:  after,
	}
	qc.ClientMessage(qc.ROOTDIR, body.Encode())

}
