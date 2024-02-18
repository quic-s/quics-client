//go:build wireinject
// +build wireinject

package qclient

import (
	"github.com/google/wire"
)

var NetSet = wire.NewSet(
	NewQPClient,
	wire.Bind(new(QPPort), new(*QPClient)),
	NewQPC,
)

// func injectDB() (*DB, error) {
// 	panic(wire.Build(DBSet))

// }

func injectQPC() (*QPC, error) {
	panic(wire.Build(NetSet))
}
