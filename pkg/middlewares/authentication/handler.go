package authentication

import (
	"backend/pkg/authentication"
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type Middleware struct {
	authenticator *authentication.TelegramMiniApp
	logger        *slog.Logger
}

func NewAuthenticationMiddleware(secret string, storage authentication.TelegramAuthenticationStorage, logger *slog.Logger) *Middleware {
	return &Middleware{
		authentication.NewTelegramMiniAppAuthenticator(secret, storage),
		logger,
	}
}

func (a *Middleware) Chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slice := strings.Split(r.Header.Get("Authorization"), " ")
		if len(slice) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization header invalid"))
			return
		}
		tma := slice[1]
		user, err := a.authenticator.Authenticate(tma)
		if err != nil {
			a.logger.Info(err.Error(), "authorization", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))

			return
		}
		ctx := context.WithValue(context.Background(), authentication.UserContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
