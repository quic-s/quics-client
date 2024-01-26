//go:build wireinject
// +build wireinject

package defaultdb

import "github.com/google/wire"

// func ProvideCRUD() (*CRUD, error) {
// 	wire.Build(CRUDProvider, wire.Bind(new(CRUD), new(*CRUD)))
// 	return &CRUD{}, nil
// }

var set = wire.NewSet(
	DBProvider,
	Badger.BadgerProvider,
	wire.Bind(new(CRUD), new(*Badger.Badger)),
)

func injectDB() (*DB, error) {
	panic(wire.Build(set))

}
