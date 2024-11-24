package event

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	errEmptyGroupsList   = errors.New("number of groups in assignment should be > 0")
	errDuplicateGroupIDs = errors.New("all group ids should be unique")
	errInvalidGroupID    = errors.New("all group ids should be valid uuids")
)

type AssignmentCreateRequest struct {
	FormID   string   `json:"form_id,omitempty" validate:"required,uuid"`
	GroupIDs []string `json:"group_ids" validate:"required"`
}

func (request *AssignmentCreateRequest) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("uuid", ValidateUUID())

	if err := validate.Struct(*request); err != nil {
		return err
	}

	if len(request.GroupIDs) == 0 {
		return errEmptyGroupsList
	}

	visitedGroupIds := make(map[string]struct{})
	for _, groupID := range request.GroupIDs {
		if _, ok := visitedGroupIds[groupID]; ok {
			return errDuplicateGroupIDs
		}

		visitedGroupIds[groupID] = struct{}{}

		if _, err := uuid.Parse(groupID); err != nil {
			return errInvalidGroupID
		}
	}

	return nil
}
