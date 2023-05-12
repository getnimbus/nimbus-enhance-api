package gorm

import (
	"context"
	"time"

	"github.com/tikivn/ultrago/u_time"
	"gorm.io/gorm"

	"go-nimeth/internal/entity"
	"go-nimeth/internal/repo"
	"go-nimeth/internal/repo/gorm_scope"
	"go-nimeth/internal/setting"
)

func NewDailyMetaRepo(
	baseRepo *baseRepo,
	s *gorm_scope.DailyMetaScope,
) repo.DailyMetaRepo {
	return &dailyMetaRepo{
		baseRepo: baseRepo,
		s:        s,
	}
}

type dailyMetaRepo struct {
	*baseRepo
	s *gorm_scope.DailyMetaScope
}

func (repo *dailyMetaRepo) S() *gorm_scope.DailyMetaScope {
	return repo.s
}

func (repo *dailyMetaRepo) GetOne(ctx context.Context, scopes ...func(db *gorm.DB) *gorm.DB) (*entity.DailyMeta, error) {
	if len(scopes) == 0 {
		return nil, setting.MissingConditionErr
	}
	db, cancel := repo.WithTimeoutCtx(ctx)
	defer cancel()

	var row DailyMetaDao
	q := db.Model(&DailyMetaDao{}).
		Scopes(scopes...).
		First(&row)
	if err := q.Error; err != nil {
		return nil, err
	}
	return row.toStruct()
}

func (repo *dailyMetaRepo) GetList(ctx context.Context, scopes ...func(db *gorm.DB) *gorm.DB) ([]*entity.DailyMeta, error) {
	if len(scopes) == 0 {
		return nil, setting.MissingConditionErr
	}
	db, cancel := repo.WithTimeoutCtx(ctx)
	defer cancel()

	var rows []*DailyMetaDao
	q := db.Model(&DailyMetaDao{}).
		Scopes(scopes...).
		Find(&rows)
	if err := q.Error; err != nil {
		return nil, err
	}

	res := make([]*entity.DailyMeta, 0, q.RowsAffected)
	for _, row := range rows {
		item, err := row.toStruct()
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (repo *dailyMetaRepo) Create(ctx context.Context, entity *entity.DailyMeta) error {
	row, err := new(DailyMetaDao).fromStruct(entity)
	if err != nil {
		return err
	}
	db, cancel := repo.WithTimeoutCtx(ctx)
	defer cancel()

	q := db.Create(&row)
	return q.Error
}

func (repo *dailyMetaRepo) Update(ctx context.Context, entity *entity.DailyMeta) error {
	row, err := new(DailyMetaDao).fromStruct(entity)
	if err != nil {
		return err
	}
	db, cancel := repo.WithTimeoutCtx(ctx)
	defer cancel()

	q := db.Updates(&row)
	return q.Error
}

type DailyMetaDao struct {
	BaseDao
	CoinName string    `gorm:"column:coin_name;type:varchar(25);not null;index:idx_block_info,priority:2,unique;<-create"`
	Date     time.Time `gorm:"column:date;type:date;not null;index:idx_block_info,priority:1,unique;<-create"`
	From     uint64    `gorm:"column:from;type:bigint;not null;<-create"`
	To       uint64    `gorm:"column:to;type:bigint;not null;<-create"`
	Status   int       `gorm:"column:status;type:int;not null;default:1"`
}

func (dao *DailyMetaDao) TableName() string {
	return "daily_meta"
}

func (dao *DailyMetaDao) fromStruct(item *entity.DailyMeta) (*DailyMetaDao, error) {
	dao.BaseDao = *new(BaseDao).fromEntity(&item.Base)
	dao.CoinName = item.CoinName
	dao.Date = item.Date.ToTime()
	dao.From = item.From
	dao.To = item.To
	dao.Status = int(item.Status)

	return dao, nil
}

func (dao *DailyMetaDao) toStruct() (*entity.DailyMeta, error) {
	return &entity.DailyMeta{
		Base:     *dao.BaseDao.toEntity(),
		CoinName: dao.CoinName,
		Date:     *u_time.FromDateTime(dao.Date),
		From:     dao.From,
		To:       dao.To,
		Status:   entity.DailyMetaStatus(dao.Status),
	}, nil
}
