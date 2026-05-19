package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (r Repo) ListProjectsByWorkspace(ctx context.Context, workspaceID string) ([]models.Project, error) {
	var ps []models.Project
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&ps, "workspace_id = ?", workspaceID).Error; err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return ps, nil
}

func (r Repo) GetProjectByID(ctx context.Context, id string) (models.Project, error) {
	var p models.Project
	err := r.DB.WithContext(ctx).First(&p, "id = ?", id).Error
	if err == nil {
		return p, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Project{}, repositories.ErrNotFound
	}
	return models.Project{}, fmt.Errorf("get project: %w", err)
}

func (r Repo) CreateProject(ctx context.Context, p models.Project) (models.Project, error) {
	if err := r.DB.WithContext(ctx).Create(&p).Error; err != nil {
		return models.Project{}, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

func (r Repo) UpdateProject(ctx context.Context, id string, patch map[string]any) (models.Project, error) {
	res := r.DB.WithContext(ctx).Model(&models.Project{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return models.Project{}, fmt.Errorf("update project: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return models.Project{}, repositories.ErrNotFound
	}
	return r.GetProjectByID(ctx, id)
}

func (r Repo) DeleteProject(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&models.Project{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete project: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return repositories.ErrNotFound
	}
	return nil
}
