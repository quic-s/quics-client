package qclient

import (
	"log"
	"os"

	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func MustSyncRecvHandler(stream *qp.Stream) (*qstypes.MustSyncReq, error) {

	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	req := qstypes.MustSyncReq{}
	req.Decode(data)
	return &req, nil

}

func MustSyncHandler(stream *qp.Stream, UUID string, AfterPath string, LastSyncTimestamp uint64, LastSyncHash string) error {
	bres := qstypes.MustSyncRes{
		UUID:                UUID,
		AfterPath:           AfterPath,
		LatestSyncHash:      LastSyncHash,
		LatestSyncTimestamp: LastSyncTimestamp,
	}

	res, err := bres.Encode()
	if err != nil {
		return err
	}

	err = stream.SendBMessage(res)
	if err != nil {
		return err
	}

	return nil

}

func GiveYouRecvHandler(stream *qp.Stream, path string, Isremoved bool) (*qstypes.GiveYouReq, error) {
	data, fileInfo, fileContent, err := stream.RecvFileBMessage()

	if err != nil {
		log.Println("quics-client: ", err)
		return nil, err
	}

	log.Println("quics-client: ", "file received")
	req := qstypes.GiveYouReq{}
	req.Decode(data)

	err = fileInfo.WriteFileWithInfo(path, fileContent)
	if err != nil {
		return nil, err
	}
	log.Println("quics-client: ", "file saved")

	if Isremoved {
		err = os.Remove(path)
		if err != nil {
			return nil, err
		}
		log.Println("quics-client: ", "file removed")
		return nil, nil
	}

	return &req, nil

}

func GiveYouHandler(stream *qp.Stream, UUID string, AfterPath string, LastSyncTimestamp uint64, LastSyncHash string) error {

	bres := qstypes.GiveYouRes{
		UUID:              UUID,
		AfterPath:         AfterPath,
		LastHash:          LastSyncHash,
		LastSyncTimestamp: LastSyncTimestamp,
	}

	res, err := bres.Encode()
	if err != nil {
		return err
	}

	err = stream.SendBMessage(res)
	if err != nil {
		return err
	}
	return nil
}
