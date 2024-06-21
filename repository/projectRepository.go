package repository

import (
	"database/sql"
	"github.com/drTragger/MykroTask/models"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	CreateProjectTx(tx *sql.Tx, project *models.Project) (*models.Project, error)
	GetProjectsForUser(userId uuid.UUID, page uint, perPage uint) ([]*models.Project, error)
	GetProjectById(projectId, memberId uuid.UUID) (*models.Project, error)
	UpdateProject(project *models.Project) (*models.Project, error)
	DeleteProject(projectId uuid.UUID) error
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProjectTx(tx *sql.Tx, project *models.Project) (*models.Project, error) {
	query := `INSERT INTO projects (id, name, description, start_date, end_date, owner_id) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	row := tx.QueryRow(query, project.ID, project.Name, project.Description, project.StartDate, project.EndDate, project.OwnerId)

	var p models.Project
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.StartDate, &p.EndDate, &p.OwnerId, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) GetProjectsForUser(userId uuid.UUID, page uint, perPage uint) ([]*models.Project, error) {
	var projects []*models.Project
	query := `SELECT p.id, p.name, p.description, p.start_date, p.end_date, p.owner_id, p.created_at, p.updated_at FROM projects AS p LEFT JOIN project_members AS pm ON p.id = pm.project_id WHERE pm.user_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userId, perPage, page*perPage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.StartDate, &p.EndDate, &p.OwnerId, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil
}

func (r *projectRepository) GetProjectById(projectId, memberId uuid.UUID) (*models.Project, error) {
	query := `SELECT p.id, p.name, p.description, p.start_date, p.end_date, p.owner_id, p.created_at, p.updated_at FROM projects AS p JOIN project_members AS pm ON p.id = pm.project_id WHERE p.id = $1 AND pm.user_id = $2`
	row := r.db.QueryRow(query, projectId, memberId)

	var p models.Project
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.StartDate, &p.EndDate, &p.OwnerId, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) UpdateProject(project *models.Project) (*models.Project, error) {
	query := `UPDATE projects SET name = $2, description = $3, end_date = $4, updated_at = $5 
              WHERE id = $1 RETURNING id, name, description, start_date, end_date, owner_id, created_at, updated_at`
	row := r.db.QueryRow(query, project.ID, project.Name, project.Description, project.EndDate, project.UpdatedAt)

	var p models.Project
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.StartDate, &p.EndDate, &p.OwnerId, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) DeleteProject(projectId uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(query, projectId)
	if err != nil {
		return err
	}
	return nil
}
