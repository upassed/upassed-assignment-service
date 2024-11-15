package assignment_test

import (
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
	request.GroupID = "invalid_uuid"

	err := request.Validate()
	require.Error(t, err)
}

func TestAssignmentCreateRequestValidation_Valid(t *testing.T) {
	request := util.RandomEventAssignmentCreateRequest()

	err := request.Validate()
	require.NoError(t, err)
}
