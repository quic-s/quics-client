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
func SendDisconnectRootDir(stream *qp.Stream, UUID string, AfterPath string) (qstypes.DisconnectRootDirRes, error) {
	breq := qstypes.DisconnectRootDirReq{
		UUID:      UUID,
		AfterPath: AfterPath,
	}

	req, err := breq.Encode()
	if err != nil {
		return qstypes.DisconnectRootDirRes{}, err
	}

	err = stream.SendBMessage(req)
	if err != nil {

		return qstypes.DisconnectRootDirRes{}, err
	}

	bres, err := stream.RecvBMessage()
	if err != nil {

		return qstypes.DisconnectRootDirRes{}, err
	}

	res := qstypes.DisconnectRootDirRes{}
	res.Decode(bres)

	return res, nil

}

func SendDisconnectClient(stream *qp.Stream, UUID string) (qstypes.DisconnectClientRes, error) {
	breq := qstypes.DisconnectClientReq{
		UUID: UUID,
	}

	req, err := breq.Encode()
	if err != nil {
		return qstypes.DisconnectClientRes{}, err
	}
	err = stream.SendBMessage(req)
	if err != nil {
		return qstypes.DisconnectClientRes{}, err
	}

	bres, err := stream.RecvBMessage()
	if err != nil {
		return qstypes.DisconnectClientRes{}, err
	}

	res := qstypes.DisconnectClientRes{}
	res.Decode(bres)
	return res, nil
}
