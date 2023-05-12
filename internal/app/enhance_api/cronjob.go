package enhance_api

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/tikivn/ultrago/u_logger"
)

func NewCronjob() Cronjob {
	return &cronjob{
		scheduler: gocron.NewScheduler(time.UTC),
	}
}

type Cronjob interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type cronjob struct {
	scheduler *gocron.Scheduler
}

func (c *cronjob) Start(ctx context.Context) error {
	ctx, logger := u_logger.GetLogger(ctx)

	j, err := c.scheduler.SingletonMode().
		Cron("*/15 * * * *").
		Do(func() {
			// TODO: write function to do cronjob
			//logger.Infof("start updating new token")
			//err := a.tokenUpdateSvc.SyncTokens(ctx, carbon.Now(carbon.UTC).StartOfDay().ToStdTime())
			//if err != nil {
			//	logger.Errorf("failed to sync new tokens: %v", err)
			//	return
			//}
			//logger.Infof("done updating new token")
		})
	if err != nil {
		logger.Errorf("failed to registered cronjob %#v: %v", j, err)
		return err
	}

	logger.Info("start cronjob scheduler")
	c.scheduler.StartBlocking()

	return nil
}

func (c *cronjob) Stop(ctx context.Context) error {
	ctx, logger := u_logger.GetLogger(ctx)
	logger.Infof("stop cronjob scheduler")
	c.scheduler.Stop()
	return nil
}
