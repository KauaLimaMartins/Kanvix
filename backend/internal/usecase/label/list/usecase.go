package list

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID      string
	WorkspaceID string
}

type Out struct {
	Labels []entity.Label
}

type labelRepo interface {
	ListByWorkspace(ctx context.Context, workspaceID string) ([]entity.Label, error)
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
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin", "member"); err != nil {
		return Out{}, err
	}
	labels, err := uc.Labels.ListByWorkspace(ctx, in.WorkspaceID)
	if err != nil {
		return Out{}, err
	}
	return Out{Labels: labels}, nil
}

