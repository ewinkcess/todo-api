package service

import (
	"errors"
	"os"
	"time"
	"todo-api/internal/domain"
	"todo-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(name, email, password string) (*domain.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(name, email, password string) (*domain.User, error) {
	existingUSer, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUSer.ID != 0 {
		return nil, errors.New("email sudah terdaftar")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("gagal mendaftarkan user")
	}
	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", errors.New("gagal memuat token")
	}
	return tokenString, nil
}
