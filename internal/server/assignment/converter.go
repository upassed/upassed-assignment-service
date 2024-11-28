package assignment

import (
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
	"github.com/upassed/upassed-assignment-service/pkg/client"
)

func ConvertToFindByFormIDResponse(assignment *business.FormAssignment) *client.AssignmentFindByFormIDResponse {
	groupIDs := make([]string, 0, len(assignment.GroupIDs))
	for _, groupID := range assignment.GroupIDs {
		groupIDs = append(groupIDs, groupID.String())
	}

	return &client.AssignmentFindByFormIDResponse{
		GroupIds: groupIDs,
	}
}

func ConvertToFindByGroupIDResponse(assignment *business.GroupAssignment) *client.AssignmentFindByGroupIDResponse {
	formIDs := make([]string, 0, len(assignment.FormIDs))
	for _, formID := range assignment.FormIDs {
		formIDs = append(formIDs, formID.String())
	}

	return &client.AssignmentFindByGroupIDResponse{
		FormIds: formIDs,
	}
}
