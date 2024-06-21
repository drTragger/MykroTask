package services

import (
	"database/sql"
	"errors"
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repository"
	"github.com/google/uuid"
	"time"
)

const ProjectsPerPage uint = 10

type ProjectService interface {
	CreateProject(project *models.Project, ownerId uuid.UUID) (*models.Project, error)
	GetProjectsForUser(memberId uuid.UUID, page uint) ([]*models.Project, error)
	GetProjectById(projectId, memberId uuid.UUID) (*models.Project, error)
	UpdateProject(project *models.Project, memberId uuid.UUID) (*models.Project, bool, error)
	DeleteProject(projectId, memberId uuid.UUID) (bool, error)
}

type projectService struct {
	projectRepository       repository.ProjectRepository
	projectMemberRepository repository.ProjectMemberRepository
	db                      *sql.DB
}

func NewProjectService(projectRepo repository.ProjectRepository, projectMemberRepo repository.ProjectMemberRepository, db *sql.DB) ProjectService {
	return &projectService{projectRepository: projectRepo, projectMemberRepository: projectMemberRepo, db: db}
}

func (s *projectService) CreateProject(project *models.Project, ownerId uuid.UUID) (*models.Project, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	project.ID = uuid.New()
	project, err = s.projectRepository.CreateProjectTx(tx, project)
	if err != nil {
		return nil, err
	}

	member := models.ProjectMember{
		ProjectId: project.ID,
		UserId:    ownerId,
		Role:      models.RoleOwner,
		JoinedAt:  time.Now(),
	}
	_, err = s.projectMemberRepository.CreateMemberTx(tx, &member)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *projectService) GetProjectsForUser(memberId uuid.UUID, page uint) ([]*models.Project, error) {
	return s.projectRepository.GetProjectsForUser(memberId, page, ProjectsPerPage)
}

func (s *projectService) GetProjectById(projectId, memberId uuid.UUID) (*models.Project, error) {
	return s.projectRepository.GetProjectById(projectId, memberId)
}

func (s *projectService) UpdateProject(project *models.Project, memberId uuid.UUID) (*models.Project, bool, error) {
	member, err := s.projectMemberRepository.GetMember(project.ID, memberId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, true, nil
		}
		return nil, false, err
	}

	if !member.CanEditProject() {
		return nil, true, nil
	}

	project.UpdatedAt = time.Now()
	p, err := s.projectRepository.UpdateProject(project)
	return p, false, err
}

func (s *projectService) DeleteProject(projectId, memberId uuid.UUID) (bool, error) {
	member, err := s.projectMemberRepository.GetMember(projectId, memberId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}

	if !member.CanEditProject() {
		return true, nil
	}

	return false, s.projectRepository.DeleteProject(projectId)
}
