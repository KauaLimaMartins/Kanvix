package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

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
	Email  string `json:"email"`
}

func (s AuthService) Login(ctx context.Context, email string) (models.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return models.User{}, "", fmt.Errorf("email required")
	}

	u, err := s.Repo.GetUserByID(ctx, "u1")
	if err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			return models.User{}, "", err
		}
		u = models.User{
			ID:          "u1",
			Email:       "you@demo.local",
			Name:        "You",
			AvatarColor: "#6366f1",
		}
		u, err = s.Repo.CreateUser(ctx, u)
		if err != nil {
			return models.User{}, "", err
		}
	}
	u.Email = email

	token, err := newSessionToken()
	if err != nil {
		return models.User{}, "", err
	}

	key := sessionKey(token)
	raw, _ := json.Marshal(sessionData{UserID: u.ID, Email: email})
	if err := s.Redis.Set(ctx, key, raw, s.SessionTTL).Err(); err != nil {
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
	raw, err := s.Redis.Get(ctx, sessionKey(token)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.User{}, repositories.ErrForbidden
		}
		return models.User{}, fmt.Errorf("get session: %w", err)
	}
	var sdata sessionData
	if err := json.Unmarshal(raw, &sdata); err != nil || sdata.UserID == "" {
		return models.User{}, repositories.ErrForbidden
	}
	u, err := s.Repo.GetUserByID(ctx, sdata.UserID)
	if err != nil {
		return models.User{}, err
	}
	if sdata.Email != "" {
		u.Email = sdata.Email
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
