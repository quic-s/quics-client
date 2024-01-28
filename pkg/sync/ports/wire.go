//go:build wireinject
// +build wireinject

package sync

import "github.com/google/wire"

var set = wire.NewSet(
	DBProvider,
	Badger.BadgerProvider,
	wire.Bind(new(CRUD), new(*Badger.Badger)),
)

func injectDB() (*DB, error) {
	panic(wire.Build(set))

}
