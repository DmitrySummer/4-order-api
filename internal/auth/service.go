package auth

import (
	"4-order-api/internal/user"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.Repository
}

func NewAuthService(userRepository *user.Repository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Login(phone, password string) (uint, error) {
	existedUser, _ := service.UserRepository.FindByPhone(phone)
	if existedUser == nil {
		return 0, errors.New(ErrWrongCreatetials)
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return 0, errors.New(ErrWrongCreatetials)
	}
	return existedUser.ID, nil
}

func (service *AuthService) Register(phone, password, name string) (uint, error) {
	existedUser, _ := service.UserRepository.FindByPhone(phone)
	if existedUser != nil {
		return 0, errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &user.User{
		Phone:    phone,
		Password: string(hashedPassword),
		Name:     name,
	}
	err = service.UserRepository.Create(user)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
