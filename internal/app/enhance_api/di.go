package enhance_api

import (
	"github.com/google/wire"

	"nimbus-enhance-api/internal/api"
	"nimbus-enhance-api/internal/infra"
	"nimbus-enhance-api/internal/service"
)

var deps = wire.NewSet(
	infra.GraphSet,
	//gorm_scope.GraphSet,
	//gorm.GraphSet,
	service.GraphSet,
	api.GraphSet,
)

var GraphSet = wire.NewSet(
	deps,
	NewHttpServer,
	NewCronjob,
	NewApp,
)
