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

		convey.Convey("TestChainService_GetLatestBlock", func() {
			res, err := svc.GetLatestBlock(ctx, BSC)
			convey.So(err, convey.ShouldBeNil)
			data, _ := json.Marshal(res)
			fmt.Println(string(data))
		})

		convey.FocusConvey("TestChainService_SearchTransactionHash", func() {
			res, err := svc.SearchTransactionHash(ctx, "0xc92daec22a20426373626b0fee73cc168b50697f30538062687b4d02bc5a52f6")
			convey.So(err, convey.ShouldBeNil)
			data, _ := json.Marshal(res)
			fmt.Println(string(data))
		})
	})
}
