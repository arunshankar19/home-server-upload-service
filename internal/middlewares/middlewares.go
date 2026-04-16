package middlewares

import (
	"context"
	"net/http"

	"github.com/arunshankar19/home-server-upload-service/internal/ctxkey"
	"github.com/arunshankar19/home-server-upload-service/internal/response"
)

const (
	userIdHeaderKey = "X-User-Id"
)

// SetUserIdInCtx extracts userId from header and injects it into ctx
func SetUserIdInCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIdHeaderKey)
		if userID == "" {
			response.Fail(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		ctx := context.WithValue(r.Context(), ctxkey.UserID, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
