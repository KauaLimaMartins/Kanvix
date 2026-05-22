package create_in_workspace

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID      string
	WorkspaceID string
	Email       string
	Password    string
	Name        string
	Role        string
}

type Out struct {
	User entity.WorkspaceUser
}

type userRepo interface {
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Create(ctx context.Context, u entity.User) (entity.User, error)
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type membershipRepo interface {
	Upsert(ctx context.Context, m entity.WorkspaceMember) error
}

type UseCase struct {
	Users       userRepo
	Members     memberRepo
	Workspaces  workspaceRepo
	Memberships membershipRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}

	email := strings.TrimSpace(strings.ToLower(in.Email))
	name := strings.TrimSpace(in.Name)
	roleStr := strings.TrimSpace(strings.ToLower(in.Role))
	if email == "" || name == "" {
		return Out{}, fmt.Errorf("%w: email and name required", domain.ErrBadRequest)
	}
	if roleStr != "admin" && roleStr != "member" {
		roleStr = "member"
	}
	if len(in.Password) < 8 {
		return Out{}, fmt.Errorf("%w: password must be at least 8 characters", domain.ErrBadRequest)
	}

	if _, err := uc.Users.GetByEmail(ctx, email); err == nil {
		return Out{}, fmt.Errorf("%w: email already exists", domain.ErrConflict)
	} else if !errors.Is(err, domain.ErrNotFound) {
		return Out{}, err
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
		AvatarColor:  "#64748b",
		Role:         roleStr,
		PasswordHash: string(pwHashBytes),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return Out{}, err
	}

	if err := uc.Memberships.Upsert(ctx, entity.WorkspaceMember{
		WorkspaceID: in.WorkspaceID,
		UserID:      u.ID,
		Role:        roleStr,
		CreatedAt:   now,
	}); err != nil {
		return Out{}, err
	}

	return Out{User: entity.WorkspaceUser{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		AvatarColor: u.AvatarColor,
		Role:        roleStr,
		Disabled:    u.Disabled,
	}}, nil
}
