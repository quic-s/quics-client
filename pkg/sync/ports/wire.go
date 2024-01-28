//go:build wireinject
// +build wireinject

package sync

import (
	"github.com/google/wire"
)

func injectDB() (*DB, error) {
	panic(wire.Build(DBSet))

}
