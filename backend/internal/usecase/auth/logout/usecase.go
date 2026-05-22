package logout

import (
	"context"
	"fmt"

	"kanvix/backend/internal/usecase/auth/session"
)

type sessionStore interface {
	Delete(ctx context.Context, key string) error
}

type UseCase struct {
	Sessions sessionStore
}

func (uc UseCase) Execute(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	if err := uc.Sessions.Delete(ctx, session.Key(token)); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

