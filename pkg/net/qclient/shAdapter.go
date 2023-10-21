package qclient

import (
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func SendLinkShare(stream *qp.Stream, UUID string, AfterPath string, MaxCnt uint) (*qstypes.LinkShareRes, error) {
	breq := qstypes.LinkShareReq{
		UUID:      UUID,
		AfterPath: AfterPath,
		MaxCount:  MaxCnt,
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
	res := qstypes.LinkShareRes{}
	res.Decode(data)
	return &res, nil
}

func SendStopShare(stream *qp.Stream, UUID string, Link string) (*qstypes.StopLinkShareReq, error) {
	breq := qstypes.StopLinkShareReq{
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

	res := qstypes.StopLinkShareRes{}
	res.Decode(data)
	return &res, nil
}
