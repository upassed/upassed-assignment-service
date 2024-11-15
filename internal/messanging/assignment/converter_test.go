package assignment_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-assignment-service/internal/messanging/assignment"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"testing"
)

func TestConvertToAssignmentCreateRequest_InvalidBytes(t *testing.T) {
	invalidBytes := make([]byte, 10)
	_, err := assignment.ConvertToAssignmentCreateRequest(invalidBytes)
	require.Error(t, err)
}

func TestConvertToAssignmentCreateRequest_ValidBytes(t *testing.T) {
	initialRequest := util.RandomEventAssignmentCreateRequest()
	initialRequestBytes, err := json.Marshal(initialRequest)
	require.NoError(t, err)

	convertedRequest, err := assignment.ConvertToAssignmentCreateRequest(initialRequestBytes)
	require.NoError(t, err)

	assert.Equal(t, initialRequest.FormID, convertedRequest.FormID)
	assert.Equal(t, initialRequest.GroupID, convertedRequest.GroupID)
}

func TestConvertToBusinessAssignment(t *testing.T) {
	request := util.RandomEventAssignmentCreateRequest()
	businessAssignment := assignment.ConvertToBusinessAssignment(request)

	require.NotNil(t, businessAssignment.ID)

	assert.Equal(t, request.FormID, businessAssignment.FormID.String())
	assert.Equal(t, request.GroupID, businessAssignment.GroupID.String())
}