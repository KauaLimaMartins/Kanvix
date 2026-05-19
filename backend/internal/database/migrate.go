package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
)

func AutoMigrate(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.Project{},
		&models.Column{},
		&models.Task{},
		&models.Label{},
		&models.TaskLabel{},
	); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}
	return nil
}
