package services

import (
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repository"
	"github.com/google/uuid"
)

const ProjectsPerPage uint = 10

type ProjectService interface {
	CreateProject(project *models.Project) (*models.Project, error)
	GetProjectsForUser(userId uuid.UUID, page uint) ([]*models.Project, error)
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
