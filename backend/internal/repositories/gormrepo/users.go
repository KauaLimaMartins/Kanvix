package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (r Repo) GetUserByID(ctx context.Context, id string) (models.User, error) {
	var u models.User
	err := r.DB.WithContext(ctx).First(&u, "id = ?", id).Error
	if err == nil {
		return u, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, repositories.ErrNotFound
	}
	return models.User{}, fmt.Errorf("get user: %w", err)
}

func (r Repo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := r.DB.WithContext(ctx).First(&u, "email = ?", email).Error
	if err == nil {
		return u, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, repositories.ErrNotFound
	}
	return models.User{}, fmt.Errorf("get user: %w", err)
}

func (r Repo) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.DB.WithContext(ctx).Order("name asc").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

func (r Repo) CreateUser(ctx context.Context, u models.User) (models.User, error) {
	if err := r.DB.WithContext(ctx).Create(&u).Error; err != nil {
		return models.User{}, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}
