package models

import(

	"time"
	"github.com/google/uuid"

)

type Songs struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Link string `gorm:"type:varchar(255);not null"`
	UserEmail string `gorm:"type:varchar(255);not null;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateSongSchema struct {
	Link string `json:"link" validate:"required"`
	UserEmail string `json:"useremail" validate:"required"`
}
