package services

import (
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uuid.UUID) (*models.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepository: userRepo}
}

func (s *userService) RegisterUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.ID = uuid.New()
	return s.userRepository.CreateUser(user)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *userService) GetUserById(id uuid.UUID) (*models.User, error) {
	return s.userRepository.GetUserById(id)
}
