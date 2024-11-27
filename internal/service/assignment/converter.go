package assignment

import (
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
)

func ConvertToDomainAssignments(assignment *business.Assignment) []*domain.Assignment {
	domainAssignments := make([]*domain.Assignment, 0, len(assignment.GroupIDs))
	for _, groupID := range assignment.GroupIDs {
		domainAssignments = append(domainAssignments, &domain.Assignment{
			ID:      assignment.ID,
			FormID:  assignment.FormID,
			GroupID: groupID,
		})
	}

	return domainAssignments
}
