package auth_utils

import (
	"context"
	"encoding/base64"
)

// ctxKey is a private type to avoid key collisions in context.
type ctxKey string

const (
	ctxAuthKey ctxKey = "auth"
)

// GetAuthClaim extracts AuthClaim from a standard context.Context.
// It supports both pointer and value storage forms and returns zero-value when absent.
func GetAuthClaim(ctx context.Context) AuthClaim {
	if ctx == nil {
		return AuthClaim{}
	}
	v := ctx.Value(ctxAuthKey)
	switch t := v.(type) {
	case *AuthClaim:
		if t != nil {
			return *t
		}
	case AuthClaim:
		return t
	}
	return AuthClaim{}
}

// WithAuthClaim returns a new context with the provided AuthClaim attached.
func WithAuthClaim(ctx context.Context, claims *AuthClaim) context.Context {
	if claims == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxAuthKey, claims)
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
