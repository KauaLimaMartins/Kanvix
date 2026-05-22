package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories/gormrepo"
)

type Bootstrap struct {
	User       dto.MeUser      `json:"user"`
	Workspaces []dto.Workspace `json:"workspaces"`
	Projects   []dto.Project   `json:"projects"`
	Columns    []dto.Column    `json:"columns"`
	Tasks      []dto.Task      `json:"tasks"`
	Labels     []dto.Label     `json:"labels"`
	Users      []dto.User      `json:"users"`
}

type AppService struct {
	Repo     gormrepo.Repo
	Redis    *redis.Client
	CacheTTL time.Duration
}

func (s AppService) Bootstrap(ctx context.Context, owner models.User) (Bootstrap, error) {
	ownerID := owner.ID
	var memberships []models.WorkspaceMember
	if err := s.Repo.DB.WithContext(ctx).Find(&memberships, "user_id = ?", ownerID).Error; err != nil {
		return Bootstrap{}, fmt.Errorf("list memberships: %w", err)
	}
	roleByWS := map[string]string{}
	wsIDs := make([]string, 0, len(memberships))
	for _, m := range memberships {
		wsIDs = append(wsIDs, m.WorkspaceID)
		roleByWS[m.WorkspaceID] = m.Role
	}

	var workspaces []models.Workspace
	if len(wsIDs) > 0 {
		if err := s.Repo.DB.WithContext(ctx).Order("created_at asc").Find(&workspaces, "id IN ?", wsIDs).Error; err != nil {
			return Bootstrap{}, fmt.Errorf("list workspaces: %w", err)
		}
	}

	var projects []models.Project
	var labels []models.Label
	if len(wsIDs) > 0 {
		if err := s.Repo.DB.WithContext(ctx).Find(&projects, "workspace_id IN ?", wsIDs).Error; err != nil {
			return Bootstrap{}, fmt.Errorf("list projects: %w", err)
		}
		if err := s.Repo.DB.WithContext(ctx).Find(&labels, "workspace_id IN ?", wsIDs).Error; err != nil {
			return Bootstrap{}, fmt.Errorf("list labels: %w", err)
		}
	}

	pIDs := make([]string, 0, len(projects))
	for _, p := range projects {
		pIDs = append(pIDs, p.ID)
	}

	var columns []models.Column
	var tasks []models.Task
	if len(pIDs) > 0 {
		if err := s.Repo.DB.WithContext(ctx).Find(&columns, "project_id IN ?", pIDs).Error; err != nil {
			return Bootstrap{}, fmt.Errorf("list columns: %w", err)
		}
		if err := s.Repo.DB.WithContext(ctx).Find(&tasks, "project_id IN ?", pIDs).Error; err != nil {
			return Bootstrap{}, fmt.Errorf("list tasks: %w", err)
		}
	}

	taskIDs := make([]string, 0, len(tasks))
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}

	var joins []models.TaskLabel
	if len(taskIDs) > 0 {
		if err := s.Repo.DB.WithContext(ctx).Find(&joins, "task_id IN ?", taskIDs).Error; err != nil {
			return Bootstrap{}, fmt.Errorf("list task labels: %w", err)
		}
	}
	taskToLabels := map[string][]string{}
	for _, j := range joins {
		taskToLabels[j.TaskID] = append(taskToLabels[j.TaskID], j.LabelID)
	}

	var users []models.User
	if len(wsIDs) > 0 {
		if err := s.Repo.DB.WithContext(ctx).
			Distinct("users.*").
			Joins("JOIN workspace_members wm ON wm.user_id = users.id").
			Where("wm.workspace_id IN ?", wsIDs).
			Order("users.name asc").
			Find(&users).Error; err != nil {
			return Bootstrap{}, fmt.Errorf("list users: %w", err)
		}
	}

	out := Bootstrap{
		User:       dto.MeUser{ID: owner.ID, Email: owner.Email, Name: owner.Name, AvatarColor: owner.AvatarColor, Role: owner.Role},
		Workspaces: make([]dto.Workspace, 0, len(workspaces)),
		Projects:   make([]dto.Project, 0, len(projects)),
		Columns:    make([]dto.Column, 0, len(columns)),
		Tasks:      make([]dto.Task, 0, len(tasks)),
		Labels:     make([]dto.Label, 0, len(labels)),
		Users:      make([]dto.User, 0, len(users)),
	}

	for _, u := range users {
		out.Users = append(out.Users, dto.User{ID: u.ID, Name: u.Name, AvatarColor: u.AvatarColor, Disabled: u.Disabled})
	}

	for _, w := range workspaces {
		out.Workspaces = append(out.Workspaces, dto.Workspace{ID: w.ID, Name: w.Name, Icon: w.Icon, Color: w.Color, Role: roleByWS[w.ID]})
	}
	for _, p := range projects {
		out.Projects = append(out.Projects, dto.Project{ID: p.ID, WorkspaceID: p.WorkspaceID, Name: p.Name, Description: p.Description})
	}
	for _, c := range columns {
		out.Columns = append(out.Columns, dto.Column{ID: c.ID, ProjectID: c.ProjectID, Name: c.Name, Order: c.Order})
	}
	for _, l := range labels {
		out.Labels = append(out.Labels, dto.Label{ID: l.ID, WorkspaceID: l.WorkspaceID, Name: l.Name, Color: l.Color})
	}
	for _, t := range tasks {
		lbls := taskToLabels[t.ID]
		if lbls == nil {
			lbls = []string{}
		}
		out.Tasks = append(out.Tasks, dto.Task{
			ID:          t.ID,
			ProjectID:   t.ProjectID,
			ColumnID:    t.ColumnID,
			Title:       t.Title,
			Description: t.Description,
			Labels:      lbls,
			DueDate:     t.DueDate,
			AssigneeID:  t.AssigneeID,
			Order:       t.Order,
			CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}

	return out, nil
}
