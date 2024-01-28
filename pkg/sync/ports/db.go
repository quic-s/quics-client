package sync

import (
	"github.com/google/wire"
	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/types"
)

type DBDo interface {
	GetUUID() string
	GetSyncMetadata(path string) types.SyncMetadata
	IsSyncMetadataExisted(path string) bool
	GetAllSyncMetadataInRoot(rootpath string) ([]*types.SyncMetadata, error)
	GetAllSyncMetadataAmongRoot() ([]*types.SyncMetadata, error)
	GetRootDirList() []types.RootDir
	GetRootDir(path string) types.RootDir
	GetBeforePathWithAfterPath(afterpath string) string
	SplitBeforeAfterRoot(path string) (string, string)
	AddRootDir(path string) error
	UpdateRootdirToRegistered(path string) error
	DeleteRootDir(path string)
}

// TODO : DB -> APP
type DB struct {
	dbdo DBDo
}

func DBProvider(dbdo DBDo) *DB {
	return &DB{
		dbdo: dbdo,
	}
}

var DBSet = wire.NewSet(
	badger.BadgerProvider,
	DBProvider,
	wire.Bind(new(DBDo), new(*badger.Badger)),
)
