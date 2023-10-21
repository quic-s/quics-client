package qclient

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/utils"
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

func GiveYouRecvHandler(stream *qp.Stream, path string, afterPath string, hash string, Isremoved bool) (*qstypes.GiveYouReq, error) {
	data, fileInfo, fileContent, err := stream.RecvFileBMessage()

	if err != nil {
		log.Println("quics-client: ", err)
		return nil, err
	}

	log.Println("quics-client: ", "file received")
	req := qstypes.GiveYouReq{}
	req.Decode(data)

	downloadDir := filepath.Join(utils.GetQuicsDirPath(), "download")

	err = fileInfo.WriteFileWithInfo(filepath.Join(downloadDir, afterPath), fileContent)
	if err != nil {
		return nil, err
	}
	log.Println("quics-client: ", "file downloaded")

	downloadFileInfo, err := os.Stat(filepath.Join(downloadDir, afterPath))
	if err != nil {
		return nil, err
	}

	// if file is removed, then remove file
	if Isremoved {
		err = os.Remove(filepath.Join(downloadDir, afterPath))
		if err != nil {
			return nil, err
		}
		err = os.Remove(path)
		if err != nil {
			return nil, err
		}
		log.Println("quics-client: ", "file removed")
		return nil, nil
	}

	// check hash is correct
	h := utils.MakeHash(afterPath, downloadFileInfo)
	if h != hash {
		os.Remove(filepath.Join(downloadDir, afterPath))
		return nil, errors.New("hash is not correct")
	}

	// copy file to path
	err = utils.CopyFile(filepath.Join(downloadDir, afterPath), path)
	if err != nil {
		return nil, err
	}

	err = os.Remove(filepath.Join(downloadDir, afterPath))
	if err != nil {
		return nil, err
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

func NeedContentRecvHandler(stream *qp.Stream) (*qstypes.NeedContentReq, error) {

	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	req := qstypes.NeedContentReq{}
	req.Decode(data)
	return &req, nil

}

func NeedContentHandler(stream *qp.Stream, UUID string, AfterPath string, LastUpdateTimestamp uint64, LastUpdateHash string) error {

	bres := qstypes.NeedContentRes{
		UUID:              UUID,
		AfterPath:         AfterPath,
		LastUpdateHash:    LastUpdateHash,
		LastUpdateVersion: LastUpdateTimestamp,
	}

	res, err := bres.Encode()
	if err != nil {
		return err
	}

	err = stream.SendFileBMessage(res, AfterPath)
	if err != nil {
		return err
	}
	return nil
}
