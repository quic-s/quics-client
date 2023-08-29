package sync

// 폴더 경로 | 마지막 동기화 시점
type SyncStatus struct {
	RootAbsPath string
	LastSyncd   string
	NickName    string //nullable
}

func ShowAllStatus() SyncStatus {
	return SyncStatus{}
}

func ShowStatus(DirForStatus string) SyncStatus {
	return SyncStatus{}
}
