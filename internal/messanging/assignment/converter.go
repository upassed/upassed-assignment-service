package assignment

import (
	"encoding/json"
	"github.com/google/uuid"
	event "github.com/upassed/upassed-assignment-service/internal/messanging/model"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
)

func ConvertToAssignmentCreateRequest(messageBody []byte) (*event.AssignmentCreateRequest, error) {
	var request event.AssignmentCreateRequest
	if err := json.Unmarshal(messageBody, &request); err != nil {
		return nil, err
	}

	return &request, nil
}

func ConvertToBusinessAssignment(request *event.AssignmentCreateRequest) *business.Assignment {
	groupIDs := make([]uuid.UUID, 0, len(request.GroupIDs))
	for _, groupID := range request.GroupIDs {
		groupIDs = append(groupIDs, uuid.MustParse(groupID))
	}

	return &business.Assignment{
		ID:       uuid.New(),
		FormID:   uuid.MustParse(request.FormID),
		GroupIDs: groupIDs,
	}
}
