package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (s AppService) ListWorkspaceUsers(ctx context.Context, requesterID, workspaceID string) ([]dto.UserDetail, error) {
	if _, err := s.requireWorkspaceRole(ctx, requesterID, workspaceID, "admin"); err != nil {
		return nil, err
	}

	type row struct {
		ID          string
		Email       string
		Name        string
		AvatarColor string
		Role        string
		Disabled    bool
	}
	var rows []row
	if err := s.Repo.DB.WithContext(ctx).Raw(`
		SELECT u.id, u.email, u.name, u.avatar_color, wm.role, u.disabled
		FROM users u
		JOIN workspace_members wm ON wm.user_id = u.id
		WHERE wm.workspace_id = ?
		ORDER BY u.name ASC
	`, workspaceID).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	out := make([]dto.UserDetail, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.UserDetail{
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

func (s AppService) CreateWorkspaceUser(ctx context.Context, requesterID, workspaceID, email, password, name, role string) (dto.UserDetail, error) {
	if _, err := s.requireWorkspaceRole(ctx, requesterID, workspaceID, "admin"); err != nil {
		return dto.UserDetail{}, err
	}
	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)
	role = strings.TrimSpace(strings.ToLower(role))
	if email == "" || name == "" {
		return dto.UserDetail{}, fmt.Errorf("email and name required")
	}
	if role != "admin" && role != "member" {
		role = "member"
	}
	if len(password) < 8 {
		return dto.UserDetail{}, fmt.Errorf("password must be at least 8 characters")
	}

	_, err := s.Repo.GetUserByEmail(ctx, email)
	if err == nil {
		return dto.UserDetail{}, fmt.Errorf("email already exists")
	}
	if !errors.Is(err, repositories.ErrNotFound) {
		return dto.UserDetail{}, err
	}

	pwHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserDetail{}, fmt.Errorf("hash password: %w", err)
	}
	now := time.Now().UTC()
	u := models.User{
		ID:           uuid.NewString(),
		Email:        email,
		Name:         name,
		AvatarColor:  "#64748b",
		Role:         role,
		PasswordHash: string(pwHashBytes),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	u, err = s.Repo.CreateUser(ctx, u)
	if err != nil {
		return dto.UserDetail{}, err
	}
	if err := s.Repo.UpsertWorkspaceMember(ctx, models.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      u.ID,
		Role:        role,
		CreatedAt:   now,
	}); err != nil {
		return dto.UserDetail{}, err
	}
	return dto.UserDetail{ID: u.ID, Email: u.Email, Name: u.Name, AvatarColor: u.AvatarColor, Role: role}, nil
}

func (s AppService) UpdateWorkspaceUserRole(ctx context.Context, requesterID, workspaceID, userID, role string) error {
	if _, err := s.requireWorkspaceRole(ctx, requesterID, workspaceID, "admin"); err != nil {
		return err
	}
	role = strings.TrimSpace(strings.ToLower(role))
	if role != "admin" && role != "member" {
		return fmt.Errorf("invalid role")
	}
	if err := s.Repo.UpsertWorkspaceMember(ctx, models.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		CreatedAt:   time.Now().UTC(),
	}); err != nil {
		return err
	}
	_, err := s.Repo.UpdateUser(ctx, userID, map[string]any{"role": role, "updated_at": time.Now().UTC()})
	return err
}

