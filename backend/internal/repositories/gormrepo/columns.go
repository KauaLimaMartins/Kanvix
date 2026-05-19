package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (r Repo) ListColumnsByProject(ctx context.Context, projectID string) ([]models.Column, error) {
	var cs []models.Column
	if err := r.DB.WithContext(ctx).Order("`order` asc").Find(&cs, "project_id = ?", projectID).Error; err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}
	return cs, nil
}

func (r Repo) GetColumnByID(ctx context.Context, id string) (models.Column, error) {
	var col models.Column
	err := r.DB.WithContext(ctx).First(&col, "id = ?", id).Error
	if err == nil {
		return col, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Column{}, repositories.ErrNotFound
	}
	return models.Column{}, fmt.Errorf("get column: %w", err)
}

func (r Repo) CreateColumn(ctx context.Context, col models.Column) (models.Column, error) {
	if err := r.DB.WithContext(ctx).Create(&col).Error; err != nil {
		return models.Column{}, fmt.Errorf("create column: %w", err)
	}
	return col, nil
}

func (r Repo) UpdateColumn(ctx context.Context, id string, patch map[string]any) (models.Column, error) {
	res := r.DB.WithContext(ctx).Model(&models.Column{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return models.Column{}, fmt.Errorf("update column: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return models.Column{}, repositories.ErrNotFound
	}
	return r.GetColumnByID(ctx, id)
}

func (r Repo) DeleteColumn(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&models.Column{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete column: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return repositories.ErrNotFound
	}
	return nil
}

func (r Repo) CountColumnsInProject(ctx context.Context, projectID string) (int64, error) {
	var n int64
	if err := r.DB.WithContext(ctx).Model(&models.Column{}).Where("project_id = ?", projectID).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count columns: %w", err)
	}
	return n, nil
}
