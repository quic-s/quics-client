package sync

import "github.com/quic-s/quics-client/pkg/connection"

// @URL /api/v1/rescan
// ex) Rescan()
func Rescan() {
	connection.Conn.SendMessage(connection.RESCAN, []byte(""))
}
