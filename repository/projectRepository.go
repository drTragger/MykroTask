package repository

import (
	"database/sql"
	"github.com/drTragger/MykroTask/models"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) (*models.Project, error)
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(project *models.Project) (*models.Project, error) {
	query := `INSERT INTO projects (id, name, description, start_date, end_date, owner_id) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	row := r.db.QueryRow(query, project.ID, project.Name, project.Description, project.StartDate, project.EndDate, project.OwnerId)

	var p models.Project
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.StartDate, &p.EndDate, &p.OwnerId, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}