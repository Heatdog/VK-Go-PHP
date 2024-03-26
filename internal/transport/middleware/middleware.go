package middleware_transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Heatdog/VK-Go-PHP/internal/transport"
	"github.com/Heatdog/VK-Go-PHP/pkg/jwt"
	"github.com/google/uuid"
)

type Middleware struct {
	logger *slog.Logger
	key    string
}

func NewMiddleware(logger *slog.Logger, key string) *Middleware {
	return &Middleware{
		logger: logger,
		key:    key,
	}
}

func (mid *Middleware) Auth(onlyAuthorize bool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		mid.logger.Debug("verify access token", slog.String("header", header))

		if header == "" {
			mid.logger.Debug("header is empty")
			if onlyAuthorize {
				transport.NewRespWriter(w, "header is empty", http.StatusUnauthorized, mid.logger)
				return
			}

			next(w, r)
			return
		}

		mid.logger.Debug("got access token", slog.String("token", header))
		fields, err := mid.verifyTokenHeader(header)
		if err != nil {
			mid.logger.Debug("auth header err", slog.Any("err", err))

			if onlyAuthorize {
				transport.NewRespWriter(w, err.Error(), http.StatusUnauthorized, mid.logger)
				return
			}

			next(w, r)
			return
		}

		id, err := uuid.Parse(fields.ID)
		if err != nil {
			mid.logger.Debug("auth header err", slog.Any("err", err))

			if onlyAuthorize {
				transport.NewRespWriter(w, err.Error(), http.StatusUnauthorized, mid.logger)
				return
			}

			next(w, r)
			return
		}

		mid.logger.Debug("set user id in context", slog.Any("id", id))

		ctx := context.WithValue(r.Context(), "user_id", id)
		next(w, r.WithContext(ctx))
	}
}

func (mid *Middleware) verifyTokenHeader(header string) (*jwt.TokenFileds, error) {
	mid.logger.Debug("check number of fields", slog.String("header", header))

	headers := strings.Split(header, " ")
	if len(headers) != 2 {
		err := fmt.Errorf("wrong scheame of auth header")
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return nil, err
	}

	mid.logger.Debug("check scheame")

	if headers[0] != "Bearer" {
		err := fmt.Errorf("wrong scheame of auth header")
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return nil, err
	}

	mid.logger.Debug("verify token", slog.String("token", string(header[1])))

	fields, err := jwt.VerifyToken(string(headers[1]), mid.key)
	if err != nil {
		mid.logger.Warn("auth header err", slog.Any("err", err))
		return nil, err
	}
	mid.logger.Debug("res fields", slog.Any("fields", fields))
	return fields, nil
}
