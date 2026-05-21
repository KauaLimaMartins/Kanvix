package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (r Repo) ListSubtasksByTask(ctx context.Context, taskID string) ([]models.Subtask, error) {
	var ss []models.Subtask
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&ss, "task_id = ?", taskID).Error; err != nil {
		return nil, fmt.Errorf("list subtasks: %w", err)
	}
	return ss, nil
}

func (r Repo) GetSubtaskByID(ctx context.Context, id string) (models.Subtask, error) {
	var s models.Subtask
	err := r.DB.WithContext(ctx).First(&s, "id = ?", id).Error
	if err == nil {
		return s, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Subtask{}, repositories.ErrNotFound
	}
	return models.Subtask{}, fmt.Errorf("get subtask: %w", err)
}

func (r Repo) CreateSubtask(ctx context.Context, s models.Subtask) (models.Subtask, error) {
	if err := r.DB.WithContext(ctx).Create(&s).Error; err != nil {
		return models.Subtask{}, fmt.Errorf("create subtask: %w", err)
	}
	return s, nil
}

func (r Repo) UpdateSubtaskDone(ctx context.Context, id string, done bool, updatedAt time.Time) (models.Subtask, error) {
	res := r.DB.WithContext(ctx).Model(&models.Subtask{}).Where("id = ?", id).Updates(map[string]any{
		"done":       done,
		"updated_at": updatedAt,
	})
	if res.Error != nil {
		return models.Subtask{}, fmt.Errorf("update subtask: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return models.Subtask{}, repositories.ErrNotFound
	}
	return r.GetSubtaskByID(ctx, id)
}
