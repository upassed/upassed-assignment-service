package business

import "github.com/google/uuid"

type FormAssignment struct {
	ID       uuid.UUID
	FormID   uuid.UUID
	GroupIDs []uuid.UUID
}

type GroupAssignment struct {
	ID      uuid.UUID
	GroupID uuid.UUID
	FormIDs []uuid.UUID
}

type AssignmentCreateResponse struct {
	CreatedAssignmentIDs []uuid.UUID
}
