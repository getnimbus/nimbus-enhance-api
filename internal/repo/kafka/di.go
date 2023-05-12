package kafka

import (
	"github.com/google/wire"
)

var GraphSet = wire.NewSet(
	NewBlockEventRepo,
)
