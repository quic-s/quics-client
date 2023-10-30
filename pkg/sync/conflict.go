package sync

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-protocol/pkg/stream"
	qstypes "github.com/quic-s/quics/pkg/types"
)

const (
	SERVER = "SERVER"
	LOCAL  = "LOCAL"
)

// @URL /api/v1/conflict/download
func ConflictDownload(path string) error {
	_, Afterpath := badger.SplitBeforeAfterRoot(path)

	err := Conn.OpenTransaction(qstypes.CONFLICTDOWNLOAD, func(stream *stream.Stream, transactionName string, transactionID []byte) error {
		log.Println("quics-client : [CONFLICTDOWNLOAD] transaction start")
		UUID := badger.GetUUID()
		res, err := qclient.SendConflictDownload(stream, UUID, Afterpath)
		if err != nil {
			return err
		}
		if reflect.ValueOf(res).IsZero() {
			return fmt.Errorf("cannot download conflict file")
		}

		return nil
	})
	if err != nil {
		return err
	}
	log.Println("quics-client : [CONFLICTDOWNLOAD] transaction success")
	return nil

}

// @URL /api/v1/conflict/list
func GetConflictList() ([]qstypes.Conflict, error) {

	UUID := badger.GetUUID()
	result := []qstypes.Conflict{}

	err := Conn.OpenTransaction(qstypes.CONFLICTLIST, func(stream *stream.Stream, transactionName string, transactionID []byte) error {
		log.Println("quics-client : [CONFLICTLIST] transaction start")
		cflist, err := qclient.SendAskConflictList(stream, UUID)
		if err != nil {
			return err
		}
		result = cflist.Conflicts
		return nil

	})
	if err != nil {
		return nil, err
	}
	log.Println("quics-client : [CONFLICTLIST] transaction success")
	return result, nil

}

// @URL /api/v1/conflict/list
func PrintCFOptions() (string, error) {
	cflist, err := GetConflictList()
	if err != nil {
		return "", err
	}
	result := "\n\n"

	for i, conflictFile := range cflist {
		afterPath := conflictFile.AfterPath
		beforePath := badger.GetBeforePathWithAfterPath(afterPath)
		result += fmt.Sprintf("%d. %s\n", i+1, filepath.Join(beforePath, afterPath))
		y := 1
		for UUID, candidate := range conflictFile.StagingFiles {
			time, err := candidate.File.ModTime.MarshalText()
			if err != nil {
				return "", err
			}

			result += fmt.Sprintf("\t(%d) Candidate > %s \n\t\tSize > %d ModTime > %s\n", y, UUID, candidate.File.Size, time)
			y++
		}
	}
	return result, nil
}

// @URL /api/v1/conflict/choose
func ChooseOne(path string, Side string) error {

	UUID := badger.GetUUID()
	_, AfterPath := badger.SplitBeforeAfterRoot(path)

	err := Conn.OpenTransaction(qstypes.CHOOSEONE, func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		// Send ChooseOne
		res, err := qclient.SendChooseOne(stream, UUID, AfterPath, Side)
		if err != nil {
			return err
		}
		if reflect.ValueOf(res).IsZero() {
			return fmt.Errorf("cannot choose one")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[CONFLICT] ", err)
	}
	return nil
}
