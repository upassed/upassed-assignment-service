package assignment

import (
	"context"
	"github.com/google/uuid"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
)

func (repository *repositoryImpl) FindByGroupID(ctx context.Context, groupID uuid.UUID) ([]*domain.Assignment, error) {
	//TODO implement me
	panic("implement me")
}
