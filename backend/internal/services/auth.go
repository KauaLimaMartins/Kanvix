package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
	"kanvix/backend/internal/repositories/gormrepo"
)

type AuthService struct {
	Repo       gormrepo.Repo
	Redis      *redis.Client
	SessionTTL time.Duration
}

type sessionData struct {
	UserID string `json:"userId"`
}

func (s AuthService) NeedsFirstSignup(ctx context.Context) (bool, error) {
	var n int64
	if err := s.Repo.DB.WithContext(ctx).Model(&models.User{}).Count(&n).Error; err != nil {
		return false, fmt.Errorf("count users: %w", err)
	}
	return n == 0, nil
}

func (s AuthService) FirstSignup(ctx context.Context, email, password, name string) (models.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)
	if email == "" {
		return models.User{}, "", fmt.Errorf("email required")
	}
	if name == "" {
		return models.User{}, "", fmt.Errorf("name required")
	}
	if len(password) < 8 {
		return models.User{}, "", fmt.Errorf("password must be at least 8 characters")
	}

	need, err := s.NeedsFirstSignup(ctx)
	if err != nil {
		return models.User{}, "", err
	}
	if !need {
		return models.User{}, "", repositories.ErrForbidden
	}

	pwHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, "", fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	u := models.User{
		ID:           uuid.NewString(),
		Email:        email,
		Name:         name,
		AvatarColor:  "#6366f1",
		Role:         "admin",
		PasswordHash: string(pwHashBytes),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	u, err = s.Repo.CreateUser(ctx, u)
	if err != nil {
		return models.User{}, "", err
	}

	ws := models.Workspace{
		ID:        uuid.NewString(),
		OwnerID:   u.ID,
		Name:      "My Workspace",
		Icon:      "LayoutGrid",
		Color:     "#6366f1",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := s.Repo.CreateWorkspace(ctx, ws); err != nil {
		return models.User{}, "", err
	}
	m := models.WorkspaceMember{WorkspaceID: ws.ID, UserID: u.ID, Role: "admin", CreatedAt: now}
	if err := s.Repo.DB.WithContext(ctx).Create(&m).Error; err != nil {
		return models.User{}, "", fmt.Errorf("create membership: %w", err)
	}

	token, err := newSessionToken()
	if err != nil {
		return models.User{}, "", err
	}
	if err := s.Redis.Set(ctx, sessionKey(token), u.ID, s.SessionTTL).Err(); err != nil {
		return models.User{}, "", fmt.Errorf("store session: %w", err)
	}
	return u, token, nil
}

func (s AuthService) Login(ctx context.Context, email string, password string) (models.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return models.User{}, "", fmt.Errorf("email required")
	}

	u, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, "", err
	}
	if u.PasswordHash == "" {
		return models.User{}, "", repositories.ErrForbidden
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return models.User{}, "", repositories.ErrForbidden
	}

	token, err := newSessionToken()
	if err != nil {
		return models.User{}, "", err
	}

	if err := s.Redis.Set(ctx, sessionKey(token), u.ID, s.SessionTTL).Err(); err != nil {
		return models.User{}, "", fmt.Errorf("store session: %w", err)
	}

	return u, token, nil
}

func (s AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	if err := s.Redis.Del(ctx, sessionKey(token)).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (s AuthService) Me(ctx context.Context, token string) (models.User, error) {
	if token == "" {
		return models.User{}, repositories.ErrForbidden
	}
	userID, err := s.Redis.Get(ctx, sessionKey(token)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.User{}, repositories.ErrForbidden
		}
		return models.User{}, fmt.Errorf("get session: %w", err)
	}
	if userID == "" {
		return models.User{}, repositories.ErrForbidden
	}
	u, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func sessionKey(token string) string { return "sess:" + token }

func newSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
