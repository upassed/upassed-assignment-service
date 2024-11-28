package assignment_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-assignment-service/internal/server/assignment"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"testing"
)

func TestConvertToFindByFormIDResponse(t *testing.T) {
	businessAssignment := util.RandomBusinessFormAssignment()
	convertedResponse := assignment.ConvertToFindByFormIDResponse(businessAssignment)

	assert.Equal(t, len(businessAssignment.GroupIDs), len(convertedResponse.GetGroupIds()))
	for idx, groupID := range businessAssignment.GroupIDs {
		assert.Equal(t, groupID.String(), convertedResponse.GetGroupIds()[idx])
	}
}

func TestConvertToFindByGroupIDResponse(t *testing.T) {
	businessAssignment := util.RandomBusinessGroupAssignment()
	convertedResponse := assignment.ConvertToFindByGroupIDResponse(businessAssignment)

	assert.Equal(t, len(businessAssignment.FormIDs), len(convertedResponse.GetFormIds()))
	for idx, formID := range businessAssignment.FormIDs {
		assert.Equal(t, formID.String(), convertedResponse.GetFormIds()[idx])
	}
}
