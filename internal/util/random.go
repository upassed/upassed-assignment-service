package util

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	event "github.com/upassed/upassed-assignment-service/internal/messanging/model"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
)

func RandomEventAssignmentCreateRequest() *event.AssignmentCreateRequest {
	groupIDsCount := gofakeit.IntRange(10, 20)
	groupIDs := make([]string, 0, groupIDsCount)

	for i := 0; i < groupIDsCount; i++ {
		groupIDs = append(groupIDs, uuid.NewString())
	}

	return &event.AssignmentCreateRequest{
		FormID:   uuid.NewString(),
		GroupIDs: groupIDs,
	}
}

func RandomBusinessAssignment() *business.Assignment {
	groupIDsCount := gofakeit.IntRange(10, 20)
	groupIDs := make([]uuid.UUID, 0, groupIDsCount)

	for i := 0; i < groupIDsCount; i++ {
		groupIDs = append(groupIDs, uuid.New())
	}

	return &business.Assignment{
		ID:       uuid.New(),
		FormID:   uuid.New(),
		GroupIDs: groupIDs,
	}
}
