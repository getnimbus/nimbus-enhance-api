package gorm

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-nimeth/internal/repo"
	"go-nimeth/internal/setting"
)

type contextKey string

const (
	txContextKey contextKey = "DBTX"
)

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) *baseRepo {
	return &baseRepo{db: db}
}

func NewTransactor(baseRepo *baseRepo) repo.Transactor {
	return baseRepo
}

// beware using Transactor will slow down performance because of type assertion db instance from context.
// https://stackoverflow.com/questions/28024884/does-a-type-assertion-type-switch-have-bad-performance-is-slow-in-go
func (repo *baseRepo) getDB(ctx context.Context) *gorm.DB {
	tx := ctx.Value(txContextKey)
	if tx != nil {
		return tx.(*gorm.DB)
	}
	return repo.db
}

// BeginTx returns new context as transaction context => use this context inside service.
//
// If you want to start another transaction context in same function then pass the parent
// context as parameter, not the old transaction context.
func (repo *baseRepo) BeginTx(ctx context.Context) (newCtx context.Context, err error) {
	if ctx.Value(txContextKey) != nil {
		err = setting.TransactionInProgressErr
		return
	}
	tx := repo.db.WithContext(ctx).Begin()
	newCtx = context.WithValue(ctx, txContextKey, tx)
	err = tx.Error
	return
}

func (repo *baseRepo) CommitTx(ctx context.Context) error {
	if ctx.Value(txContextKey) == nil {
		return setting.TransactionNotStartedErr
	}
	return repo.getDB(ctx).Commit().Error
}

func (repo *baseRepo) RollbackTx(ctx context.Context) error {
	if ctx.Value(txContextKey) == nil {
		return setting.TransactionNotStartedErr
	}
	return repo.getDB(ctx).Rollback().Error
}

// WithTimeoutCtx use context to get db transaction if existed
func (repo *baseRepo) WithTimeoutCtx(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctxTimeout, cancel := context.WithTimeout(ctx, setting.QueryTimeout*time.Second)

	return repo.getDB(ctx).WithContext(ctxTimeout), cancel
}

func (repo *baseRepo) WithBackgroundTimeout() (*gorm.DB, context.CancelFunc) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), setting.QueryTimeout*time.Second)

	return repo.db.WithContext(ctxTimeout), cancel
}

func (repo *baseRepo) WithTransactionTimeout(ctx context.Context, tx *gorm.DB) (*gorm.DB, context.CancelFunc) {
	ctxTimeout, cancel := context.WithTimeout(ctx, setting.QueryTimeout*time.Second)
	ctxTx := ctx.Value(txContextKey)
	if ctxTx != nil {
		return ctxTx.(*gorm.DB).WithContext(ctxTimeout), cancel
	}
	return tx.WithContext(ctxTimeout), cancel
}
