package update

import (
	"context"
	"fmt"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID      string
	WorkspaceID string
	Name        *string
	Icon        *string
	Color       *string
}

type Out struct {
	Workspace entity.Workspace
}

type workspaceRepo interface {
	GetByID(ctx context.Context, id string) (entity.Workspace, error)
	Update(ctx context.Context, id string, patch map[string]any) (entity.Workspace, error)
}

type memberRepo interface {
	role.MemberRepo
}

type UseCase struct {
	Workspaces workspaceRepo
	Members    memberRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}

	patch := map[string]any{}
	if in.Name != nil {
		patch["name"] = *in.Name
	}
	if in.Icon != nil {
		patch["icon"] = *in.Icon
	}
	if in.Color != nil {
		patch["color"] = *in.Color
	}
	if len(patch) == 0 {
		return Out{}, fmt.Errorf("empty patch")
	}
	patch["updated_at"] = time.Now().UTC()

	w, err := uc.Workspaces.Update(ctx, in.WorkspaceID, patch)
	if err != nil {
		return Out{}, err
	}
	return Out{Workspace: w}, nil
}

