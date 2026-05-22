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

type WorkspaceMembers struct {
	DB *gorm.DB
}

func (r WorkspaceMembers) CountAdmins(ctx context.Context, workspaceID string) (int64, error) {
	var n int64
	if err := r.DB.WithContext(ctx).
		Model(&model.WorkspaceMember{}).
		Where("workspace_id = ? AND role = ?", workspaceID, "admin").
		Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count admins: %w", err)
	}
	return n, nil
}

func (r WorkspaceMembers) Get(ctx context.Context, workspaceID, userID string) (entity.WorkspaceMember, error) {
	var m model.WorkspaceMember
	err := r.DB.WithContext(ctx).First(&m, "workspace_id = ? AND user_id = ?", workspaceID, userID).Error
	if err == nil {
		return workspaceMemberToEntity(m), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkspaceMember{}, domain.ErrNotFound
	}
	return entity.WorkspaceMember{}, fmt.Errorf("get workspace member: %w", err)
}

func (r WorkspaceMembers) Upsert(ctx context.Context, m entity.WorkspaceMember) error {
	row := workspaceMemberFromEntity(m)
	if err := r.DB.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("upsert workspace member: %w", err)
	}
	return nil
}

func (r WorkspaceMembers) ListByWorkspace(ctx context.Context, workspaceID string) ([]entity.WorkspaceMember, error) {
	var rows []model.WorkspaceMember
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&rows, "workspace_id = ?", workspaceID).Error; err != nil {
		return nil, fmt.Errorf("list workspace members: %w", err)
	}
	out := make([]entity.WorkspaceMember, 0, len(rows))
	for _, m := range rows {
		out = append(out, workspaceMemberToEntity(m))
	}
	return out, nil
}

func (r WorkspaceMembers) ListWorkspacesForUser(ctx context.Context, userID string) ([]entity.Workspace, error) {
	var rows []model.Workspace
	if err := r.DB.WithContext(ctx).
		Joins("JOIN workspace_members wm ON wm.workspace_id = workspaces.id").
		Where("wm.user_id = ?", userID).
		Order("workspaces.created_at asc").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	out := make([]entity.Workspace, 0, len(rows))
	for _, w := range rows {
		out = append(out, workspaceToEntity(w))
	}
	return out, nil
}

func (r WorkspaceMembers) ListByUser(ctx context.Context, userID string) ([]entity.WorkspaceMember, error) {
	var rows []model.WorkspaceMember
	if err := r.DB.WithContext(ctx).Find(&rows, "user_id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("list memberships: %w", err)
	}
	out := make([]entity.WorkspaceMember, 0, len(rows))
	for _, m := range rows {
		out = append(out, workspaceMemberToEntity(m))
	}
	return out, nil
}

func (r WorkspaceMembers) Delete(ctx context.Context, workspaceID, userID string) error {
	res := r.DB.WithContext(ctx).Delete(&model.WorkspaceMember{}, "workspace_id = ? AND user_id = ?", workspaceID, userID)
	if res.Error != nil {
		return fmt.Errorf("delete workspace member: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
