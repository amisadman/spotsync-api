package reservation

import (
	"net/http"
	"strconv"
	"spotsync/internal/domain/reservation/dto"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateReservation(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Unauthorized",
			"errors":  "missing user id in context",
		})
	}

	var req dto.CreateReservationRequest
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

	res, err := h.svc.CreateReservation(userID, req)
	if err != nil {
		if err == ErrZoneFull {
			return c.JSON(http.StatusConflict, map[string]any{
				"success": false,
				"message": "Reservation failed",
				"errors":  err.Error(),
			})
		}
		if err == ErrZoneNotFound {
			return c.JSON(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Reservation failed",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Reservation failed",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"message": "Reservation confirmed successfully",
		"data":    res,
	})
}

func (h *Handler) GetMyReservations(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Unauthorized",
			"errors":  "missing user id in context",
		})
	}

	res, err := h.svc.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to retrieve reservations",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "My reservations retrieved successfully",
		"data":    res,
	})
}

func (h *Handler) GetAllReservations(c *echo.Context) error {
	res, err := h.svc.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to retrieve all reservations",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "All reservations retrieved successfully",
		"data":    res,
	})
}

func (h *Handler) CancelReservation(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Unauthorized",
			"errors":  "missing user id in context",
		})
	}

	userRoleVal := c.Get("user_role")
	userRole, ok := userRoleVal.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Unauthorized",
			"errors":  "missing user role in context",
		})
	}

	idParam := c.Param("id")
	resID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Invalid reservation ID",
			"errors":  "ID must be an integer",
		})
	}

	err = h.svc.CancelReservation(userID, userRole, uint(resID))
	if err != nil {
		if err == ErrReservationNotFound {
			return c.JSON(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Failed to cancel reservation",
				"errors":  err.Error(),
			})
		}
		if err == ErrForbidden {
			return c.JSON(http.StatusForbidden, map[string]any{
				"success": false,
				"message": "Failed to cancel reservation",
				"errors":  err.Error(),
			})
		}
		if err == ErrCannotCancel {
			return c.JSON(http.StatusConflict, map[string]any{
				"success": false,
				"message": "Failed to cancel reservation",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to cancel reservation",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Reservation cancelled successfully",
	})
}
