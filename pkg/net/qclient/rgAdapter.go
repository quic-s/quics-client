package qclient

import (
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func SendClientRegister(stream *qp.Stream, UUID string, ClientPassword string) (qstypes.ClientRegisterRes, error) {
	breq := qstypes.ClientRegisterReq{
		UUID:           UUID,
		ClientPassword: ClientPassword,
	}
	req, err := breq.Encode()
	if err != nil {
		return qstypes.ClientRegisterRes{}, err
	}

	stream.SendBMessage(req)
	bres, err := stream.RecvBMessage()
	if err != nil {
		return qstypes.ClientRegisterRes{}, err
	}

	res := qstypes.ClientRegisterRes{}
	res.Decode(bres)
	return res, nil

}

func SendAskRootList(stream *qp.Stream, UUID string) (*qstypes.AskRootDirRes, error) {
	breq := qstypes.AskRootDirReq{
		UUID: UUID,
	}
	req, err := breq.Encode()
	if err != nil {
		return &qstypes.AskRootDirRes{}, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return &qstypes.AskRootDirRes{}, err
	}
	bres, err := stream.RecvBMessage()
	if err != nil {
		return &qstypes.AskRootDirRes{}, err
	}

	res := qstypes.AskRootDirRes{}
	res.Decode(bres)
	return &res, nil

}

func SendRootDirRegister(stream *qp.Stream, UUID string, RootDirPassword string, BeforePath string, AfterPath string) (qstypes.RootDirRegisterRes, error) {
	breq := qstypes.RootDirRegisterReq{
		UUID:            UUID,
		RootDirPassword: RootDirPassword,
		BeforePath:      BeforePath,
		AfterPath:       AfterPath,
	}
	req, err := breq.Encode()
	if err != nil {
		return qstypes.RootDirRegisterRes{}, err
	}

	err = stream.SendBMessage(req)
	if err != nil {
		return qstypes.RootDirRegisterRes{}, err
	}

	bres, err := stream.RecvBMessage()
	if err != nil {
		return qstypes.RootDirRegisterRes{}, err
	}

	res := qstypes.RootDirRegisterRes{}
	res.Decode(bres)
	return res, nil

}
