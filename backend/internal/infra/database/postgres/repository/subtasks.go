package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/infra/database/postgres/model"
)

type Subtasks struct {
	DB *gorm.DB
}

func (r Subtasks) ListByTask(ctx context.Context, taskID string) ([]entity.Subtask, error) {
	var rows []model.Subtask
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&rows, "task_id = ?", taskID).Error; err != nil {
		return nil, fmt.Errorf("list subtasks: %w", err)
	}
	out := make([]entity.Subtask, 0, len(rows))
	for _, s := range rows {
		out = append(out, subtaskToEntity(s))
	}
	return out, nil
}

func (r Subtasks) GetByID(ctx context.Context, id string) (entity.Subtask, error) {
	var s model.Subtask
	err := r.DB.WithContext(ctx).First(&s, "id = ?", id).Error
	if err == nil {
		return subtaskToEntity(s), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Subtask{}, domain.ErrNotFound
	}
	return entity.Subtask{}, fmt.Errorf("get subtask: %w", err)
}

func (r Subtasks) Create(ctx context.Context, s entity.Subtask) (entity.Subtask, error) {
	row := subtaskFromEntity(s)
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return entity.Subtask{}, fmt.Errorf("create subtask: %w", err)
	}
	return subtaskToEntity(row), nil
}

func (r Subtasks) Update(ctx context.Context, id string, patch map[string]any) (entity.Subtask, error) {
	res := r.DB.WithContext(ctx).Model(&model.Subtask{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return entity.Subtask{}, fmt.Errorf("update subtask: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.Subtask{}, domain.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r Subtasks) Delete(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&model.Subtask{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete subtask: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r Subtasks) UpdateFields(ctx context.Context, id string, title *string, done *bool, updatedAt time.Time) (entity.Subtask, error) {
	patch := map[string]any{"updated_at": updatedAt}
	if title != nil {
		patch["title"] = *title
	}
	if done != nil {
		patch["done"] = *done
	}
	return r.Update(ctx, id, patch)
}
