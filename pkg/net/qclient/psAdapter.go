package qclient

import (
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func SendFileMeta(stream *qp.Stream, UUID string, AfterPath string) (*qstypes.PleaseFileMetaRes, error) {

	breq := qstypes.PleaseFileMetaReq{
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
	res := qstypes.PleaseFileMetaRes{}
	res.Decode(data)
	return &res, nil

}

func SendPleaseSync(stream *qp.Stream, UUID string, Event string, BeforePath string, AfterPath string, LastUpdateTimestamp uint64, LastUpdateHash string) (*qstypes.PleaseSyncRes, error) {

	breq := qstypes.PleaseSyncReq{
		UUID:                UUID,
		Event:               Event,
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: LastUpdateTimestamp,
		LastUpdateHash:      LastUpdateHash,
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
	res := qstypes.PleaseSyncRes{}
	res.Decode(data)
	return &res, nil

}

func SendPleaseTake(stream *qp.Stream, UUID string, AfterPath string, path string) (*qstypes.PleaseTakeRes, error) {

	breq := qstypes.PleaseTakeReq{
		UUID:      UUID,
		AfterPath: AfterPath,
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
	res := qstypes.PleaseTakeRes{}
	res.Decode(data)
	return &res, nil

}
