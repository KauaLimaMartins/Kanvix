package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (r Repo) GetWorkspaceMember(ctx context.Context, workspaceID, userID string) (models.WorkspaceMember, error) {
	var m models.WorkspaceMember
	err := r.DB.WithContext(ctx).First(&m, "workspace_id = ? AND user_id = ?", workspaceID, userID).Error
	if err == nil {
		return m, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.WorkspaceMember{}, repositories.ErrNotFound
	}
	return models.WorkspaceMember{}, fmt.Errorf("get workspace member: %w", err)
}

func (r Repo) UpsertWorkspaceMember(ctx context.Context, m models.WorkspaceMember) error {
	if err := r.DB.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("upsert workspace member: %w", err)
	}
	return nil
}

func (r Repo) ListWorkspaceMembers(ctx context.Context, workspaceID string) ([]models.WorkspaceMember, error) {
	var ms []models.WorkspaceMember
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&ms, "workspace_id = ?", workspaceID).Error; err != nil {
		return nil, fmt.Errorf("list workspace members: %w", err)
	}
	return ms, nil
}

func (r Repo) ListWorkspacesForUser(ctx context.Context, userID string) ([]models.Workspace, error) {
	var ws []models.Workspace
	if err := r.DB.WithContext(ctx).
		Joins("JOIN workspace_members wm ON wm.workspace_id = workspaces.id").
		Where("wm.user_id = ?", userID).
		Order("workspaces.created_at asc").
		Find(&ws).Error; err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	return ws, nil
}

