package httpx

import "context"

type bearerCtxKey struct{}

func WithBearerToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, bearerCtxKey{}, token)
}

func BearerTokenFromContext(ctx context.Context) string {
	token, _ := ctx.Value(bearerCtxKey{}).(string)
	return token
}
