package services

import (
	"errors"
	"github.com/drTragger/MykroTask/models"
	"github.com/drTragger/MykroTask/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	RegisterUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uuid.UUID) (*models.User, error)
	LoginUser(email, password string) (*models.JwtToken, error)
}

type userService struct {
	userRepository repository.UserRepository
	jwtKey         []byte
}

func NewUserService(userRepo repository.UserRepository, jwtKey []byte) UserService {
	return &userService{userRepository: userRepo, jwtKey: jwtKey}
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

func (s *userService) LoginUser(email, password string) (*models.JwtToken, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	expiresAt := time.Now().Add(time.Hour * 72)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID.String(),
		"exp":    expiresAt.Unix(),
	})

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	return &models.JwtToken{Token: tokenString, ExpiresAt: &expiresAt}, nil
}
