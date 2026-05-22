package create

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
)

type In struct {
	UserID string
	Name   string
	Icon   *string
	Color  *string
}

type Out struct {
	Workspace entity.Workspace
}

type userRepo interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type workspaceRepo interface {
	Create(ctx context.Context, w entity.Workspace) (entity.Workspace, error)
}

type memberRepo interface {
	Upsert(ctx context.Context, m entity.WorkspaceMember) error
}

type UseCase struct {
	Users      userRepo
	Workspaces workspaceRepo
	Members    memberRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	u, err := uc.Users.GetByID(ctx, in.UserID)
	if err != nil {
		return Out{}, err
	}
	if u.Role != "admin" {
		return Out{}, domain.ErrForbidden
	}

	now := time.Now().UTC()
	w := entity.Workspace{
		ID:        uuid.NewString(),
		OwnerID:   in.UserID,
		Name:      in.Name,
		Icon:      "LayoutGrid",
		Color:     "#6366f1",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if in.Icon != nil && *in.Icon != "" {
		w.Icon = *in.Icon
	}
	if in.Color != nil && *in.Color != "" {
		w.Color = *in.Color
	}

	created, err := uc.Workspaces.Create(ctx, w)
	if err != nil {
		return Out{}, err
	}
	if err := uc.Members.Upsert(ctx, entity.WorkspaceMember{
		WorkspaceID: created.ID,
		UserID:      in.UserID,
		Role:        "admin",
		CreatedAt:   now,
	}); err != nil {
		return Out{}, fmt.Errorf("create membership: %w", err)
	}

	return Out{Workspace: created}, nil
}
