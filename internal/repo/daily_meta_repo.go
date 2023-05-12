package repo

import (
	"context"

	"gorm.io/gorm"

	"nimbus-enhance-api/internal/entity"
	"nimbus-enhance-api/internal/repo/gorm_scope"
)

type DailyMetaRepo interface {
	S() *gorm_scope.DailyMetaScope
	GetOne(ctx context.Context, scopes ...func(db *gorm.DB) *gorm.DB) (*entity.DailyMeta, error)
	GetList(ctx context.Context, scopes ...func(db *gorm.DB) *gorm.DB) ([]*entity.DailyMeta, error)
	Create(ctx context.Context, entity *entity.DailyMeta) error
	Update(ctx context.Context, entity *entity.DailyMeta) error
}
