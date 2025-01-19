package database

import (
	"context"
	"fmt"
	"go-codebase/pkg/logger"
	"go-codebase/pkg/utils"
	"math"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	TransactionContextKey                 ContextTransaction = "postgres:transaction"
	TransactionLockTypeContextKey         ContextTransaction = "postgres:transaction_lock_type"
	DefaultPostgresMaxOpenConnections                        = float64(10)
	DefaultPostgresMaxIdleConnections                        = float64(5)
	DefaultPostgresMaxIdleTimeConnections                    = float64(19)
	PostgresLockTypeForUpdate                                = "UPDATE"
	PostgresLockTypeForShare                                 = "SHARE"
	TracePostgresServiceName                                 = "postgres"
)

func NewPostgres(opts *DBServiceOpts) (IDBService, error) {
	l := logger.WithId(opts.Logger, ContextName, "NewPostgres")

	if opts.PostgresUri == nil {
		l.Debug("Postgres URI is not set, skipping")
		return nil, nil
	}

	if opts.ApplicationName == nil {
		opts.ApplicationName = utils.ToPointer(TracePostgresServiceName)
	} else {
		opts.ApplicationName = utils.ToPointer(fmt.Sprintf("%s_postgres", *opts.ApplicationName))
	}

	// register instrumentation
	// sqltrace.Register("pgx", &stdlib.Driver{}, sqltrace.WithServiceName(*opts.ApplicationName))

	// create connection
	_, err := pgxpool.Connect(context.Background(), *opts.PostgresUri)
	if err != nil {
		l.Error("Cannot open database connection", zap.Error(err))
		return nil, err
	}

	// create logger
	gormOpts := &gorm.Config{
		PrepareStmt: true,
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  *opts.PostgresUri,
		PreferSimpleProtocol: true,
	})

	// create GORM
	db, err := gorm.Open(dialector, gormOpts)
	if err != nil {
		l.Error("Cannot open Gorm", zap.Error(err))
		return nil, err
	}

	// get connection
	sqlDB, err := db.DB()
	if err != nil {
		l.Error("Cannot retrieve database from connection", zap.Error(err))
		return nil, err
	}

	// set max open connection
	sqlDB.SetMaxOpenConns(int(math.Max(float64(opts.PostgresMaxOpenConnections), DefaultPostgresMaxOpenConnections)))
	sqlDB.SetMaxIdleConns(int(math.Max(float64(opts.PostgresMaxIdleConnections), DefaultPostgresMaxIdleConnections)))
	sqlDB.SetConnMaxIdleTime(time.Duration(DefaultPostgresMaxIdleTimeConnections) * time.Second)

	// setup cancellation
	ctx, cancel := context.WithTimeout(context.Background(), PingTimeout)
	defer cancel()

	// ping database
	if err := sqlDB.PingContext(ctx); err != nil {
		l.Error("Database ping timed out")
		return nil, err
	}

	l.Info("Database connected")
	return &DBService{
		Gorm: db,
	}, nil
}

func (db *DBService) Ping() bool {
	l := logger.NewLogger(ContextName, "PingPostgres")
	ctx, cancel := context.WithTimeout(context.Background(), PingTimeout)
	defer cancel()

	_, err := db.Gorm.ConnPool.ExecContext(ctx, "SELECT 1")
	if err != nil {
		l.Error("Cannot check postgres", zap.Error(err))
	}
	return err == nil
}

func (db *DBService) BeginTransaction(ctx context.Context) (context.Context, *gorm.DB) {
	tx := ctx.Value(TransactionContextKey)
	if tx == nil {
		tx = db.Gorm.Begin()
		ctx = context.WithValue(ctx, TransactionContextKey, tx)
	}

	return ctx, tx.(*gorm.DB)
}

func (db *DBService) SetUpdateLockType(ctx context.Context) context.Context {
	return context.WithValue(ctx, TransactionLockTypeContextKey, PostgresLockTypeForUpdate)
}

func (db *DBService) SetShareLockType(ctx context.Context) context.Context {
	return context.WithValue(ctx, TransactionLockTypeContextKey, PostgresLockTypeForShare)
}

func (db *DBService) GetLockType(ctx context.Context) *string {
	lockType := ctx.Value(TransactionLockTypeContextKey)
	if lockType == nil {
		return nil
	}
	return utils.ToPointer(lockType.(string))
}

func (db *DBService) Defer(ctx context.Context) {
	tx := ctx.Value(TransactionContextKey).(*gorm.DB)
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	}

	tx.Commit()
}

func (db *DBService) GetTransaction(ctx context.Context) *gorm.DB {
	tx := ctx.Value(TransactionContextKey)
	if tx == nil {
		tx = db.Gorm
	}

	return tx.(*gorm.DB).WithContext(ctx)
}

func (db *DBService) RollbackTransaction(ctx context.Context) *gorm.DB {
	tx := ctx.Value(TransactionContextKey).(*gorm.DB)
	return tx.Rollback()
}

func (db *DBService) CommitTransaction(ctx context.Context) *gorm.DB {
	tx := ctx.Value(TransactionContextKey).(*gorm.DB)
	return tx.Commit()
}
