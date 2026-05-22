package register

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	Name     string
}

type Out struct {
	User  User
	Token string
}

type userRepo interface {
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, u entity.User) (entity.User, error)
}

type workspaceRepo interface {
	Create(ctx context.Context, w entity.Workspace) (entity.Workspace, error)
}

type workspaceMemberRepo interface {
	Upsert(ctx context.Context, m entity.WorkspaceMember) error
}

type sessionStore interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type UseCase struct {
	Users      userRepo
	Workspaces workspaceRepo
	Members    workspaceMemberRepo
	Sessions   sessionStore
	SessionTTL time.Duration
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	name := strings.TrimSpace(in.Name)
	if email == "" {
		return Out{}, fmt.Errorf("email required")
	}
	if name == "" {
		return Out{}, fmt.Errorf("name required")
	}
	if len(in.Password) < 8 {
		return Out{}, fmt.Errorf("password must be at least 8 characters")
	}

	n, err := uc.Users.Count(ctx)
	if err != nil {
		return Out{}, fmt.Errorf("count users: %w", err)
	}
	if n != 0 {
		return Out{}, domain.ErrForbidden
	}

	pwHashBytes, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return Out{}, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	u, err := uc.Users.Create(ctx, entity.User{
		ID:           uuid.NewString(),
		Email:        email,
		Name:         name,
		AvatarColor:  "#6366f1",
		Role:         "admin",
		PasswordHash: string(pwHashBytes),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return Out{}, err
	}

	ws, err := uc.Workspaces.Create(ctx, entity.Workspace{
		ID:        uuid.NewString(),
		OwnerID:   u.ID,
		Name:      "My Workspace",
		Icon:      "LayoutGrid",
		Color:     "#6366f1",
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return Out{}, err
	}

	if err := uc.Members.Upsert(ctx, entity.WorkspaceMember{WorkspaceID: ws.ID, UserID: u.ID, Role: "admin", CreatedAt: now}); err != nil {
		return Out{}, fmt.Errorf("create membership: %w", err)
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
