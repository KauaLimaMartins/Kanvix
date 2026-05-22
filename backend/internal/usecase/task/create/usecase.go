package create

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID    string
	ProjectID string
	ColumnID  string
	Title     string
}

type Out struct {
	Task entity.Task
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type taskRepo interface {
	CountInColumn(ctx context.Context, columnID string) (int64, error)
	Create(ctx context.Context, t entity.Task) (entity.Task, error)
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Projects   projectRepo
	Tasks      taskRepo
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	p, err := uc.Projects.GetByID(ctx, in.ProjectID)
	if err != nil {
		return Out{}, err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin", "member"); err != nil {
		return Out{}, err
	}

	n, err := uc.Tasks.CountInColumn(ctx, in.ColumnID)
	if err != nil {
		return Out{}, err
	}

	now := time.Now().UTC()
	t, err := uc.Tasks.Create(ctx, entity.Task{
		ID:          uuid.NewString(),
		ProjectID:   in.ProjectID,
		ColumnID:    in.ColumnID,
		Title:       in.Title,
		Description: "",
		Order:       int(n),
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return Out{}, err
	}
	return Out{Task: t}, nil
}

