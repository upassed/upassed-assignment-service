package event

import "github.com/go-playground/validator/v10"

type AssignmentCreateRequest struct {
	FormID  string `json:"form_id,omitempty" validate:"required,uuid"`
	GroupID string `json:"group_id,omitempty" validate:"required,uuid"`
}

func (request *AssignmentCreateRequest) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("uuid", ValidateUUID())

	if err := validate.Struct(*request); err != nil {
		return err
	}

	return nil
}
