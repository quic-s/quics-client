package sync

import (
	"github.com/quic-s/quics-client/pkg/types"
)

type Repository interface {
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
