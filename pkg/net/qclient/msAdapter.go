package qclient

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func (qc *QPClient) MustSyncRecvHandler(stream *qp.Stream) (*qstypes.MustSyncReq, error) {

	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	req := qstypes.MustSyncReq{}
	req.Decode(data)
	return &req, nil

}

func (qc *QPClient) MustSyncHandler(stream *qp.Stream, UUID string, AfterPath string, LastSyncTimestamp uint64, LastSyncHash string) error {
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

func (qc *QPClient) GiveYouRecvHandler(stream *qp.Stream, path string, afterPath string, hash string, Isremoved bool) (*qstypes.GiveYouReq, error) {
	data, fileInfo, fileContent, err := stream.RecvFileBMessage()

	if err != nil {
		log.Println("quics-client: ", err)
		return nil, err
	}

	log.Println("quics-client: ", "file received")
	req := qstypes.GiveYouReq{}
	req.Decode(data)

	tempDir := utils.GetQuicsTempDirPath()

	err = fileInfo.WriteFileWithInfo(filepath.Join(tempDir, afterPath), fileContent)
	if err != nil {
		return nil, err
	}
	log.Println("quics-client: ", "file downloaded")

	tempFileInfo, err := os.Stat(filepath.Join(tempDir, afterPath))
	if err != nil {
		return nil, err
	}

	// if file is removed, then remove file
	if Isremoved {
		err = os.Remove(filepath.Join(tempDir, afterPath))
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		err = os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		// If Case in Dir is empty
		for dirPath, _ := filepath.Split(path); dirPath[:len(dirPath)-1] != filepath.Join(path[:len(path)-len(afterPath)], strings.Split(afterPath, "/")[1]); dirPath, _ = filepath.Split(dirPath[:len(dirPath)-1]) {

			dir, err := os.Open(dirPath)
			if err != nil && !os.IsNotExist(err) {
				log.Println("quics err: ", err)
				return nil, err
			} else if os.IsNotExist(err) {
				continue
			}

			// Delete directory when it is empty
			files, err := dir.Readdir(-1)
			if err != nil {
				return nil, err
			}
			if len(files) == 0 {
				os.Remove(dirPath)
			}
			dir.Close()
		}
		log.Println("quics-client: ", "file removed")
		return nil, nil
	}

	// check hash is correct
	h := utils.MakeHash(afterPath, tempFileInfo)
	if h != hash {
		os.Remove(filepath.Join(tempDir, afterPath))
		return nil, errors.New("hash is not correct")
	}

	// copy file to path
	err = utils.CopyFile(filepath.Join(tempDir, afterPath), path)
	if err != nil {
		return nil, err
	}

	err = os.Remove(filepath.Join(tempDir, afterPath))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return &req, nil
}

func (qc *QPClient) GiveYouHandler(stream *qp.Stream, UUID string, AfterPath string, LastSyncTimestamp uint64, LastSyncHash string) error {

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

func (qc *QPClient) NeedContentRecvHandler(stream *qp.Stream) (*qstypes.NeedContentReq, error) {

	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	req := qstypes.NeedContentReq{}
	req.Decode(data)
	return &req, nil

}

func (qc *QPClient) NeedContentHandler(stream *qp.Stream, path string, UUID string, AfterPath string, LastUpdateTimestamp uint64, LastUpdateHash string) error {

	bres := qstypes.NeedContentRes{
		UUID:                UUID,
		AfterPath:           AfterPath,
		LastUpdateHash:      LastUpdateHash,
		LastUpdateTimestamp: LastUpdateTimestamp,
	}

	res, err := bres.Encode()
	if err != nil {
		return err
	}

	err = stream.SendFileBMessage(res, path)
	if err != nil {
		return err
	}
	return nil
}
