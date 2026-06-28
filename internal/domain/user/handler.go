package user

import (
	"net/http"
	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

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
		return c.JSON(http.StatusBadRequest, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	res, err := h.svc.CreateUser(req)
	if err != nil {
		if err == ErrEmailAlreadyExists {
			return c.JSON(http.StatusBadRequest, httpresponse.APIErrorResponse{
				Success: false,
				Message: "Registration failed",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Registration failed",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    res,
	})
}

func (h *Handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	res, err := h.svc.LoginUser(req)
	if err != nil {
		if err == ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, httpresponse.APIErrorResponse{
				Success: false,
				Message: "Login failed",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Login failed",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    res,
	})
}

func (h *Handler) RefreshToken(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, httpresponse.APIErrorResponse{
		Success: false,
		Message: "Refresh token not fully supported",
	})
}

func (h *Handler) GetMe(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user id in context",
		})
	}

	user, err := h.svc.GetUserByID(userID)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, httpresponse.APIErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  "unable to retrieve user details",
		})
	}

	return c.JSON(http.StatusOK, httpresponse.APIResponse{
		Success: true,
		Message: "User details retrieved successfully",
		Data: dto.UserShortResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}