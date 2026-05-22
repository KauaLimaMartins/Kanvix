package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/infra/database/postgres/model"
)

func AutoMigrate(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).AutoMigrate(
		&model.User{},
		&model.Workspace{},
		&model.WorkspaceMember{},
		&model.Project{},
		&model.Column{},
		&model.Task{},
		&model.Subtask{},
		&model.Label{},
		&model.TaskLabel{},
	); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}
	return nil
}

