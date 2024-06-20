package models

import (
	"github.com/google/uuid"
	"time"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required,min=3,max=255"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	OwnerId     uuid.UUID `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
