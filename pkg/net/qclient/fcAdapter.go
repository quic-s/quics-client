package qclient

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func ForceSyncRecvHandler(stream *qp.Stream) (*qstypes.MustSyncReq, string, error) {

	data, fileInfo, fileContent, err := stream.RecvFileBMessage()

	if err != nil {
		log.Println("", err)
		return nil, "", err
	}

	log.Println("quics-client: ", "file received")
	req := qstypes.MustSyncReq{}
	req.Decode(data)

	BeforePath := badger.GetBeforePathWithAfterPath(req.AfterPath)
	path := filepath.Join(BeforePath, req.AfterPath)

	downloadDir := filepath.Join(utils.GetQuicsDirPath(), "download")

	err = fileInfo.WriteFileWithInfo(filepath.Join(downloadDir, req.AfterPath), fileContent)
	if err != nil {
		return nil, "", err
	}
	log.Println("quics-client: ", "file downloaded")

	downloadFileInfo, err := os.Stat(filepath.Join(downloadDir, req.AfterPath))
	if err != nil {
		return nil, "", err
	}

	// check hash is correct
	h := utils.MakeHash(req.AfterPath, downloadFileInfo)
	if h != req.LatestHash {
		os.Remove(filepath.Join(downloadDir, req.AfterPath))
		return nil, "", errors.New("hash is not correct")
	}

	// copy file to path
	err = utils.CopyFile(filepath.Join(downloadDir, req.AfterPath), path)
	if err != nil {
		return nil, "", err
	}

	err = os.Remove(filepath.Join(downloadDir, req.AfterPath))
	if err != nil {
		return nil, "", err
	}

	// err = fileInfo.WriteFileWithInfo(path, fileContent)
	// if err != nil {
	// 	return nil, "", err
	// }
	log.Println("quics-client: ", "file saved")

	return &req, BeforePath, nil

}

func ForceSyncHandler(stream *qp.Stream, UUID string, AfterPath string, LastSyncTimestamp uint64, LastSyncHash string) error {
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
