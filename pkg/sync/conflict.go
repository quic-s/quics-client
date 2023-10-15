package sync

import (
	"fmt"
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

func GetConflictList() ([]*qstypes.Conflict, error) {

	UUID := badger.GetUUID()
	result := []*qstypes.Conflict{}

	err := Conn.OpenTransaction(qstypes.CONFLICTLIST, func(stream *stream.Stream, transactionName string, transactionID []byte) error {
		cflist, err := qclient.SendAskConflictList(stream, UUID)
		if err != nil {
			return err
		}
		// Create a new slice of pointers to qstypes.Conflict
		conflicts := make([]*qstypes.Conflict, len(cflist.Conflicts))
		// Copy the elements from cflist.Conflicts to the new slice
		for i, c := range cflist.Conflicts {
			conflicts[i] = &c
		}
		result = conflicts
		return nil

	})
	if err != nil {
		return nil, err
	}
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
		for UUID, candidate := range conflictFile.StagingFiles {
			time, err := candidate.File.ModTime.MarshalText()
			if err != nil {
				return "", err
			}

			result += fmt.Sprintf("\t%s, Size > %d ModTime > %s\n", UUID, candidate.File.Size, time)
		}
	}
	return result, nil
}

// @URL /api/v1/conflict/choose
func ChooseOne(path string, Side string) error {

	UUID := badger.GetUUID()
	_, AfterPath := badger.SplitBeforeAfterRoot(path)

	err := Conn.OpenTransaction("CONFLICT", func(stream *stream.Stream, transactionName string, transactionID []byte) error {

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
