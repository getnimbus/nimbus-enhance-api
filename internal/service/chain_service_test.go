package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestChainService(t *testing.T) {
	ctx := context.Background()

	convey.FocusConvey("TestChainService", t, func() {
		svc := NewChainService()

		convey.FocusConvey("TestChainService_GetLatestBlock", func() {
			res, err := svc.GetLatestBlock(ctx, ArbitrumNitro)
			convey.So(err, convey.ShouldBeNil)
			data, _ := json.Marshal(res)
			fmt.Println(string(data))
		})

		convey.Convey("TestChainService_SearchTransactionHash", func() {
			res, err := svc.SearchTransactionHash(ctx, "")
			convey.So(err, convey.ShouldBeNil)
			data, _ := json.Marshal(res)
			fmt.Println(string(data))
		})
	})
}
