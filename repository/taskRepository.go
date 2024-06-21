package repository

import (
	"database/sql"
	"github.com/drTragger/MykroTask/models"
	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateTask(task *models.Task) (*models.Task, error)
	GetTasksForUser(projectId, userId uuid.UUID) ([]*models.Task, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task *models.Task) (*models.Task, error) {
	query := `INSERT INTO tasks(id, title, description, status, priority, assignee, due_date, project_id, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;`
	row := r.db.QueryRow(query, task.ID, task.Title, task.Description, task.Status, task.Priority, task.Assignee, task.DueDate, task.ProjectID, task.CreatedBy)

	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.Assignee, &task.DueDate, &task.ProjectID, &task.CreatedBy, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *taskRepository) GetTasksForUser(projectId, userId uuid.UUID) ([]*models.Task, error) {
	var tasks []*models.Task
	query := `SELECT * FROM tasks WHERE project_id = $1 AND assignee = $2;`
	rows, err := r.db.Query(query, projectId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Assignee, &t.DueDate, &t.ProjectID, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}
