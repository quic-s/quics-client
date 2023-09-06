package history

import (
	"github.com/quic-s/quics-client/pkg/connection"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
	"github.com/quic-s/quics-client/pkg/types"
)
func PleaseFileRequest(version uint64, filepath string) {
	before, after := utils.SplitBeforeAfterRoot(filepath)
	pleaseFile := types.PleaseFile{
		Uuid:          viper.GetViperEnvVariables("UUID"),
		SyncTimestamp: version,
		BeforePath:    before,
		AfterPath:     after,
	}
	connection.Conn.SendMessage(connection.PLEASEFILE, pleaseFile.Encode())
}
