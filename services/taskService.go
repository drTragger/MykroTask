package services

import (
	"database/sql"
	"errors"
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repository"
	"github.com/google/uuid"
)

type TaskService interface {
	CreateTask(task *models.Task) (*models.Task, bool, error)
	GetTasksForUser(projectId, userId uuid.UUID) ([]*models.Task, error)
}

type taskService struct {
	taskRepository          repository.TaskRepository
	projectMemberRepository repository.ProjectMemberRepository
}

func NewTaskService(taskRepository repository.TaskRepository, projectMemberRepository repository.ProjectMemberRepository) TaskService {
	return &taskService{taskRepository: taskRepository, projectMemberRepository: projectMemberRepository}
}

func (s *taskService) CreateTask(task *models.Task) (*models.Task, bool, error) {
	_, err := s.projectMemberRepository.GetMember(task.ProjectID, task.CreatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, true, nil
		}
		return nil, false, err
	}

	task.ID = uuid.New()
	task, err = s.taskRepository.CreateTask(task)
	return task, false, err
}

func (s *taskService) GetTasksForUser(projectId, userId uuid.UUID) ([]*models.Task, error) {
	return s.taskRepository.GetTasksForUser(projectId, userId)
}
