package business

import "github.com/google/uuid"

type FormAssignment struct {
	FormID   uuid.UUID
	GroupIDs []uuid.UUID
}

type GroupAssignment struct {
	GroupID uuid.UUID
	FormIDs []uuid.UUID
}

type AssignmentCreateResponse struct {
	CreatedAssignmentIDs []uuid.UUID
}
