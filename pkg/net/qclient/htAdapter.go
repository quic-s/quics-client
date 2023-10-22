package qclient

import (
	"path/filepath"
	"strconv"

	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func SendRollBack(stream *qp.Stream, UUID string, AfterPath string, Version uint64) (*qstypes.RollBackRes, error) {
	breq := qstypes.RollBackReq{
		UUID:      UUID,
		AfterPath: AfterPath,
		Version:   Version,
	}

	req, err := breq.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return nil, err
	}

	rollbackRes := &qstypes.RollBackRes{}
	bres, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	rollbackRes.Decode(bres)
	return rollbackRes, nil
}

func SendShowHistory(stream *qp.Stream, UUID string, Afterpath string, CntFromHead uint64) (*qstypes.ShowHistoryRes, error) {

	breq := qstypes.ShowHistoryReq{
		UUID:        UUID,
		AfterPath:   Afterpath,
		CntFromHead: CntFromHead,
	}

	req, err := breq.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return nil, err
	}

	showHistoryRes := &qstypes.ShowHistoryRes{}
	bres, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	showHistoryRes.Decode(bres)
	return showHistoryRes, nil
}

func SendDownloadHistory(stream *qp.Stream, UUID string, AfterPath string, Version uint64) (*qstypes.DownloadHistoryRes, error) {

	breq := qstypes.DownloadHistoryReq{
		UUID:      UUID,
		AfterPath: AfterPath,
		Version:   Version,
	}
	req, err := breq.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return nil, err
	}

	downloadHistoryRes := &qstypes.DownloadHistoryRes{}
	bres, fileInfo, fileContent, err := stream.RecvFileBMessage()
	if err != nil {
		return nil, err
	}

	downloadHistoryRes.Decode(bres)

	base := filepath.Base(AfterPath)
	name := "History_ver" + strconv.Itoa(int(Version)) + "_" + base
	err = fileInfo.WriteFileWithInfo(filepath.Join(utils.GetDownloadDirPath(), name), fileContent)
	if err != nil {
		return nil, err
	}

	return downloadHistoryRes, nil

}
