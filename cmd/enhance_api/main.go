package main

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tikivn/ultrago/u_graceful"
	"github.com/tikivn/ultrago/u_logger"
	"golang.org/x/sync/errgroup"

	"nimbus-enhance-api/internal/conf"
	"nimbus-enhance-api/internal/repo/gorm"
	"nimbus-enhance-api/pkg/mem_stats"
)

func init() {
	os.Setenv("TZ", "UTC")
	_, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	if conf.Config.IsDebug() {
		u_logger.WithFormatter(logrus.DebugLevel)
	} else {
		u_logger.WithFormatter(logrus.InfoLevel)
	}
}

func main() {
	ctx, logger := u_logger.GetLogger(u_graceful.NewCtx())

	// load config env
	if err := conf.LoadConfig("."); err != nil {
		logger.Fatalf("cannot load config: %v", err)
	}

	// migration code
	if conf.Config.IsMigration() {
		if err := gorm.RunMigration(ctx); err != nil {
			logger.Fatalf(err.Error())
		}
	}

	app, cleanup, err := initApp(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		shutDownErr := app.Stop(context.Background()) // graceful shutdown
		logger.Infof("master is shutdown with err=%v", shutDownErr)
		cleanup() // close connection such as mysql, redis,...
	}()

	// for monitoring memory
	go mem_stats.Monitor(ctx, 30*time.Second)

	eg, childCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return app.Start(childCtx)
	})
	if err = eg.Wait(); err != nil {
		logger.Fatal(err)
	}
}
