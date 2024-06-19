package repository

import (
	"database/sql"
	"github.com/drTragger/MykroTask/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uuid.UUID) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Password)
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserById(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
