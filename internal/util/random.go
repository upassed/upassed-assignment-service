package util

import (
	"github.com/google/uuid"
	event "github.com/upassed/upassed-assignment-service/internal/messanging/model"
)

func RandomEventAssignmentCreateRequest() *event.AssignmentCreateRequest {
	return &event.AssignmentCreateRequest{
		FormID:  uuid.NewString(),
		GroupID: uuid.NewString(),
	}
}
