package business

import "github.com/google/uuid"

type Assignment struct {
	ID      uuid.UUID
	FormID  uuid.UUID
	GroupID uuid.UUID
}

type AssignmentCreateResponse struct {
	CreatedAssignmentID uuid.UUID
}
