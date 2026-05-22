package list

import (
	"context"

	"kanvix/backend/internal/domain/entity"
)

type Out struct {
	Workspaces []entity.Workspace
}

type membershipRepo interface {
	ListWorkspacesForUser(ctx context.Context, userID string) ([]entity.Workspace, error)
}

type UseCase struct {
	Memberships membershipRepo
}

func (uc UseCase) Execute(ctx context.Context, userID string) (Out, error) {
	ws, err := uc.Memberships.ListWorkspacesForUser(ctx, userID)
	if err != nil {
		return Out{}, err
	}
	return Out{Workspaces: ws}, nil
}

