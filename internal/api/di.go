package api

import (
	"github.com/google/wire"
	"github.com/tikivn/ultrago/u_handler"
)

var GraphSet = wire.NewSet(
	u_handler.NewBaseHandler,
	NewEnhanceApiHandler,
)
