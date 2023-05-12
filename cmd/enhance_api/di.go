//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"go-nimeth/internal/app/enhance_api"
)

func initApp(ctx context.Context) (enhance_api.App, func(), error) {
	wire.Build(enhance_api.GraphSet)
	return nil, nil, nil
}
