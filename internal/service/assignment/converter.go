package assignment

import (
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
)

func ConvertToDomainAssignment(assignment *business.Assignment) *domain.Assignment {
	return &domain.Assignment{
		ID:       assignment.ID,
		FormID:   assignment.FormID,
		GroupIDs: assignment.GroupIDs,
	}
}
