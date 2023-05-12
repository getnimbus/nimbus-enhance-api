package entity

import (
	"fmt"

	"go-nimeth/internal/setting"

	"github.com/tikivn/ultrago/u_time"
	"github.com/tikivn/ultrago/u_validator"
)

type DailyMetaStatus int

const (
	DailyMetaStatus_NOT_READY  DailyMetaStatus = setting.StatusNotReady // default status
	DailyMetaStatus_PROCESSING DailyMetaStatus = setting.StatusProcessing
	DailyMetaStatus_DONE       DailyMetaStatus = setting.StatusDone
	DailyMetaStatus_FAIL       DailyMetaStatus = setting.StatusFail
)

type DailyMeta struct {
	Base
	CoinName string          `json:"coin_name" validate:"required"`
	Date     u_time.Date     `json:"date" validate:"required"`
	From     uint64          `json:"from" validate:"omitempty,gte=0"`
	To       uint64          `json:"to" validate:"required,gt=0"`
	Status   DailyMetaStatus `json:"status" validate:"required"`
}

func (d *DailyMeta) Validate() error {
	if err := u_validator.Struct(d); err != nil {
		return err
	}

	if d.From >= d.To {
		return fmt.Errorf("invalid block number range from block %v -> %v", d.From, d.To)
	}

	switch d.Status {
	case DailyMetaStatus_NOT_READY,
		DailyMetaStatus_PROCESSING,
		DailyMetaStatus_DONE,
		DailyMetaStatus_FAIL:
	default:
		return fmt.Errorf("invalid daily meta status %v", d.Status)
	}

	return nil
}
