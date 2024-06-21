package services

import (
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repository"
	"github.com/google/uuid"
)

type ProjectMemberService interface {
	CreateMember(member *models.ProjectMember, userId uuid.UUID) (*models.ProjectMember, bool, error)
	GetMember(projectId, userId uuid.UUID) (*models.ProjectMember, error)
	GetMembers(projectId uuid.UUID) ([]*models.ProjectMember, error)
	DeleteMember(projectId, memberId, userId uuid.UUID) (bool, error)
}

type projectMemberService struct {
	projectMemberRepository repository.ProjectMemberRepository
}

func NewProjectMemberService(projectMemberRepo repository.ProjectMemberRepository) ProjectMemberService {
	return &projectMemberService{projectMemberRepository: projectMemberRepo}
}

func (s *projectMemberService) CreateMember(member *models.ProjectMember, userId uuid.UUID) (*models.ProjectMember, bool, error) {
	m, err := s.GetMember(member.ProjectId, userId)
	if err != nil {
		return nil, false, err
	}

	if !m.CanEditProject() {
		return nil, true, nil
	}

	member, err = s.projectMemberRepository.CreateMember(member)

	return member, false, err
}

func (s *projectMemberService) GetMember(projectId, userId uuid.UUID) (*models.ProjectMember, error) {
	return s.projectMemberRepository.GetMember(projectId, userId)
}

func (s *projectMemberService) GetMembers(projectId uuid.UUID) ([]*models.ProjectMember, error) {
	return s.projectMemberRepository.GetMembers(projectId)
}

func (s *projectMemberService) DeleteMember(projectId, memberId, userId uuid.UUID) (bool, error) {
	member, err := s.GetMember(projectId, userId)
	if err != nil {
		return false, err
	}

	if !member.CanEditProject() || memberId == userId {
		return true, nil
	}

	return false, s.projectMemberRepository.DeleteMember(projectId, memberId)
}
