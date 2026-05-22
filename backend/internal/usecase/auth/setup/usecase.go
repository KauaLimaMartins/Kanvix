package setup

import (
	"context"
	"fmt"
)

type Out struct {
	NeedsFirstSignup bool `json:"needsFirstSignup"`
}

type userCounter interface {
	Count(ctx context.Context) (int64, error)
}

type UseCase struct {
	Users userCounter
}

func (uc UseCase) Execute(ctx context.Context) (Out, error) {
	n, err := uc.Users.Count(ctx)
	if err != nil {
		return Out{}, fmt.Errorf("count users: %w", err)
	}
	return Out{NeedsFirstSignup: n == 0}, nil
}

