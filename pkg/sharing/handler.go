package sharing

import (
	"log"

	"github.com/quic-s/quics-client/pkg/badger"
	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-client/pkg/utils"
)

func GetDownloadLink(filepath string, maxCnt uint32) string {

	uuid, err := badger.View("uuid")
	if err != nil {
		return ""
	}

	beforePath, afterPath := utils.SplitBeforeAfterRoot(filepath)
	fileDownloadRequest := types.FileDownloadRequest{
		Uuid:       string(uuid),
		BeforePath: beforePath,
		AfterPath:  afterPath,
		MaxCnt:     maxCnt,
	}
	log.Println(fileDownloadRequest)
	//TODO can get return message
	//connection.ClientMessage(connection.SHARING, fileDownloadRequest.Encode())

	return ""

}
