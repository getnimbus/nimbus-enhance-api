package gorm_scope

import (
	"github.com/tikivn/ultrago/u_time"
	"gorm.io/gorm"
)

type DailyMetaScope struct {
	*base
}

func NewDailyMeta(b *base) *DailyMetaScope {
	return &DailyMetaScope{base: b}
}

func (s *DailyMetaScope) FilterDate(date u_time.Date) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("date = ?", date.ToString())
	}
}

func (s *DailyMetaScope) FilterCoinName(coinName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("coin_name = ?", coinName)
	}
}
