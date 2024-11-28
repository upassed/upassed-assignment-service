package assignment

import (
	"context"
	"github.com/google/uuid"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
)

func (service *serviceImpl) FindByGroupID(ctx context.Context, groupID uuid.UUID) (*business.GroupAssignment, error) {
	//TODO implement me
	panic("implement me")
}
