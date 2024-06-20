package services

import (
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repository"
	"github.com/google/uuid"
)

type ProjectService interface {
	CreateProject(project *models.Project) (*models.Project, error)
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
