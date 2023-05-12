package enhance_api

import (
	"github.com/google/wire"
	"go-nimeth/internal/api"

	"go-nimeth/internal/infra"
	"go-nimeth/internal/repo/gorm"
	"go-nimeth/internal/repo/gorm_scope"
	"go-nimeth/internal/service"
)

var deps = wire.NewSet(
	infra.GraphSet,
	gorm_scope.GraphSet,
	gorm.GraphSet,
	service.GraphSet,
	api.GraphSet,
)

var GraphSet = wire.NewSet(
	deps,
	NewHttpServer,
	NewApp,
)
