package qclient

import (
	"log"

	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func SendPing(stream *qp.Stream, UUID string) (*qstypes.Ping, error) {

	req := qstypes.Ping{
		UUID: UUID,
	}

	breq, err := req.Encode()
	if err != nil {
		return nil, err
	}

	err = stream.SendBMessage(breq)
	if err != nil {
		log.Println("error occurred when send msg ; ", err)
		return nil, err
	}

	bres, err := stream.RecvBMessage()
	if err != nil {
		log.Println("error occurred when recv msg ; ", err)
		return nil, err
	}

	result := qstypes.Ping{}
	err = result.Decode(bres)
	if err != nil {
		return nil, err
	}

	return &result, nil

}
