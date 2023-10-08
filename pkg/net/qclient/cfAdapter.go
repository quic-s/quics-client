package qclient

import (
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func SendPleaseServerFile(stream *qp.Stream, path string, UUID string, AfterPath string, SelectedTimestamp uint64, NewTimestamp uint64, NewHash string) (*qstypes.PleaseFileRes, error) {

	breq := qstypes.PleaseFileReq{
		UUID:              UUID,
		AfterPath:         AfterPath,
		SelectedTimestamp: SelectedTimestamp,
		NewTimestamp:      NewTimestamp,
		NewHash:           NewHash,
		Side:              "SERVER",
	}

	req, err := breq.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return nil, err
	}
	data, fileinfo, filecontent, err := stream.RecvFileBMessage()
	if err != nil {

		return nil, err
	}
	res := qstypes.PleaseFileRes{}
	res.Decode(data)

	err = fileinfo.WriteFileWithInfo(path, filecontent)
	if err != nil {
		return nil, err
	}
	return &res, nil

}

func SendPleaseLocalFile(stream *qp.Stream, path string, UUID string, AfterPath string, SelectedTimestamp uint64, NewTimstamp uint64, NewHash string) (*qstypes.PleaseFileRes, error) {

	breq := qstypes.PleaseFileReq{
		UUID:              UUID,
		AfterPath:         AfterPath,
		SelectedTimestamp: SelectedTimestamp,
		NewTimestamp:      NewTimstamp,
		NewHash:           NewHash,
		Side:              "LOCAL",
	}
	req, err := breq.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendFileBMessage(req, path)
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
