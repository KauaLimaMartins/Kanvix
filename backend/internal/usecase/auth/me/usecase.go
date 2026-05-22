package me

import (
	"context"
	"fmt"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/auth/session"
)

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AvatarColor string `json:"avatarColor"`
	Role        string `json:"role"`
}

type Out struct {
	User User `json:"user"`
}

type userRepo interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type sessionStore interface {
	Get(ctx context.Context, key string) (value string, ok bool, err error)
}

type UseCase struct {
	Users    userRepo
	Sessions sessionStore
}

func (uc UseCase) Execute(ctx context.Context, token string) (Out, error) {
	if token == "" {
		return Out{}, domain.ErrForbidden
	}

	userID, ok, err := uc.Sessions.Get(ctx, session.Key(token))
	if err != nil {
		return Out{}, fmt.Errorf("get session: %w", err)
	}
	if !ok || userID == "" {
		return Out{}, domain.ErrForbidden
	}

	u, err := uc.Users.GetByID(ctx, userID)
	if err != nil {
		return Out{}, err
	}
	if u.Disabled {
		return Out{}, domain.ErrForbidden
	}

	return Out{
		User: User{
			ID:          u.ID,
			Email:       u.Email,
			Name:        u.Name,
			AvatarColor: u.AvatarColor,
			Role:        u.Role,
		},
	}, nil
}
