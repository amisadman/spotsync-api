package user

import (
	"net/http"
	"spotsync/internal/domain/user/dto"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Invalid request payload",
			"errors":  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	res, err := h.svc.CreateUser(req)
	if err != nil {
		if err == ErrEmailAlreadyExists {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"success": false,
				"message": "Registration failed",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Registration failed",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"message": "User registered successfully",
		"data":    res,
	})
}

func (h *Handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Invalid request payload",
			"errors":  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	res, err := h.svc.LoginUser(req)
	if err != nil {
		if err == ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Login failed",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Login failed",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Login successful",
		"data":    res,
	})
}

func (h *Handler) RefreshToken(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]any{
		"success": false,
		"message": "Refresh token not fully supported",
	})
}

func (h *Handler) GetMe(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Unauthorized",
			"errors":  "missing user id in context",
		})
	}

	user, err := h.svc.GetUserByID(userID)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "User not found",
			"errors":  "unable to retrieve user details",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "User details retrieved successfully",
		"data": dto.UserShortResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}