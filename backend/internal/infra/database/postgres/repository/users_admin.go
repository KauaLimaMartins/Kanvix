package repository

import (
	"context"
	"fmt"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/infra/database/postgres/model"
)

func (s Store) Update(ctx context.Context, id string, patch map[string]any) (entity.User, error) {
	res := s.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return entity.User{}, fmt.Errorf("update user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.User{}, domain.ErrNotFound
	}
	return s.GetByID(ctx, id)
}

func (s Store) ListWorkspaceUsers(ctx context.Context, workspaceID string) ([]entity.WorkspaceUser, error) {
	type row struct {
		ID          string
		Email       string
		Name        string
		AvatarColor string
		Role        string
		Disabled    bool
	}
	var rows []row
	if err := s.DB.WithContext(ctx).Raw(`
		SELECT u.id, u.email, u.name, u.avatar_color, wm.role, u.disabled
		FROM users u
		JOIN workspace_members wm ON wm.user_id = u.id
		WHERE wm.workspace_id = ?
		ORDER BY u.name ASC
	`, workspaceID).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	out := make([]entity.WorkspaceUser, 0, len(rows))
	for _, r := range rows {
		out = append(out, entity.WorkspaceUser{
			ID:          r.ID,
			Email:       r.Email,
			Name:        r.Name,
			AvatarColor: r.AvatarColor,
			Role:        r.Role,
			Disabled:    r.Disabled,
		})
	}
	return out, nil
}
