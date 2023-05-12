package infra

import (
	"github.com/google/wire"
)

var GraphSet = wire.NewSet(
	NewPostgresSession,
)
