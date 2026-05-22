package create

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID      string
	WorkspaceID string
	Name        string
	Color       string
}

type Out struct {
	Label entity.Label
}

type labelRepo interface {
	Create(ctx context.Context, l entity.Label) (entity.Label, error)
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Labels     labelRepo
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}
	now := time.Now().UTC()
	l, err := uc.Labels.Create(ctx, entity.Label{
		ID:          uuid.NewString(),
		WorkspaceID: in.WorkspaceID,
		Name:        in.Name,
		Color:       in.Color,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return Out{}, err
	}
	return Out{Label: l}, nil
}

