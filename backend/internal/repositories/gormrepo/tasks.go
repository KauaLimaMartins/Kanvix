package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (r Repo) ListTasksByProject(ctx context.Context, projectID string) ([]models.Task, error) {
	var ts []models.Task
	if err := r.DB.WithContext(ctx).Order("column_id asc, `order` asc").Find(&ts, "project_id = ?", projectID).Error; err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return ts, nil
}

func (r Repo) ListTasksByColumn(ctx context.Context, columnID string) ([]models.Task, error) {
	var ts []models.Task
	if err := r.DB.WithContext(ctx).Order("`order` asc").Find(&ts, "column_id = ?", columnID).Error; err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return ts, nil
}

func (r Repo) GetTaskByID(ctx context.Context, id string) (models.Task, error) {
	var t models.Task
	err := r.DB.WithContext(ctx).First(&t, "id = ?", id).Error
	if err == nil {
		return t, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Task{}, repositories.ErrNotFound
	}
	return models.Task{}, fmt.Errorf("get task: %w", err)
}

func (r Repo) CreateTask(ctx context.Context, t models.Task) (models.Task, error) {
	if err := r.DB.WithContext(ctx).Create(&t).Error; err != nil {
		return models.Task{}, fmt.Errorf("create task: %w", err)
	}
	return t, nil
}

func (r Repo) UpdateTask(ctx context.Context, id string, patch map[string]any) (models.Task, error) {
	res := r.DB.WithContext(ctx).Model(&models.Task{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return models.Task{}, fmt.Errorf("update task: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return models.Task{}, repositories.ErrNotFound
	}
	return r.GetTaskByID(ctx, id)
}

func (r Repo) DeleteTask(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&models.Task{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete task: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return repositories.ErrNotFound
	}
	return nil
}

func (r Repo) CountTasksInColumn(ctx context.Context, columnID string) (int64, error) {
	var n int64
	if err := r.DB.WithContext(ctx).Model(&models.Task{}).Where("column_id = ?", columnID).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count tasks: %w", err)
	}
	return n, nil
}
