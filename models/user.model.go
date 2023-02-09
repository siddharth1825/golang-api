package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email     string    `gorm:"uniqueIndex;not null;primary_key"`
	Songs []Songs `gorm:"foreignKey:UserEmail; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type CreateUserSchema struct {
	Email     string `json:"email" validate:"required"`
}

type UpdateUserSchema struct {
	Email     string `json:"email,omitempty"`
}



