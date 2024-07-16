package services

import (
	"database/sql"
	"errors"
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repository"
	"github.com/google/uuid"
	"time"
)

type TaskService interface {
	CreateTask(task *models.Task) (*models.Task, bool, error)
	GetTasksForUser(projectId, memberId, userId uuid.UUID) ([]*models.Task, bool, error)
	GetTaskById(projectId, taskId, memberId uuid.UUID) (*models.Task, bool, error)
	DeleteTask(projectId, taskId, memberId uuid.UUID) (bool, error)
	GetTasksForProject(projectId, userId uuid.UUID) ([]*models.Task, bool, error)
	UpdateTask(projectId, taskId, userId uuid.UUID, task *models.Task) (*models.Task, bool, error)
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

	_, err = s.projectMemberRepository.GetMember(task.ProjectID, task.Assignee)
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

func (s *taskService) GetTasksForUser(projectId, memberId, userId uuid.UUID) ([]*models.Task, bool, error) {
	_, err := s.projectMemberRepository.GetMember(projectId, memberId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, true, nil
		}
		return nil, false, err
	}

	tasks, err := s.taskRepository.GetTasksForUser(projectId, userId)

	return tasks, false, err
}

func (s *taskService) GetTaskById(projectId, taskId, memberId uuid.UUID) (*models.Task, bool, error) {
	_, err := s.projectMemberRepository.GetMember(projectId, memberId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, true, nil
		}
		return nil, false, err
	}

	task, err := s.taskRepository.GetTaskById(projectId, taskId)
	if err != nil {
		return nil, false, err
	}

	return task, false, nil
}

func (s *taskService) DeleteTask(projectId, taskId, memberId uuid.UUID) (bool, error) {
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

	err = s.taskRepository.DeleteTask(projectId, taskId)
	if err != nil {
		return false, err
	}
	return false, nil
}

func (s *taskService) GetTasksForProject(projectId, userId uuid.UUID) ([]*models.Task, bool, error) {
	_, err := s.projectMemberRepository.GetMember(projectId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, true, nil
		}
		return nil, false, err
	}

	tasks, err := s.taskRepository.GetTasksForProject(projectId, userId)
	if err != nil {
		return nil, false, err
	}
	return tasks, false, nil
}

func (s *taskService) UpdateTask(projectId, taskId, userId uuid.UUID, task *models.Task) (*models.Task, bool, error) {
	_, err := s.projectMemberRepository.GetMember(projectId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, true, nil
		}
		return nil, false, err
	}

	task.ID = taskId
	task.ProjectID = projectId
	task.UpdatedAt = time.Now()
	task, err = s.taskRepository.UpdateTask(task)
	return task, false, err
}
