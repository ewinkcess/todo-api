package service

import (
	"errors"
	"todo-api/internal/domain"
	"todo-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetProfile(userID uint) (*domain.User, error)
	UpdateProfile(userID uint, name, email string) (*domain.User, error)
	UpdatePassword(userID uint, oldPassword, newPassword string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uint) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}

func (s *userService) UpdateProfile(userID uint, name, email string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(email)
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("email sudah digunakan")
		}
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("email sudah digunakan")
		}
	}
	user.Name = name
	user.Email = email
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, errors.New("gagal mengupdate profil")
	}

	return user, nil
}

func (s *userService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("password lama tidak sesuai")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}

	user.Password = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return errors.New("gagal mengupdate password")
	}

	return nil
}
