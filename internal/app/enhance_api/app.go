package enhance_api

import (
	"context"
	"fmt"
	"net"

	"github.com/tikivn/ultrago/u_graceful"
	"github.com/tikivn/ultrago/u_logger"
	"golang.org/x/sync/errgroup"

	"nimbus-enhance-api/internal/setting"
)

func NewApp(
	httpServer *HttpServer,
) App {
	return &app{
		httpServer: httpServer,
	}
}

type App interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type app struct {
	httpServer *HttpServer
}

func (a *app) Start(ctx context.Context) error {
	ctx, logger := u_logger.GetLogger(ctx)
	httpLis, err := net.Listen("tcp", setting.HttpPort)
	if err != nil {
		logger.Fatal("failed to listen http port %s: %v", setting.HttpPort, err)
		return err
	}

	eg, childCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return u_graceful.BlockListen(childCtx, func() error {
			logger.Infof("start listening http request on %s", setting.HttpPort)
			if err := a.httpServer.Serve(httpLis); err != nil {
				return fmt.Errorf("failed to serve http: %v", err)
			}
			return nil
		})
	})

	logger.Info("server started!")
	return eg.Wait()
}

func (a *app) Stop(ctx context.Context) error {
	ctx, logger := u_logger.GetLogger(ctx)
	logger.Infof("stop listening http request on %s", setting.HttpPort)
	err := a.httpServer.Shutdown(ctx)
	return err
}
