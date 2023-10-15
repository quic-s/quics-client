package qclient

import (
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
