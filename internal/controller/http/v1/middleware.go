package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/passionde/user-segmentation-service/internal/service"
	"net/http"
	"strings"
)

const (
	userIdCtx = "keyId"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (h *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			return nil
		}

		userId, err := h.authService.TokenExist(c.Request().Context(), token)
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, ErrCannotParseToken.Error())
			return err
		}

		c.Set(userIdCtx, userId)

		return next(c)
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
