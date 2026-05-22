package delete_from_workspace

import (
	"context"
	"fmt"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type Action string

const (
	ActionUnassign Action = "unassign"
	ActionReassign Action = "reassign"
	ActionDisable  Action = "disable"
)

type In struct {
	UserID       string
	WorkspaceID  string
	TargetUserID string
	Action       Action
	ReassignToUserID *string
}

type membershipRepo interface {
	Get(ctx context.Context, workspaceID, userID string) (entity.WorkspaceMember, error)
	CountAdmins(ctx context.Context, workspaceID string) (int64, error)
	Delete(ctx context.Context, workspaceID, userID string) error
}

type userRepo interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
	Update(ctx context.Context, id string, patch map[string]any) (entity.User, error)
}

type taskRepo interface {
	UnassignForWorkspace(ctx context.Context, workspaceID, userID string, updatedAt time.Time) error
	ReassignForWorkspace(ctx context.Context, workspaceID, fromUserID, toUserID string, updatedAt time.Time) error
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
	Tasks       taskRepo
	Members     memberRepo
	Workspaces  workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) error {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin"); err != nil {
		return err
	}
	if in.UserID == in.TargetUserID {
		return fmt.Errorf("cannot modify yourself")
	}

	member, err := uc.Memberships.Get(ctx, in.WorkspaceID, in.TargetUserID)
	if err != nil {
		return err
	}
	if member.Role == "admin" {
		adminCount, err := uc.Memberships.CountAdmins(ctx, in.WorkspaceID)
		if err != nil {
			return err
		}
		if adminCount <= 1 {
			return fmt.Errorf("cannot remove last admin")
		}
	}

	switch in.Action {
	case ActionDisable:
		_, err := uc.Users.Update(ctx, in.TargetUserID, map[string]any{
			"disabled":   true,
			"updated_at": time.Now().UTC(),
		})
		return err
	case ActionUnassign:
		now := time.Now().UTC()
		if err := uc.Tasks.UnassignForWorkspace(ctx, in.WorkspaceID, in.TargetUserID, now); err != nil {
			return err
		}
		return uc.Memberships.Delete(ctx, in.WorkspaceID, in.TargetUserID)
	case ActionReassign:
		if in.ReassignToUserID == nil || *in.ReassignToUserID == "" {
			return fmt.Errorf("reassignToUserId required")
		}
		targetID := *in.ReassignToUserID
		if targetID == in.TargetUserID {
			return fmt.Errorf("invalid reassign target")
		}
		if _, err := uc.Memberships.Get(ctx, in.WorkspaceID, targetID); err != nil {
			return err
		}
		u, err := uc.Users.GetByID(ctx, targetID)
		if err != nil {
			return err
		}
		if u.Disabled {
			return fmt.Errorf("cannot reassign to disabled user")
		}
		now := time.Now().UTC()
		if err := uc.Tasks.ReassignForWorkspace(ctx, in.WorkspaceID, in.TargetUserID, targetID, now); err != nil {
			return err
		}
		return uc.Memberships.Delete(ctx, in.WorkspaceID, in.TargetUserID)
	default:
		return fmt.Errorf("invalid action")
	}
}

