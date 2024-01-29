package qclient

import (
	"github.com/quic-s/quics-client/pkg/db/badger"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

type NetPort interface {
}

type QPC struct {
	qpport QPPort
}

func NewQPC(qpport QPPort) *QPC {
	return &QPC{qpport: qpport}
}

type QPPort interface {
	// cfAdpater
	SendChooseOne(stream *qp.Stream, UUID string, AfterPath string, side string) (*qstypes.PleaseFileRes, error)
	SendAskConflictList(stream *qp.Stream, UUID string) (*qstypes.AskConflictListRes, error)
	SendConflictDownload(stream *qp.Stream, UUID string, AfterPath string) ([]*qstypes.ConflictDownloadReq, error)
	// fcAdapter
	ForceSyncRecvHandler(stream *qp.Stream, badger *badger.Badger) (*qstypes.MustSyncReq, string, error)
	ForceSyncHandler(stream *qp.Stream, UUID string, AfterPath string, LastSyncTimestamp uint64, LastSyncHash string) error
	// fsAdapter
	AskAllMetaRecvHandler(stream *qp.Stream) (*qstypes.AskAllMetaReq, error)
	AskAllMetaHandler(stream *qp.Stream, UUID string, syncMetaList []qstypes.SyncMetadata) error
	SendRescan(stream *qp.Stream, UUID string, RootAfterPath []string) (*qstypes.RescanRes, error)
	// htAdapter
	SendRollBack(stream *qp.Stream, UUID string, AfterPath string, Version uint64) (*qstypes.RollBackRes, error)
	SendShowHistory(stream *qp.Stream, UUID string, Afterpath string, CntFromHead uint64) (*qstypes.ShowHistoryRes, error)
	SendDownloadHistory(stream *qp.Stream, UUID string, AfterPath string, Version uint64) (*qstypes.DownloadHistoryRes, error)
	// msAdapter
	MustSyncRecvHandler(stream *qp.Stream) (*qstypes.MustSyncReq, error)
	MustSyncHandler(stream *qp.Stream, UUID string, AfterPath string, LastSyncTimestamp uint64, LastSyncHash string) error
	GiveYouRecvHandler(stream *qp.Stream, path string, afterPath string, hash string, Isremoved bool) (*qstypes.GiveYouReq, error)
	GiveYouHandler(stream *qp.Stream, UUID string, AfterPath string, LastSyncTimestamp uint64, LastSyncHash string) error
	NeedContentRecvHandler(stream *qp.Stream) (*qstypes.NeedContentReq, error)
	NeedContentHandler(stream *qp.Stream, path string, UUID string, AfterPath string, LastUpdateTimestamp uint64, LastUpdateHash string) error
	// psAdapter
	SendPleaseSync(stream *qp.Stream, UUID string, Event string, AfterPath string, LastUpdateTimestamp uint64, LastUpdateHash string, LastSyncHash string, fileMetadata qstypes.FileMetadata) (*qstypes.PleaseSyncRes, error)
	SendPleaseTake(stream *qp.Stream, UUID string, AfterPath string, path string) (*qstypes.PleaseTakeRes, error)
	// rgAdapter
	SendClientRegister(stream *qp.Stream, UUID string, ClientPassword string) (qstypes.ClientRegisterRes, error)
	SendAskRootList(stream *qp.Stream, UUID string) (*qstypes.AskRootDirRes, error)
	SendRootDirRegister(stream *qp.Stream, UUID string, RootDirPassword string, BeforePath string, AfterPath string) (qstypes.RootDirRegisterRes, error)
	SendDisconnectRootDir(stream *qp.Stream, UUID string, AfterPath string) (*qstypes.DisconnectRootDirRes, error)
	SendDisconnectClient(stream *qp.Stream, UUID string) (*qstypes.DisconnectClientRes, error)
	// rsAdapter
	SendPing(stream *qp.Stream, UUID string) (*qstypes.Ping, error)
	// shAdapter
	SendLinkShare(stream *qp.Stream, UUID string, AfterPath string, MaxCnt uint64) (*qstypes.ShareRes, error)
	SendStopShare(stream *qp.Stream, UUID string, Link string) (*qstypes.StopShareRes, error)
}
