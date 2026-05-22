package update_in_workspace

import (
	"context"
	"fmt"
	"strings"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID      string
	WorkspaceID string
	TargetUserID string
	Name        *string
	Role        *string
}

type membershipRepo interface {
	Get(ctx context.Context, workspaceID, userID string) (entity.WorkspaceMember, error)
	Upsert(ctx context.Context, m entity.WorkspaceMember) error
}

type userRepo interface {
	Update(ctx context.Context, id string, patch map[string]any) (entity.User, error)
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Memberships membershipRepo
	Users       userRepo
	Members     memberRepo
	Workspaces  workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) error {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin"); err != nil {
		return err
	}
	if _, err := uc.Memberships.Get(ctx, in.WorkspaceID, in.TargetUserID); err != nil {
		return err
	}
	if in.Name == nil && in.Role == nil {
		return fmt.Errorf("invalid request")
	}

	patch := map[string]any{"updated_at": time.Now().UTC()}
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if n == "" {
			return fmt.Errorf("invalid name")
		}
		patch["name"] = n
	}
	if in.Role != nil {
		r := strings.TrimSpace(strings.ToLower(*in.Role))
		if r != "admin" && r != "member" {
			return fmt.Errorf("invalid role")
		}
		if err := uc.Memberships.Upsert(ctx, entity.WorkspaceMember{
			WorkspaceID: in.WorkspaceID,
			UserID:      in.TargetUserID,
			Role:        r,
			CreatedAt:   time.Now().UTC(),
		}); err != nil {
			return err
		}
		patch["role"] = r
	}

	_, err := uc.Users.Update(ctx, in.TargetUserID, patch)
	return err
}

