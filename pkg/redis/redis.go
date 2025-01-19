package redis

import (
	"context"
	"go-codebase/infrastructure/logger"

	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
)

func NewRedis(opts *DBServiceOpts) (IDBService, error) {
	l := logger.WithId(opts.Logger, ContextName, "NewRedis")

	if opts.RedisUri == nil {
		l.Debug("Redis URI is not set, skipping")
		return nil, nil
	}

	// create redis client
	opt, _ := redis.ParseURL(*opts.RedisUri)
	cl := redis.NewClient(opt)

	// setup cancellation
	ctxTimeout, cancel := context.WithTimeout(context.Background(), PingTimeout)
	defer cancel()

	// ping redis
	if res := cl.Ping(ctxTimeout); res.Err() != nil {
		l.Error("Cannot ping redis")
		return nil, res.Err()
	}

	return &DBService{
		Redis: cl,
	}, nil
}

func (db *DBService) PingRedis() bool {
	l := logger.NewLogger(ContextName, "PingRedis")
	ctx, cancel := context.WithTimeout(context.Background(), PingTimeout)
	defer cancel()

	res := db.Redis.Ping(ctx)
	if res.Err() != nil {
		l.Error("Cannot check redis", zap.Error(res.Err()))
	}
	return res.Err() == nil
}

func (db *DBService) PingRedisWithError() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), PingTimeout)
	defer cancel()

	res := db.Redis.Ping(ctx)
	return res.Err() == nil, res.Err()
}
