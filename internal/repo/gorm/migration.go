package gorm

import (
	"context"
	"time"

	"github.com/tikivn/ultrago/u_logger"

	"nimbus-enhance-api/internal/conf"
	"nimbus-enhance-api/internal/infra"
)

func RunMigration(ctx context.Context) error {
	if !conf.Config.IsMigration() {
		return nil
	}

	ctx, logger := u_logger.GetLogger(ctx)
	logger.Infof("Run migration")
	db, _, err := infra.NewPostgresSession()
	if err != nil {
		logger.Error(err)
		return err
	}
	// db.DisableForeignKeyConstraintWhenMigrating = true
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal(err)
		return err
	}
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(time.Minute)

	err = db.WithContext(ctx).AutoMigrate(
		&DailyMetaDao{},
	)
	if err != nil {
		logger.Fatal(err)
		return err
	}
	return nil
}
