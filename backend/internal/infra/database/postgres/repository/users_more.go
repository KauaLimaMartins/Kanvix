package repository

import (
	"context"
	"fmt"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/infra/database/postgres/model"
)

func (s Store) ListDistinctByWorkspaceIDs(ctx context.Context, workspaceIDs []string) ([]entity.User, error) {
	var rows []model.User
	if len(workspaceIDs) == 0 {
		return nil, nil
	}
	if err := s.DB.WithContext(ctx).
		Distinct("users.*").
		Joins("JOIN workspace_members wm ON wm.user_id = users.id").
		Where("wm.workspace_id IN ?", workspaceIDs).
		Order("users.name asc").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	out := make([]entity.User, 0, len(rows))
	for _, u := range rows {
		out = append(out, userToEntity(u))
	}
	return out, nil
}

