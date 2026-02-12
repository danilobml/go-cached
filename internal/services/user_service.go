package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/danilobml/go-cached/internal/models"
	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	FindById(ctx context.Context, id string) (*models.User, error)
	List(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, name, email string) error
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type UserService struct {
	userRepo UserRepository
	cache    Cache
}

func NewUserService(userRepo UserRepository, cache Cache) *UserService {
	return &UserService{
		userRepo: userRepo,
		cache:    cache,
	}
}

func (us *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	key := fmt.Sprintf("q:UserByID:%s", id)
	cachedUser, err := us.cache.Get(ctx, key)
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
		_ = us.cache.Del(ctx, key)
	} else if err != redis.Nil {
		log.Println("cache GET failed", err)
	}

	user, err := us.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	userToCache, err := json.Marshal(user)
	if err == nil {
		err = us.cache.Set(ctx, key, userToCache, 2*time.Minute)
		if err != nil {
			log.Println("Failed to create cache")
		}
	}

	return user, nil
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	key := "q:Users:All"
	cached, err := us.cache.Get(ctx, key)
	if err == nil {
		var users []*models.User
		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			return users, nil
		}
		_ = us.cache.Del(ctx, key)
	} else if err != redis.Nil {
		log.Println("cache GET failed", err)
	}

	users, err := us.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	if toCache, err := json.Marshal(users); err == nil {
		if err := us.cache.Set(ctx, key, toCache, 30*time.Second); err != nil {
			log.Println("cache SET failed", err)
		}
	}

	return users, nil
}

func (us *UserService) CreateUser(ctx context.Context, name, email string) error {
	err := us.userRepo.Create(ctx, name, email)
	if err != nil {
		return err
	}

	_ = us.cache.Del(ctx, "q:Users:All")

	return nil
}