func (s AppService) UpdateWorkspaceUser(ctx context.Context, requesterID, workspaceID, userID string, name *string, role *string) error {
	if _, err := s.requireWorkspaceRole(ctx, requesterID, workspaceID, "admin"); err != nil {
		return err
	}
	if _, err := s.Repo.GetWorkspaceMember(ctx, workspaceID, userID); err != nil {
		return err
	}
	if name == nil && role == nil {
		return fmt.Errorf("invalid request")
	}

	patch := map[string]any{
		"updated_at": time.Now().UTC(),
	}
	if name != nil {
		n := strings.TrimSpace(*name)
		if n == "" {
			return fmt.Errorf("invalid name")
		}
		patch["name"] = n
	}
	if role != nil {
		r := strings.TrimSpace(strings.ToLower(*role))
		if r != "admin" && r != "member" {
			return fmt.Errorf("invalid role")
		}
		if err := s.Repo.UpsertWorkspaceMember(ctx, models.WorkspaceMember{
			WorkspaceID: workspaceID,
			UserID:      userID,
			Role:        r,
			CreatedAt:   time.Now().UTC(),
		}); err != nil {
			return err
		}
		patch["role"] = r
	}

	_, err := s.Repo.UpdateUser(ctx, userID, patch)
	return err
}

type DeleteWorkspaceUserAction string

const (
	DeleteWorkspaceUserActionUnassign DeleteWorkspaceUserAction = "unassign"
	DeleteWorkspaceUserActionReassign DeleteWorkspaceUserAction = "reassign"
	DeleteWorkspaceUserActionDisable  DeleteWorkspaceUserAction = "disable"
)

func (s AppService) DeleteWorkspaceUser(ctx context.Context, requesterID, workspaceID, userID string, action DeleteWorkspaceUserAction, reassignToUserID *string) error {
	if _, err := s.requireWorkspaceRole(ctx, requesterID, workspaceID, "admin"); err != nil {
		return err
	}
	if requesterID == userID {
		return fmt.Errorf("cannot modify yourself")
	}

	member, err := s.Repo.GetWorkspaceMember(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if member.Role == "admin" {
		var adminCount int64
		if err := s.Repo.DB.WithContext(ctx).
			Model(&models.WorkspaceMember{}).
			Where("workspace_id = ? AND role = ?", workspaceID, "admin").
			Count(&adminCount).Error; err != nil {
			return fmt.Errorf("count admins: %w", err)
		}
		if adminCount <= 1 {
			return fmt.Errorf("cannot remove last admin")
		}
	}

	switch action {
	case DeleteWorkspaceUserActionDisable:
		_, err := s.Repo.UpdateUser(ctx, userID, map[string]any{
			"disabled":   true,
			"updated_at": time.Now().UTC(),
		})
		return err
	case DeleteWorkspaceUserActionUnassign:
		now := time.Now().UTC()
		if err := s.Repo.DB.WithContext(ctx).Exec(`
			UPDATE tasks
			SET assignee_id = NULL, updated_at = ?
			WHERE assignee_id = ?
			AND project_id IN (SELECT id FROM projects WHERE workspace_id = ?)
		`, now, userID, workspaceID).Error; err != nil {
			return fmt.Errorf("unassign tasks: %w", err)
		}
		return s.Repo.DeleteWorkspaceMember(ctx, workspaceID, userID)
	case DeleteWorkspaceUserActionReassign:
		if reassignToUserID == nil || *reassignToUserID == "" {
			return fmt.Errorf("reassignToUserId required")
		}
		targetID := *reassignToUserID
		if targetID == userID {
			return fmt.Errorf("invalid reassign target")
		}
		if _, err := s.Repo.GetWorkspaceMember(ctx, workspaceID, targetID); err != nil {
			return err
		}
		u, err := s.Repo.GetUserByID(ctx, targetID)
		if err != nil {
			return err
		}
		if u.Disabled {
			return fmt.Errorf("cannot reassign to disabled user")
		}
		now := time.Now().UTC()
		if err := s.Repo.DB.WithContext(ctx).Exec(`
			UPDATE tasks
			SET assignee_id = ?, updated_at = ?
			WHERE assignee_id = ?
			AND project_id IN (SELECT id FROM projects WHERE workspace_id = ?)
		`, targetID, now, userID, workspaceID).Error; err != nil {
			return fmt.Errorf("reassign tasks: %w", err)
		}
		return s.Repo.DeleteWorkspaceMember(ctx, workspaceID, userID)
	default:
		return fmt.Errorf("invalid action")
	}
}
