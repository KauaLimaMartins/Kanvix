package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (r Repo) ListWorkspacesByOwner(ctx context.Context, ownerID string) ([]models.Workspace, error) {
	var ws []models.Workspace
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&ws, "owner_id = ?", ownerID).Error; err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	return ws, nil
}

func (r Repo) GetWorkspaceByID(ctx context.Context, id string) (models.Workspace, error) {
	var w models.Workspace
	err := r.DB.WithContext(ctx).First(&w, "id = ?", id).Error
	if err == nil {
		return w, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Workspace{}, repositories.ErrNotFound
	}
	return models.Workspace{}, fmt.Errorf("get workspace: %w", err)
}

func (r Repo) CreateWorkspace(ctx context.Context, w models.Workspace) (models.Workspace, error) {
	if err := r.DB.WithContext(ctx).Create(&w).Error; err != nil {
		return models.Workspace{}, fmt.Errorf("create workspace: %w", err)
	}
	return w, nil
}

func (r Repo) UpdateWorkspace(ctx context.Context, id string, patch map[string]any) (models.Workspace, error) {
	res := r.DB.WithContext(ctx).Model(&models.Workspace{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return models.Workspace{}, fmt.Errorf("update workspace: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return models.Workspace{}, repositories.ErrNotFound
	}
	return r.GetWorkspaceByID(ctx, id)
}

func (r Repo) DeleteWorkspace(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&models.Workspace{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete workspace: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return repositories.ErrNotFound
	}
	return nil
}
