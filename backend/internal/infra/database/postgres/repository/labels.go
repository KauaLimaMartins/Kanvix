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

type Labels struct {
	DB *gorm.DB
}

func (r Labels) ListByWorkspaceIDs(ctx context.Context, workspaceIDs []string) ([]entity.Label, error) {
	var rows []model.Label
	if len(workspaceIDs) == 0 {
		return nil, nil
	}
	if err := r.DB.WithContext(ctx).Find(&rows, "workspace_id IN ?", workspaceIDs).Error; err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	out := make([]entity.Label, 0, len(rows))
	for _, l := range rows {
		out = append(out, labelToEntity(l))
	}
	return out, nil
}

func (r Labels) ListByWorkspace(ctx context.Context, workspaceID string) ([]entity.Label, error) {
	var rows []model.Label
	if err := r.DB.WithContext(ctx).Order("name asc").Find(&rows, "workspace_id = ?", workspaceID).Error; err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	out := make([]entity.Label, 0, len(rows))
	for _, l := range rows {
		out = append(out, labelToEntity(l))
	}
	return out, nil
}

func (r Labels) GetByID(ctx context.Context, id string) (entity.Label, error) {
	var l model.Label
	err := r.DB.WithContext(ctx).First(&l, "id = ?", id).Error
	if err == nil {
		return labelToEntity(l), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Label{}, domain.ErrNotFound
	}
	return entity.Label{}, fmt.Errorf("get label: %w", err)
}

func (r Labels) Create(ctx context.Context, l entity.Label) (entity.Label, error) {
	row := labelFromEntity(l)
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return entity.Label{}, fmt.Errorf("create label: %w", err)
	}
	return labelToEntity(row), nil
}

func (r Labels) Update(ctx context.Context, id string, patch map[string]any) (entity.Label, error) {
	res := r.DB.WithContext(ctx).Model(&model.Label{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return entity.Label{}, fmt.Errorf("update label: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.Label{}, domain.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r Labels) Delete(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&model.Label{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete label: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
