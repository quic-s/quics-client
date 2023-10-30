package qclient

import (
	"log"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func SendChooseOne(stream *qp.Stream, UUID string, AfterPath string, side string) (*qstypes.PleaseFileRes, error) {

	breq := qstypes.PleaseFileReq{
		UUID:      UUID,
		AfterPath: AfterPath,
		Side:      side,
	}
	req, err := breq.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return nil, err
	}
	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	res := qstypes.PleaseFileRes{}
	res.Decode(data)
	return &res, nil

}

func SendAskConflictList(stream *qp.Stream, UUID string) (*qstypes.AskConflictListRes, error) {
	breq := qstypes.AskConflictListReq{
		UUID: UUID,
	}
	req, err := breq.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return nil, err
	}

	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	res := qstypes.AskConflictListRes{}
	res.Decode(data)
	return &res, nil

}

func SendConflictDownload(stream *qp.Stream, UUID string, AfterPath string) ([]*qstypes.ConflictDownloadReq, error) {

	breq := qstypes.AskStagingNumReq{
		UUID:      UUID,
		AfterPath: AfterPath,
	}

	req, err := breq.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return nil, err
	}

	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	res := qstypes.AskStagingNumRes{}
	res.Decode(data)

	result := []*qstypes.ConflictDownloadReq{}

	if res.ConflictNum == 0 {
		return result, nil
	}

	for i := uint64(0); i < res.ConflictNum; i++ {

		data, fileInfo, fileContent, err := stream.RecvFileBMessage()
		if err != nil {
			continue
		}
		res := qstypes.ConflictDownloadReq{}
		res.Decode(data)

		base := filepath.Base(AfterPath)
		name := "Conflict_" + res.Candidate + "_" + base
		err = fileInfo.WriteFileWithInfo(filepath.Join(utils.GetDownloadDirPath(), name), fileContent)
		if err != nil {
			continue
		}
		log.Println("Conflict files recieved >> ", name)
		result = append(result, &res)
	}
	return result, nil

}
