package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"kayam-be/internal/repository"
)

type AuthService struct {
	Queries *repository.Queries
}

// Login is deliberately minimal for this prototype phase — email only, no
// password, "chosen for prototype speed, not production trust" per
// kayam-be/.adlc/spec/kayam-be/requirements.md US-2. Creates the user on
// first login. The returned user's ID doubles as the bearer token (see
// internal/middleware/auth.go) — an explicit, documented shortcut, not
// something to carry into a real auth system.
func (s *AuthService) Login(ctx context.Context, name, email string) (repository.User, error) {
	user, err := s.Queries.GetUserByEmail(ctx, &email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return repository.User{}, fmt.Errorf("lookup user: %w", err)
		}
		user, err = s.Queries.CreateUser(ctx, repository.CreateUserParams{
			Name:       name,
			Email:      &email,
			AuthMethod: "email",
		})
		if err != nil {
			return repository.User{}, fmt.Errorf("create user: %w", err)
		}
	}

	if err := s.Queries.RecordLogin(ctx, user.ID); err != nil {
		return repository.User{}, fmt.Errorf("record login: %w", err)
	}

	// Re-fetch rather than returning the pre-RecordLogin `user` value — the
	// caller (and any FE showing login_count) needs this login reflected,
	// not the state from before it.
	user, err = s.Queries.GetUserByID(ctx, user.ID)
	if err != nil {
		return repository.User{}, fmt.Errorf("refetch user after login: %w", err)
	}

	return user, nil
}
