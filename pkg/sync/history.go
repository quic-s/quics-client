package sync

import (
	"reflect"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

// @URL /api/v1/history/rollback
func RollBack(path string, version uint64) error {

	_, AfterPath := badger.SplitBeforeAfterRoot(path)

	err := QPClient.RecvTransactionHandleFunc(qstypes.ROLLBACK, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		rollbackres, err := qclient.SendRollBack(stream, badger.GetUUID(), AfterPath, version)
		if err != nil {
			return err
		}
		if reflect.ValueOf(rollbackres).IsZero() {
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// @URL /api/v1/history/show
func HistoryShow(path string, cntfromhead uint64) ([]qstypes.FileHistory, error) {
	historyShowRes := []qstypes.FileHistory{}
	err := QPClient.RecvTransactionHandleFunc(qstypes.HISTORYSHOW, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		historyshowres, err := qclient.SendShowHistory(stream, badger.GetUUID(), path, cntfromhead)
		if err != nil {
			return err
		}
		if reflect.ValueOf(historyshowres).IsZero() {
			return nil
		}
		historyShowRes = historyshowres.History
		return nil
	})
	if err != nil {
		return nil, err
	}
	return historyShowRes, nil
}

// @URL /api/v1/history/download
func HistoryDownload(path string, version uint64) error {

	err := QPClient.RecvTransactionHandleFunc(qstypes.HISTORYDOWNLOAD, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		historydownloadres, err := qclient.SendDownloadHistory(stream, badger.GetUUID(), path, version)
		if err != nil {
			return err
		}
		if reflect.ValueOf(historydownloadres).IsZero() {
			return nil
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
