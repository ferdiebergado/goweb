package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/goexpress"
	"github.com/go-playground/validator/v10"
)

type ctxKey int

var paramsCtxKey ctxKey

func DecodeJSON[T any]() goexpress.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") == jsonCT {
				slog.Info("Decoding json body...")
				var decoded T
				decoder := json.NewDecoder(r.Body)
				decoder.DisallowUnknownFields()
				if err := decoder.Decode(&decoded); err != nil {
					badRequestError(w, r, err)
					return
				}
				ctx := context.WithValue(r.Context(), paramsCtxKey, decoded)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ValidateInput[T any](validate *validator.Validate) goexpress.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("Validating input...")
			ctxVal := r.Context().Value(paramsCtxKey)
			params, ok := ctxVal.(T)

			if !ok {
				var t T
				badRequestError(w, r, fmt.Errorf("cannot type assert context value %v to %T", ctxVal, t))
				return
			}

			if err := validate.Struct(params); err != nil {
				validationError(w, r, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func NewParamsContext[T any](ctx context.Context, t T) context.Context {
	return context.WithValue(ctx, paramsCtxKey, t)
}

func FromParamsContext[T any](ctx context.Context) (T, bool) {
	t, ok := ctx.Value(paramsCtxKey).(T)
	return t, ok
}
