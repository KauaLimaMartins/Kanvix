package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
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
	}
	var rows []row
	if err := s.Repo.DB.WithContext(ctx).Raw(`
		SELECT u.id, u.email, u.name, u.avatar_color, wm.role
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
