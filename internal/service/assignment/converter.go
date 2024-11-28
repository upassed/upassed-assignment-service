package assignment

import (
	"github.com/google/uuid"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
)

func ConvertToDomainAssignments(assignment *business.FormAssignment) []*domain.Assignment {
	domainAssignments := make([]*domain.Assignment, 0, len(assignment.GroupIDs))
	for _, groupID := range assignment.GroupIDs {
		domainAssignments = append(domainAssignments, &domain.Assignment{
			ID:      uuid.New(),
			FormID:  assignment.FormID,
			GroupID: groupID,
		})
	}

	return domainAssignments
}

func ConvertToAssignmentCreateResponse(domainAssignments []*domain.Assignment) *business.AssignmentCreateResponse {
	createdAssignmentIDs := make([]uuid.UUID, 0, len(domainAssignments))
	for _, domainAssignment := range domainAssignments {
		createdAssignmentIDs = append(createdAssignmentIDs, domainAssignment.ID)
	}

	return &business.AssignmentCreateResponse{
		CreatedAssignmentIDs: createdAssignmentIDs,
	}
}

func ConvertToBusinessFormAssignment(domainAssignments []*domain.Assignment) *business.FormAssignment {
	groupIDs := make([]uuid.UUID, 0, len(domainAssignments))
	for _, domainAssignment := range domainAssignments {
		groupIDs = append(groupIDs, domainAssignment.GroupID)
	}

	return &business.FormAssignment{
		GroupIDs: groupIDs,
	}
}

func ConvertToBusinessGroupAssignment(domainAssignments []*domain.Assignment) *business.GroupAssignment {
	formIDs := make([]uuid.UUID, 0, len(domainAssignments))
	for _, domainAssignment := range domainAssignments {
		formIDs = append(formIDs, domainAssignment.FormID)
	}

	return &business.GroupAssignment{
		FormIDs: formIDs,
	}
}
