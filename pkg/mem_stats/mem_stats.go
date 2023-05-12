package mem_stats

import (
	"context"
	"runtime"
	"time"

	"github.com/tikivn/ultrago/u_logger"
)

func Monitor(ctx context.Context, delay time.Duration) {
	ctx, logger := u_logger.GetLogger(ctx)
	ticker := time.NewTicker(delay)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			// For info on each, see: https://golang.org/pkg/runtime/#MemStats
			logger.Infof("Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v",
				bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
		}
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
