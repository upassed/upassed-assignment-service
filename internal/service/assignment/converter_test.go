package assignment_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-assignment-service/internal/service/assignment"
	"github.com/upassed/upassed-assignment-service/internal/util"
	"testing"
)

func TestConvertToDomainAssignment(t *testing.T) {
	businessAssignment := util.RandomBusinessAssignment()
	domainAssignment := assignment.ConvertToDomainAssignment(businessAssignment)

	assert.Equal(t, businessAssignment.ID, domainAssignment.ID)
	assert.Equal(t, businessAssignment.FormID, domainAssignment.FormID)
	assert.Equal(t, len(businessAssignment.GroupIDs), len(domainAssignment.GroupIDs))

	for idx, businessGroupID := range businessAssignment.GroupIDs {
		assert.Equal(t, businessGroupID, domainAssignment.GroupIDs[idx])
	}
}
