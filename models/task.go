package models

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description string    `json:"description"`
	Status      string    `json:"status" validate:"required"`
	Priority    string    `json:"priority" validate:"required"`
	Assignee    uuid.UUID `json:"assignee" validate:"required,uuid"`
	DueDate     time.Time `json:"dueDate"`
	ProjectID   uuid.UUID `json:"projectId"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
