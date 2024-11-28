package assignment_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-assignment-service/internal/service/assignment"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"testing"
)

func TestConvertToDomainAssignment(t *testing.T) {
	businessAssignment := util.RandomBusinessFormAssignment()
	domainAssignments := assignment.ConvertToDomainAssignments(businessAssignment)

	assert.Equal(t, len(businessAssignment.GroupIDs), len(domainAssignments))
	for idx, domainAssignment := range domainAssignments {
		assert.NotNil(t, domainAssignment.ID)
		assert.Equal(t, businessAssignment.FormID, domainAssignment.FormID)
		assert.Equal(t, businessAssignment.GroupIDs[idx], domainAssignment.GroupID)
	}
}

func TestConvertToAssignmentCreateResponse(t *testing.T) {
	domainAssignments := util.RandomDomainAssignments()
	createResponse := assignment.ConvertToAssignmentCreateResponse(domainAssignments)

	for idx, domainAssignment := range domainAssignments {
		assert.Equal(t, domainAssignment.ID, createResponse.CreatedAssignmentIDs[idx])
	}
}

func TestConvertToBusinessFormAssignment(t *testing.T) {
	domainAssignments := util.RandomDomainAssignments()
	businessFormAssignment := assignment.ConvertToBusinessFormAssignment(domainAssignments)

	for idx, domainAssignment := range domainAssignments {
		assert.Equal(t, domainAssignment.GroupID, businessFormAssignment.GroupIDs[idx])
	}
}

func TestConvertToBusinessGroupAssignment(t *testing.T) {
	domainAssignments := util.RandomDomainAssignments()
	businessGroupAssignment := assignment.ConvertToBusinessGroupAssignment(domainAssignments)

	for idx, domainAssignment := range domainAssignments {
		assert.Equal(t, domainAssignment.FormID, businessGroupAssignment.FormIDs[idx])
	}
}
