package app

import (
	"github.com/google/wire"
	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/sync"
)

type App struct {
	db sync.DB
}

// TODO : DB -> APP

var DBSet = wire.NewSet(
	badger.NewBadger,

	wire.Bind(new(sync.DB), new(*badger.Badger)),
)
