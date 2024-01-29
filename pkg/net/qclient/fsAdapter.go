package qclient

import (
	"github.com/quic-s/quics-protocol/pkg/stream"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func (qc *QPClient) AskAllMetaRecvHandler(stream *stream.Stream) (*qstypes.AskAllMetaReq, error) {
	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}

	req := qstypes.AskAllMetaReq{}
	req.Decode(data)

	return &req, nil

}

func (qc *QPClient) AskAllMetaHandler(stream *stream.Stream, UUID string, syncMetaList []qstypes.SyncMetadata) error {

	bres := qstypes.AskAllMetaRes{
		UUID:         UUID,
		SyncMetaList: syncMetaList,
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

func (qc *QPClient) SendRescan(stream *stream.Stream, UUID string, RootAfterPath []string) (*qstypes.RescanRes, error) {
	bres := qstypes.RescanReq{
		UUID:          UUID,
		RootAfterPath: RootAfterPath,
	}
	res, err := bres.Encode()
	if err != nil {
		return nil, err
	}
	err = stream.SendBMessage(res)
	if err != nil {
		return nil, err
	}

	data, err := stream.RecvBMessage()
	if err != nil {
		return nil, err
	}
	rescanRes := qstypes.RescanRes{}
	rescanRes.Decode(data)
	return &rescanRes, nil
}
