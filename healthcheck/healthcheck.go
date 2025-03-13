package healthcheck

import (
	"context"
	"time"

	"go-base-utils/constants"
	"go-base-utils/http/middlewares"
	"go-base-utils/logger"
	"go-base-utils/saga/retry"

	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Run(
	ctx context.Context,
	cfg *HealthcheckConfig,
	readDb *sqlx.DB,
	writeDb *sqlx.DB,
	redis redis.UniversalClient,
	client *kafka.Conn,
	kafkaConnectFunc func() (*kafka.Conn, error),
) func() {
	return func() {
		itv := time.Duration(cfg.Interval) * time.Second
		health := healthcheck.NewHandler()
		readyCheck(ctx, cfg.GoroutineThreshold, health)
		liveCheck(ctx, health, itv, readDb, writeDb, redis, client, kafkaConnectFunc)
		gin.SetMode(gin.ReleaseMode)
		router := gin.New()
		router.Use(middlewares.Logging())
		router.Use(gin.WrapH(health))
		logger.Info("Heath check server listening on port", zap.String("Port", cfg.Port))
		if err := router.Run(cfg.Port); err != nil {
			logger.Warn("Heath check server", zap.Error(err))
		}
	}
}

func readyCheck(ctx context.Context, goRoutinesThreshold int, health healthcheck.Handler) {
	health.AddReadinessCheck(constants.GoroutineThreshold, healthcheck.GoroutineCountCheck(goRoutinesThreshold))
}

func liveCheck(ctx context.Context, health healthcheck.Handler, interval time.Duration, readDb *sqlx.DB, writeDb *sqlx.DB, redis redis.UniversalClient, client *kafka.Conn, connectFunc func() (*kafka.Conn, error)) {
	if readDb != nil {
		health.AddLivenessCheck(constants.ReadDatabase, healthcheck.AsyncWithContext(ctx, func() (err error) {
			err = readDb.DB.PingContext(ctx)
			if err != nil {
				logger.Error("Read database", zap.Error(err))
			}
			return
		}, interval))
	}
	if writeDb != nil {
		health.AddLivenessCheck(constants.WriteDatabase, healthcheck.AsyncWithContext(ctx, func() (err error) {
			err = writeDb.DB.PingContext(ctx)
			if err != nil {
				logger.Error("Readiness check write database", zap.Error(err))
			}
			return
		}, interval))
	}

	if redis != nil {
		health.AddLivenessCheck(constants.Redis, healthcheck.AsyncWithContext(ctx, func() error {
			err := redis.Ping(ctx).Err()
			if err != nil {
				logger.Error("Redis Readiness Check Fail", zap.Error(err))
			}
			return err
		}, interval))
	}

	if client != nil {
		health.AddLivenessCheck(constants.Kafka, healthcheck.AsyncWithContext(ctx, func() error {
			_, err := client.Brokers()
			if err != nil {
				err := retry.NewBackoff(
					retry.WithBackoffInitialInterval(3*time.Second),
					retry.WithBackoffMaxRetries(3)).Retry(ctx, func() error {
					client, err = connectFunc()
					if err != nil {
						logger.Error("Retry connecting to kafka ...")
					}
					return err
				})
				logger.Error("Kafka Readiness Check Fail", zap.Error(err))
			}
			return err
		}, interval))
	}
}
