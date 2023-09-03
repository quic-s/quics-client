package history

import (
	qc "github.com/quic-s/quics-client/pkg/quic"
)

func PleaseFileRequest(version uint64) {
	pleaseFile := PleaseFile{
		SyncTimestamp: version,
	}
	qc.ClientMessage(qc.HISTORY, pleaseFile.Encode())
}
