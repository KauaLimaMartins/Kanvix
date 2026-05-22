package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"kanvix/backend/internal/infra/database/postgres/model"
)

func EnsureOwnerMemberships(ctx context.Context, db *gorm.DB) error {
	var workspaces []model.Workspace
	if err := db.WithContext(ctx).Find(&workspaces).Error; err != nil {
		return fmt.Errorf("list workspaces: %w", err)
	}
	now := time.Now().UTC()
	for _, ws := range workspaces {
		m := model.WorkspaceMember{
			WorkspaceID: ws.ID,
			UserID:      ws.OwnerID,
			Role:        "admin",
			CreatedAt:   now,
		}
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&m).Error; err != nil {
			return fmt.Errorf("ensure owner membership: %w", err)
		}
	}
	return nil
}

