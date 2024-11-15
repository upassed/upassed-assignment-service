package assignment

import (
	"context"
	"github.com/upassed/upassed-assignment-service/internal/config"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
	"log/slog"
)

type Service interface {
	Create(ctx context.Context, assignment *business.Assignment) (*business.AssignmentCreateResponse, error)
}

type serviceImpl struct {
	cfg        *config.Config
	log        *slog.Logger
	repository repository
}

type repository interface {
	Save(ctx context.Context, assignment *domain.Assignment) error
}

func New(cfg *config.Config, log *slog.Logger, repository repository) Service {
	return &serviceImpl{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}
