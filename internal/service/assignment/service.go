package assignment

import (
	"context"
	"github.com/google/uuid"
	"github.com/upassed/upassed-assignment-service/internal/config"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
	"log/slog"
)

type Service interface {
	Create(ctx context.Context, assignment *business.FormAssignment) (*business.AssignmentCreateResponse, error)
	FindByFormID(ctx context.Context, formID uuid.UUID) (*business.FormAssignment, error)
	FindByGroupID(ctx context.Context, groupID uuid.UUID) (*business.GroupAssignment, error)
}

type serviceImpl struct {
	cfg        *config.Config
	log        *slog.Logger
	repository repository
}

type repository interface {
	Save(ctx context.Context, assignment []*domain.Assignment) error
	CheckDuplicates(ctx context.Context, assignments []*domain.Assignment) ([]*domain.Assignment, error)
	FindByFormID(ctx context.Context, formID uuid.UUID) ([]*domain.Assignment, error)
	FindByGroupID(ctx context.Context, groupID uuid.UUID) ([]*domain.Assignment, error)
}

func New(cfg *config.Config, log *slog.Logger, repository repository) Service {
	return &serviceImpl{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}
