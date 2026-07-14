package auth

import (
	"context"

	"github.com/pafthang/pocketagent/pkgs/models"
)

type contextKey string

const userContextKey contextKey = "space_user"

// WithUser stores the authenticated user in context.
func WithUser(ctx context.Context, user models.AuthUser) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext returns the authenticated user.
func UserFromContext(ctx context.Context) (models.AuthUser, bool) {
	user, ok := ctx.Value(userContextKey).(models.AuthUser)
	return user, ok
}
