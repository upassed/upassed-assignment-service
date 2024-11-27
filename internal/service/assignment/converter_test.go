package assignment_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-assignment-service/internal/service/assignment"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"testing"
)

func TestConvertToDomainAssignment(t *testing.T) {
	businessAssignment := util.RandomBusinessAssignment()
	domainAssignments := assignment.ConvertToDomainAssignments(businessAssignment)

	assert.Equal(t, len(businessAssignment.GroupIDs), len(domainAssignments))
	for idx, domainAssignment := range domainAssignments {
		assert.Equal(t, businessAssignment.ID, domainAssignment.ID)
		assert.Equal(t, businessAssignment.FormID, domainAssignment.FormID)
		assert.Equal(t, businessAssignment.GroupIDs[idx], domainAssignment.GroupID)
	}
}
