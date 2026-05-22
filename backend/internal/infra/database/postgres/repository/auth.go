package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/infra/database/postgres/model"
)

func (s Store) Count(ctx context.Context) (int64, error) {
	var n int64
	if err := s.DB.WithContext(ctx).Model(&model.User{}).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return n, nil
}

func (s Store) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	var u model.User
	err := s.DB.WithContext(ctx).First(&u, "email = ?", email).Error
	if err == nil {
		return userToEntity(u), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.User{}, domain.ErrNotFound
	}
	return entity.User{}, fmt.Errorf("get user: %w", err)
}

func (s Store) GetByID(ctx context.Context, id string) (entity.User, error) {
	var u model.User
	err := s.DB.WithContext(ctx).First(&u, "id = ?", id).Error
	if err == nil {
		return userToEntity(u), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.User{}, domain.ErrNotFound
	}
	return entity.User{}, fmt.Errorf("get user: %w", err)
}

func (s Store) Create(ctx context.Context, u entity.User) (entity.User, error) {
	m := userFromEntity(u)
	if err := s.DB.WithContext(ctx).Create(&m).Error; err != nil {
		return entity.User{}, fmt.Errorf("create user: %w", err)
	}
	return userToEntity(m), nil
}
