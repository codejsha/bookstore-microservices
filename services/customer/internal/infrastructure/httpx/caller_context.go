package httpx

import "context"

type callerCtxKey struct{}

type Caller struct {
	UserUid string
	IsAdmin bool
}

func WithCaller(ctx context.Context, caller Caller) context.Context {
	return context.WithValue(ctx, callerCtxKey{}, caller)
}

func CallerFromContext(ctx context.Context) Caller {
	caller, _ := ctx.Value(callerCtxKey{}).(Caller)
	return caller
}
