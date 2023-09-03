package sync

import (
	qc "github.com/quic-s/quics-client/pkg/quic"
)

func Rescan() {
	qc.ClientMessage(qc.RESCAN, []byte("reset the time"))
}
