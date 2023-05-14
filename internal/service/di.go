package service

import (
	"github.com/google/wire"
	"github.com/tikivn/ultrago/u_http_client"
)

var GraphSet = wire.NewSet(
	u_http_client.NewHttpExecutor,
	NewChainService,
)
