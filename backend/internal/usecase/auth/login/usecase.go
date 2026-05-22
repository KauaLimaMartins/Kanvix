package login

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

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

type In struct {
	Email    string
	Password string
}

type Out struct {
	User  User
	Token string
}

type userRepo interface {
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type sessionStore interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type UseCase struct {
	Users      userRepo
	Sessions   sessionStore
	SessionTTL time.Duration
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" {
		return Out{}, fmt.Errorf("email required")
	}

	u, err := uc.Users.GetByEmail(ctx, email)
	if err != nil {
		return Out{}, err
	}
	if u.Disabled {
		return Out{}, domain.ErrForbidden
	}
	if u.PasswordHash == "" {
		return Out{}, domain.ErrForbidden
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return Out{}, domain.ErrForbidden
	}

	token, err := session.NewToken()
	if err != nil {
		return Out{}, err
	}
	if err := uc.Sessions.Set(ctx, session.Key(token), u.ID, uc.SessionTTL); err != nil {
		return Out{}, fmt.Errorf("store session: %w", err)
	}

	return Out{
		User: User{
			ID:          u.ID,
			Email:       u.Email,
			Name:        u.Name,
			AvatarColor: u.AvatarColor,
			Role:        u.Role,
		},
		Token: token,
	}, nil
}
