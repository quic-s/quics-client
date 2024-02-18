package qclient

import (
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func (qc *QPClient) SendLinkShare(stream *qp.Stream, UUID string, AfterPath string, MaxCnt uint64) (*qstypes.ShareRes, error) {
	breq := qstypes.ShareReq{
		UUID:      UUID,
		AfterPath: AfterPath,
		MaxCnt:    MaxCnt,
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
	res := qstypes.ShareRes{}
	res.Decode(data)
	return &res, nil
}

func (qc *QPClient) SendStopShare(stream *qp.Stream, UUID string, Link string) (*qstypes.StopShareRes, error) {
	breq := qstypes.StopShareReq{
		UUID: UUID,
		Link: Link,
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

	res := qstypes.StopShareRes{}
	res.Decode(data)
	return &res, nil
}
