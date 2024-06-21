package repository

import (
	"database/sql"
	"github.com/drTragger/MykroTask/models"
	"github.com/google/uuid"
)

type ProjectMemberRepository interface {
	CreateMember(member *models.ProjectMember) (*models.ProjectMember, error)
	CreateMemberTx(tx *sql.Tx, member *models.ProjectMember) (*models.ProjectMember, error)
	GetMember(projectId, userId uuid.UUID) (*models.ProjectMember, error)
	GetMembers(projectId uuid.UUID) ([]*models.ProjectMember, error)
	DeleteMember(projectId, userId uuid.UUID) error
}

type projectMemberRepository struct {
	db *sql.DB
}

func NewProjectMemberRepository(db *sql.DB) ProjectMemberRepository {
	return &projectMemberRepository{db: db}
}

func (r *projectMemberRepository) CreateMember(member *models.ProjectMember) (*models.ProjectMember, error) {
	query := `INSERT INTO project_members (project_id, user_id, role) VALUES ($1, $2, $3) RETURNING *;`
	row := r.db.QueryRow(query, member.ProjectId, member.UserId, member.Role.String())

	err := row.Scan(&member.ProjectId, &member.UserId, &member.Role, &member.JoinedAt)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *projectMemberRepository) CreateMemberTx(tx *sql.Tx, member *models.ProjectMember) (*models.ProjectMember, error) {
	query := `INSERT INTO project_members (project_id, user_id, role) VALUES ($1, $2, $3) RETURNING *;`
	row := tx.QueryRow(query, member.ProjectId, member.UserId, member.Role.String())

	err := row.Scan(&member.ProjectId, &member.UserId, &member.Role, &member.JoinedAt)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *projectMemberRepository) GetMember(projectId, userId uuid.UUID) (*models.ProjectMember, error) {
	var member models.ProjectMember
	query := `SELECT u.email, u.name, pm.role, pm.joined_at FROM project_members AS pm JOIN users AS u ON pm.user_id = u.id WHERE project_id = $1 AND user_id = $2;`

	row := r.db.QueryRow(query, projectId, userId)
	err := row.Scan(&member.Email, &member.Name, &member.Role, &member.JoinedAt)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *projectMemberRepository) GetMembers(projectId uuid.UUID) ([]*models.ProjectMember, error) {
	query := `SELECT u.id, u.email, u.name, pm.role, pm.joined_at FROM project_members AS pm JOIN users AS u ON pm.user_id = u.id WHERE project_id = $1;`
	rows, err := r.db.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.ProjectMember
	for rows.Next() {
		var m models.ProjectMember
		err := rows.Scan(&m.UserId, &m.Email, &m.Name, &m.Role, &m.JoinedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, nil
}

func (r *projectMemberRepository) DeleteMember(projectId, userId uuid.UUID) error {
	query := `DELETE FROM project_members WHERE project_id = $1 AND user_id = $2;`
	_, err := r.db.Exec(query, projectId, userId)
	if err != nil {
		return err
	}
	return nil
}
