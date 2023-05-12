package entity

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"github.com/tikivn/ultrago/u_time"

	"go-nimeth/internal/setting"
)

func TestDailyMetaEntity(t *testing.T) {
	convey.Convey("TestDailyMetaEntity", t, func() {

		convey.Convey("Empty block number to", func() {
			dailyMeta := &DailyMeta{
				Base: Base{
					ID: "ac3b1914-a122-4fa9-bccc-2db44facc7d6",
				},
				Date:     *u_time.FromDateTime(time.Now()),
				CoinName: setting.EthCoin,
				From:     0,
				To:       0,
				Status:   DailyMetaStatus_NOT_READY,
			}

			err := dailyMeta.Validate()
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("Invalid block number range from >= to", func() {
			dailyMeta := &DailyMeta{
				Base: Base{
					ID: "ac3b1914-a122-4fa9-bccc-2db44facc7d6",
				},
				Date:     *u_time.FromDateTime(time.Now()),
				CoinName: setting.EthCoin,
				From:     100,
				To:       90,
				Status:   DailyMetaStatus_NOT_READY,
			}

			err := dailyMeta.Validate()
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("Success case", func() {
			dailyMeta := &DailyMeta{
				Base: Base{
					ID: "ac3b1914-a122-4fa9-bccc-2db44facc7d6",
				},
				Date:     *u_time.FromDateTime(time.Now()),
				CoinName: setting.EthCoin,
				From:     1,
				To:       100,
				Status:   DailyMetaStatus_NOT_READY,
			}

			err := dailyMeta.Validate()
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
