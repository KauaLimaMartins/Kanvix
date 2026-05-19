package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/repositories"
)

func (s AppService) WorkspaceStats(ctx context.Context, ownerID, workspaceID string) (dto.WorkspaceStats, error) {
	ws, err := s.Repo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return dto.WorkspaceStats{}, err
	}
	if ws.OwnerID != ownerID {
		return dto.WorkspaceStats{}, repositories.ErrForbidden
	}

	key := "stats:" + workspaceID
	if s.Redis != nil {
		if raw, err := s.Redis.Get(ctx, key).Bytes(); err == nil {
			var cached dto.WorkspaceStats
			if json.Unmarshal(raw, &cached) == nil {
				return cached, nil
			}
		}
	}

	var projectCount int64
	if err := s.Repo.DB.WithContext(ctx).Table("projects").Where("workspace_id = ?", workspaceID).Count(&projectCount).Error; err != nil {
		return dto.WorkspaceStats{}, fmt.Errorf("count projects: %w", err)
	}

	type row struct {
		ProjectID string
		Count     int
	}
	var rows []row
	if err := s.Repo.DB.WithContext(ctx).Raw(`
		SELECT tasks.project_id AS project_id, COUNT(*) AS count
		FROM tasks
		JOIN projects ON projects.id = tasks.project_id
		WHERE projects.workspace_id = ?
		GROUP BY tasks.project_id
	`, workspaceID).Scan(&rows).Error; err != nil {
		return dto.WorkspaceStats{}, fmt.Errorf("count tasks: %w", err)
	}

	byProject := map[string]int{}
	total := 0
	for _, r := range rows {
		byProject[r.ProjectID] = r.Count
		total += r.Count
	}

	out := dto.WorkspaceStats{
		WorkspaceID:    workspaceID,
		ProjectCount:   int(projectCount),
		TaskCount:      total,
		TasksByProject: byProject,
	}

	if s.Redis != nil {
		if b, err := json.Marshal(out); err == nil {
			_ = s.Redis.Set(ctx, key, b, s.CacheTTL).Err()
		}
	}

	return out, nil
}

func (s AppService) Search(ctx context.Context, ownerID, workspaceID, q string, limit int) (dto.SearchResponse, error) {
	ws, err := s.Repo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return dto.SearchResponse{}, err
	}
	if ws.OwnerID != ownerID {
		return dto.SearchResponse{}, repositories.ErrForbidden
	}

	q = strings.TrimSpace(q)
	if q == "" {
		return dto.SearchResponse{Query: "", Results: []dto.SearchResult{}}, nil
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	cacheKey := ""
	if s.Redis != nil {
		cacheKey = "search:" + workspaceID + ":" + strings.ToLower(q) + ":" + fmt.Sprintf("%d", limit)
		if raw, err := s.Redis.Get(ctx, cacheKey).Bytes(); err == nil {
			var cached dto.SearchResponse
			if json.Unmarshal(raw, &cached) == nil {
				return cached, nil
			}
		}
	}

	like := "%" + q + "%"

	type projectRow struct {
		ID   string
		Name string
	}
	type taskRow struct {
		ID        string
		Title     string
		ProjectID string
	}

	var prs []projectRow
	if err := s.Repo.DB.WithContext(ctx).Raw(`
		SELECT id, name FROM projects
		WHERE workspace_id = ? AND name LIKE ?
		ORDER BY created_at DESC
		LIMIT ?
	`, workspaceID, like, limit).Scan(&prs).Error; err != nil {
		return dto.SearchResponse{}, fmt.Errorf("search projects: %w", err)
	}

	var trs []taskRow
	if err := s.Repo.DB.WithContext(ctx).Raw(`
		SELECT tasks.id, tasks.title, tasks.project_id
		FROM tasks
		JOIN projects ON projects.id = tasks.project_id
		WHERE projects.workspace_id = ? AND tasks.title LIKE ?
		ORDER BY tasks.created_at DESC
		LIMIT ?
	`, workspaceID, like, limit).Scan(&trs).Error; err != nil {
		return dto.SearchResponse{}, fmt.Errorf("search tasks: %w", err)
	}

	results := make([]dto.SearchResult, 0, len(prs)+len(trs))
	for _, p := range prs {
		results = append(results, dto.SearchResult{Type: "project", ID: p.ID, Name: p.Name, WorkspaceID: workspaceID})
	}
	for _, t := range trs {
		results = append(results, dto.SearchResult{Type: "task", ID: t.ID, Title: t.Title, ProjectID: t.ProjectID, WorkspaceID: workspaceID})
	}

	out := dto.SearchResponse{Query: q, Results: results}
	if s.Redis != nil && cacheKey != "" {
		if b, err := json.Marshal(out); err == nil {
			_ = s.Redis.Set(ctx, cacheKey, b, s.CacheTTL).Err()
		}
	}
	return out, nil
}
