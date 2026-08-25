package handler

import (
	"context"
	"net/http"
)

type requestIDKey struct{}

func withID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func requestID(r *http.Request) string {
	if value, ok := r.Context().Value(requestIDKey{}).(string); ok {
		return value
	}
	return r.Header.Get("X-Request-ID")
}
