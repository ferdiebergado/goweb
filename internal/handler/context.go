package handler

import "context"

type ctxKey int

var paramsCtxKey ctxKey

func NewParamsContext[T any](ctx context.Context, t T) context.Context {
	return context.WithValue(ctx, paramsCtxKey, t)
}

func FromParamsContext[T any](ctx context.Context) (T, bool) {
	t, ok := ctx.Value(paramsCtxKey).(T)
	return t, ok
}
