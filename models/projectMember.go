package models

import (
	"github.com/google/uuid"
	"time"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func (r Role) String() string {
	return string(r)
}

func GetValidRoles() []string {
	return []string{
		RoleOwner.String(),
		RoleAdmin.String(),
		RoleMember.String(),
	}
}

type ProjectMember struct {
	ProjectId uuid.UUID `json:"projectId,omitempty" validate:"required"`
	UserId    uuid.UUID `json:"userId" validate:"required"`
	Role      Role      `json:"role" validate:"required,role"`
	JoinedAt  time.Time `json:"joinedAt"`
	User
}

func (pm *ProjectMember) CanEditProject() bool {
	return pm.Role == RoleAdmin || pm.Role == RoleOwner
}
