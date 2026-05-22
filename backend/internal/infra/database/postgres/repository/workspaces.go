package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/domain/entity"
)

type WorkspaceRepo struct {
	DB *gorm.DB
}

func (r WorkspaceRepo) Create(ctx context.Context, w entity.Workspace) (entity.Workspace, error) {
	m := workspaceFromEntity(w)
	if err := r.DB.WithContext(ctx).Create(&m).Error; err != nil {
		return entity.Workspace{}, fmt.Errorf("create workspace: %w", err)
	}
	return entity.Workspace{
		ID:        m.ID,
		OwnerID:   m.OwnerID,
		Name:      m.Name,
		Icon:      m.Icon,
		Color:     m.Color,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

type WorkspaceMemberRepo struct {
	DB *gorm.DB
}

func (r WorkspaceMemberRepo) Create(ctx context.Context, m entity.WorkspaceMember) error {
	row := workspaceMemberFromEntity(m)
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("create membership: %w", err)
	}
	return nil
}

