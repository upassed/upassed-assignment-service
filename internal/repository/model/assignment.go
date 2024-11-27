package domain

import "github.com/google/uuid"

type Assignment struct {
	ID      uuid.UUID
	FormID  uuid.UUID
	GroupID uuid.UUID
}

func (Assignment) TableName() string {
	return "assignment"
}
