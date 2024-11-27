package assignment

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-assignment-service/internal/caching/assignment"
	"github.com/upassed/upassed-assignment-service/internal/config"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	"gorm.io/gorm"
	"log/slog"
)

type Repository interface {
	Save(ctx context.Context, assignments []*domain.Assignment) error
	CheckDuplicates(ctx context.Context, assignments []*domain.Assignment) ([]*domain.Assignment, error)
}

type repositoryImpl struct {
	db    *gorm.DB
	cache *assignment.RedisClient
	cfg   *config.Config
	log   *slog.Logger
}

func New(db *gorm.DB, cacheClient *redis.Client, cfg *config.Config, log *slog.Logger) Repository {
	cache := assignment.New(cacheClient, cfg, log)
	return &repositoryImpl{
		db:    db,
		cache: cache,
		cfg:   cfg,
		log:   log,
	}
}
