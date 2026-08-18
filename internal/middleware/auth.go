package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	"kayam-be/internal/repository"
)

type ctxKey string

const userIDCtxKey ctxKey = "user_id"

// RequireAuth is intentionally minimal for this prototype: the bearer token
// IS the user's UUID (issued at login), looked up on every request. No JWT,
// no session store — documented as a prototype-only shortcut in
// kayam-be/.adlc/spec/kayam-be/requirements.md US-2, matching the same
// "mocked, not production-trust" posture as verification elsewhere in this
// phase. Swapping in real signed tokens later touches only this file.
func RequireAuth(queries *repository.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				http.Error(w, `{"error":"missing Authorization: Bearer <token>"}`, http.StatusUnauthorized)
				return
			}

			var userID pgtype.UUID
			if err := userID.Scan(token); err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			if _, err := queries.GetUserByID(r.Context(), userID); err != nil {
				http.Error(w, `{"error":"unknown user"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDCtxKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (pgtype.UUID, bool) {
	id, ok := ctx.Value(userIDCtxKey).(pgtype.UUID)
	return id, ok
}
