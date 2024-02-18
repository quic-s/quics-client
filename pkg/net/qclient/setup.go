package qclient

import (
	qp "github.com/quic-s/quics-protocol"
)

type QPClient struct {
	Conn     *qp.Connection `wire:"-"` // ignore this field when inject
	QPClient *qp.QP
}

func NewQPClient() *QPClient {
	//newClient, err := qp.New(qp.LOG_LEVEL_INFO)
	newClient, err := qp.New(qp.LOG_LEVEL_ERROR)
	if err != nil {
		panic(err)
	}
	return &QPClient{QPClient: newClient}
}

func (qc *QPClient) closeConnnect() {
	qc.Conn.Close()
	qc.Conn = nil
}
