package assignment_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"testing"
)

func TestAssignmentCreateRequestValidation_InvalidFormID(t *testing.T) {
	request := util.RandomEventAssignmentCreateRequest()
	request.FormID = "invalid_uuid"

	err := request.Validate()
	require.Error(t, err)
}

func TestAssignmentCreateRequestValidation_InvalidGroupID(t *testing.T) {
	request := util.RandomEventAssignmentCreateRequest()
	request.GroupIDs[0] = "invalid_uuid"

	err := request.Validate()
	require.Error(t, err)
}

func TestAssignmentCreateRequestValidation_DuplicateGroupIDs(t *testing.T) {
	request := util.RandomEventAssignmentCreateRequest()
	duplicateGroupID := uuid.NewString()
	request.GroupIDs[0] = duplicateGroupID
	request.GroupIDs[1] = duplicateGroupID

	err := request.Validate()
	require.Error(t, err)
}

func TestAssignmentCreateRequestValidation_EmptyGroupID(t *testing.T) {
	request := util.RandomEventAssignmentCreateRequest()
	request.GroupIDs = []string{}

	err := request.Validate()
	require.Error(t, err)
}

func TestAssignmentCreateRequestValidation_Valid(t *testing.T) {
	request := util.RandomEventAssignmentCreateRequest()

	err := request.Validate()
	require.NoError(t, err)
}
