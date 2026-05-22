package update

import (
	"context"
	"fmt"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID  string
	LabelID string
	Patch   map[string]any
}

type Out struct {
	Label entity.Label
}

type labelRepo interface {
	GetByID(ctx context.Context, id string) (entity.Label, error)
	Update(ctx context.Context, id string, patch map[string]any) (entity.Label, error)
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
	if len(in.Patch) == 0 {
		return Out{}, fmt.Errorf("empty patch")
	}
	l, err := uc.Labels.GetByID(ctx, in.LabelID)
	if err != nil {
		return Out{}, err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, l.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}
	in.Patch["updated_at"] = time.Now().UTC()
	updated, err := uc.Labels.Update(ctx, in.LabelID, in.Patch)
	if err != nil {
		return Out{}, err
	}
	return Out{Label: updated}, nil
}

