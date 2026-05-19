package gormrepo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"kanvix/backend/internal/models"
)

func (r Repo) ListTaskLabelsByProject(ctx context.Context, projectID string) ([]models.TaskLabel, error) {
	var joins []models.TaskLabel
	if err := r.DB.WithContext(ctx).
		Joins("JOIN tasks ON tasks.id = task_labels.task_id").
		Where("tasks.project_id = ?", projectID).
		Find(&joins).Error; err != nil {
		return nil, fmt.Errorf("list task labels: %w", err)
	}
	return joins, nil
}

func (r Repo) ReplaceTaskLabels(ctx context.Context, taskID string, labelIDs []string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("task_id = ?", taskID).Delete(&models.TaskLabel{}).Error; err != nil {
			return err
		}
		if len(labelIDs) == 0 {
			return nil
		}
		now := time.Now().UTC()
		joins := make([]models.TaskLabel, 0, len(labelIDs))
		for _, lid := range labelIDs {
			if lid == "" {
				continue
			}
			joins = append(joins, models.TaskLabel{TaskID: taskID, LabelID: lid, CreatedAt: now})
		}
		if len(joins) == 0 {
			return nil
		}
		return tx.Create(&joins).Error
	})
}
