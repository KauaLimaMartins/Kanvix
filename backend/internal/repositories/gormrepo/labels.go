package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (r Repo) ListLabelsByWorkspace(ctx context.Context, workspaceID string) ([]models.Label, error) {
	var ls []models.Label
	if err := r.DB.WithContext(ctx).Order("name asc").Find(&ls, "workspace_id = ?", workspaceID).Error; err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	return ls, nil
}

func (r Repo) GetLabelByID(ctx context.Context, id string) (models.Label, error) {
	var l models.Label
	err := r.DB.WithContext(ctx).First(&l, "id = ?", id).Error
	if err == nil {
		return l, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Label{}, repositories.ErrNotFound
	}
	return models.Label{}, fmt.Errorf("get label: %w", err)
}

func (r Repo) CreateLabel(ctx context.Context, l models.Label) (models.Label, error) {
	if err := r.DB.WithContext(ctx).Create(&l).Error; err != nil {
		return models.Label{}, fmt.Errorf("create label: %w", err)
	}
	return l, nil
}

func (r Repo) UpdateLabel(ctx context.Context, id string, patch map[string]any) (models.Label, error) {
	res := r.DB.WithContext(ctx).Model(&models.Label{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return models.Label{}, fmt.Errorf("update label: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return models.Label{}, repositories.ErrNotFound
	}
	return r.GetLabelByID(ctx, id)
}

func (r Repo) DeleteLabel(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&models.Label{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete label: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return repositories.ErrNotFound
	}
	return nil
}
