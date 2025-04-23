package auth

import (
	"4-order-api/internal/user"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.Repository
}

func NewAuthService(userRepository *user.Repository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Login(phone, password string) (string, error) {
	existedUser, err := service.UserRepository.FindByPhone(phone)
	if err != nil {
		return "", fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}
	if existedUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	return existedUser.Phone, nil
}

func (service *AuthService) Register(phone, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByPhone(phone)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &user.User{
		Phone:    phone,
		Password: string(hashedPassword),
		Name:     name,
	}
	err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Phone, nil
}
