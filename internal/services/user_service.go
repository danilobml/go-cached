package services

import (
	"context"

	"github.com/danilobml/go-cached/internal/models"
)

type UserRepository interface {
	FindById(ctx context.Context, id string) (*models.User, error)
	List(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context,	name, email string) error
}

type UserService struct {
	userRepo UserRepository 
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (us *UserService) GetUser(ctx context.Context,	id string) (*models.User, error) {
	user, err := us.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	users, err := us.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserService) CreateUser(ctx context.Context, name, email string) error {
	err := us.userRepo.Create(ctx, name, email)
	if err != nil {
		return err
	}

	return nil
}
