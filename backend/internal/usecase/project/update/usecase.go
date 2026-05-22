package update

import (
	"context"
	"fmt"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID    string
	ProjectID string
	Patch     map[string]any
}

type Out struct {
	Project entity.Project
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
	Update(ctx context.Context, id string, patch map[string]any) (entity.Project, error)
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Projects   projectRepo
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	if len(in.Patch) == 0 {
		return Out{}, fmt.Errorf("empty patch")
	}

	p, err := uc.Projects.GetByID(ctx, in.ProjectID)
	if err != nil {
		return Out{}, err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}

	in.Patch["updated_at"] = time.Now().UTC()
	updated, err := uc.Projects.Update(ctx, in.ProjectID, in.Patch)
	if err != nil {
		return Out{}, err
	}
	return Out{Project: updated}, nil
}

