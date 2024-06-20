package services

import (
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repository"
	"github.com/google/uuid"
	"time"
)

const ProjectsPerPage uint = 10

type ProjectService interface {
	CreateProject(project *models.Project) (*models.Project, error)
	GetProjectsForUser(userId uuid.UUID, page uint) ([]*models.Project, error)
	GetProjectById(projectId uuid.UUID) (*models.Project, error)
	UpdateProject(project *models.Project) (*models.Project, error)
	CheckUserPermission(projectId, userId uuid.UUID) (bool, error)
	DeleteProject(projectId uuid.UUID) error
}

type projectService struct {
	projectRepository repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepository: projectRepo}
}

func (s *projectService) CreateProject(project *models.Project) (*models.Project, error) {
	project.ID = uuid.New()
	return s.projectRepository.CreateProject(project)
}

func (s *projectService) GetProjectsForUser(userId uuid.UUID, page uint) ([]*models.Project, error) {
	return s.projectRepository.GetProjectsForUser(userId, page, ProjectsPerPage)
}

func (s *projectService) GetProjectById(projectId uuid.UUID) (*models.Project, error) {
	return s.projectRepository.GetProjectById(projectId)
}

func (s *projectService) UpdateProject(project *models.Project) (*models.Project, error) {
	project.UpdatedAt = time.Now()
	return s.projectRepository.UpdateProject(project)
}

func (s *projectService) CheckUserPermission(projectId, userId uuid.UUID) (bool, error) {
	project, err := s.projectRepository.GetProjectById(projectId)
	if err != nil {
		return false, err
	}

	return project.OwnerId == userId, nil
}

func (s *projectService) DeleteProject(projectId uuid.UUID) error {
	return s.projectRepository.DeleteProject(projectId)
}
