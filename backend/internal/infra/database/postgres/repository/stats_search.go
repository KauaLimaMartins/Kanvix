package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/domain/entity"
)

type StatsSearch struct {
	DB *gorm.DB
}

func (r StatsSearch) CountProjects(ctx context.Context, workspaceID string) (int64, error) {
	var projectCount int64
	if err := r.DB.WithContext(ctx).Table("projects").Where("workspace_id = ?", workspaceID).Count(&projectCount).Error; err != nil {
		return 0, fmt.Errorf("count projects: %w", err)
	}
	return projectCount, nil
}

func (r StatsSearch) CountTasksByProject(ctx context.Context, workspaceID string) (map[string]int, int, error) {
	type row struct {
		ProjectID string `gorm:"column:project_id"`
		Count     int    `gorm:"column:count"`
	}
	var rows []row
	if err := r.DB.WithContext(ctx).Raw(`
		SELECT tasks.project_id AS project_id, COUNT(*) AS count
		FROM tasks
		JOIN projects ON projects.id = tasks.project_id
		WHERE projects.workspace_id = ?
		GROUP BY tasks.project_id
	`, workspaceID).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}

	byProject := map[string]int{}
	total := 0
	for _, r := range rows {
		byProject[r.ProjectID] = r.Count
		total += r.Count
	}
	return byProject, total, nil
}

func (r StatsSearch) SearchProjects(ctx context.Context, workspaceID string, like string, limit int) ([]entity.ProjectSearchHit, error) {
	type row struct {
		ID   string `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var rows []row
	if err := r.DB.WithContext(ctx).Raw(`
		SELECT id, name FROM projects
		WHERE workspace_id = ? AND name LIKE ?
		ORDER BY created_at DESC
		LIMIT ?
	`, workspaceID, like, limit).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("search projects: %w", err)
	}
	out := make([]entity.ProjectSearchHit, 0, len(rows))
	for _, r := range rows {
		out = append(out, entity.ProjectSearchHit{ID: r.ID, Name: r.Name})
	}
	return out, nil
}

func (r StatsSearch) SearchTasks(ctx context.Context, workspaceID string, like string, limit int) ([]entity.TaskSearchHit, error) {
	type row struct {
		ID        string `gorm:"column:id"`
		Title     string `gorm:"column:title"`
		ProjectID string `gorm:"column:project_id"`
	}
	var rows []row
	if err := r.DB.WithContext(ctx).Raw(`
		SELECT tasks.id, tasks.title, tasks.project_id
		FROM tasks
		JOIN projects ON projects.id = tasks.project_id
		WHERE projects.workspace_id = ? AND tasks.title LIKE ?
		ORDER BY tasks.created_at DESC
		LIMIT ?
	`, workspaceID, like, limit).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("search tasks: %w", err)
	}
	out := make([]entity.TaskSearchHit, 0, len(rows))
	for _, r := range rows {
		out = append(out, entity.TaskSearchHit{ID: r.ID, Title: r.Title, ProjectID: r.ProjectID})
	}
	return out, nil
}

