//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"nimbus-enhance-api/internal/app/enhance_api"
)

func initApp(ctx context.Context) (enhance_api.App, func(), error) {
	wire.Build(enhance_api.GraphSet)
	return nil, nil, nil
}
