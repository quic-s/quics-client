package sync

import (
	"log"
	"reflect"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

// @URL /api/v1/history/rollback
func RollBack(path string, version uint64) error {

	_, AfterPath := badger.SplitBeforeAfterRoot(path)

	err := Conn.OpenTransaction(qstypes.ROLLBACK, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
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
	_, AfterPath := badger.SplitBeforeAfterRoot(path)
	err := Conn.OpenTransaction(qstypes.HISTORYSHOW, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
		historyshowres, err := qclient.SendShowHistory(stream, badger.GetUUID(), AfterPath, cntfromhead)
		if err != nil {
			return err
		}
		if reflect.ValueOf(historyshowres).IsZero() {
			log.Println("dont have history")
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
	_, AfterPath := badger.SplitBeforeAfterRoot(path)

	err := Conn.OpenTransaction(qstypes.HISTORYDOWNLOAD, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
		historydownloadres, err := qclient.SendDownloadHistory(stream, badger.GetUUID(), AfterPath, version)
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
