package gorm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-nimeth/internal/entity"
)

type BaseDao struct {
	ID        string `gorm:"column:id;type:varchar(36);primaryKey;not null;<-:create"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli;<-:create"`
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (dao *BaseDao) BeforeCreate(db *gorm.DB) error {
	if dao.ID == "" {
		dao.ID = uuid.New().String()
	}
	return nil
}

func (dao *BaseDao) fromEntity(item *entity.Base) *BaseDao {
	dao.ID = item.ID
	dao.CreatedAt = item.CreatedAt
	dao.UpdatedAt = item.UpdatedAt
	return dao
}

func (dao *BaseDao) toEntity() *entity.Base {
	return &entity.Base{
		ID:        dao.ID,
		CreatedAt: dao.CreatedAt,
		UpdatedAt: dao.UpdatedAt,
	}
}
