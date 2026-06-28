package middleware

import (
	"net/http"
	"spotsync/internal/auth"
	"spotsync/internal/httpresponse"
	"strings"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, httpError("Unauthorized", "missing authorization header"))
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, httpError("Unauthorized", "invalid authorization header format"))
			}

			claims, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, httpError("Unauthorized", "invalid or expired token"))
			}

			// reject refresh tokens from being used as access tokens
			if claims.TokenType != auth.TokenTypeAccess {
				return c.JSON(http.StatusUnauthorized, httpError("Unauthorized", "invalid token type"))
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_name", claims.Name)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			roleVal := c.Get("user_role")
			role, ok := roleVal.(string)
			if !ok || role == "" {
				return c.JSON(http.StatusForbidden, httpError("Access denied", "missing user role"))
			}

			for _, allowed := range allowedRoles {
				if allowed == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, httpError("Access denied", "insufficient permissions"))
		}
	}
}

func httpError(message string, errors interface{}) httpresponse.APIErrorResponse {
	return httpresponse.APIErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}