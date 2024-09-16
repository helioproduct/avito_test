package usecase

import (
	"avito_api/internal/entities/user"
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("пользователь не существует или некорректен")
)

type userUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(repo UserRepository) *userUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (uc *userUseCase) CreateUser(user *user.User) (string, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := uc.userRepo.CreateUser(user)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (uc *userUseCase) GetUserByID(id string) (*user.User, error) {
	return uc.userRepo.GetUserByID(id)
}

func (uc *userUseCase) GetUserByUsername(username string) (*user.User, error) {
	return uc.userRepo.GetUserByUsername(username)
}

func (uc *userUseCase) GetAllUsers() ([]*user.User, error) {
	return uc.userRepo.GetAllUsers()
}
