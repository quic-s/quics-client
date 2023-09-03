package conflict

type TwoOptions struct {
	ServerSideHash          string
	ServerSideSyncTimestamp uint64
	ClientSideHash          string
	ClientSideTimestamp     uint64
}
