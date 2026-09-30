package middleware

import (
	"context"
	"net/http"
	"strings"

	"andurel-site/models"
	"andurel-site/services"

	"github.com/labstack/echo/v5"
)

type apiTokenContextKey struct{}

func APITokenFromContext(ctx context.Context) (models.Token, bool) {
	token, ok := ctx.Value(apiTokenContextKey{}).(models.Token)
	return token, ok
}

func APITokenAuth(tokens services.APITokens) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			parts := strings.Fields(header)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}

			token, err := tokens.Authenticate(c.Request().Context(), parts[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			ctx := context.WithValue(c.Request().Context(), apiTokenContextKey{}, token)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
