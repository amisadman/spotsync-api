package reservation

import (
	"net/http"
	"strconv"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpresponse"

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
		return c.JSON(http.StatusUnauthorized, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user id in context",
		})
	}

	var req dto.CreateReservationRequest
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

	res, err := h.svc.CreateReservation(userID, req)
	if err != nil {
		if err == ErrZoneFull {
			return c.JSON(http.StatusConflict, httpresponse.APIErrorResponse{
				Success: false,
				Message: "Reservation failed",
				Errors:  err.Error(),
			})
		}
		if err == ErrZoneNotFound {
			return c.JSON(http.StatusNotFound, httpresponse.APIErrorResponse{
				Success: false,
				Message: "Reservation failed",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Reservation failed",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.APIResponse{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    res,
	})
}

func (h *Handler) GetMyReservations(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user id in context",
		})
	}

	res, err := h.svc.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Failed to retrieve reservations",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.APIResponse{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    res,
	})
}

func (h *Handler) GetAllReservations(c *echo.Context) error {
	res, err := h.svc.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Failed to retrieve all reservations",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.APIResponse{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    res,
	})
}

func (h *Handler) CancelReservation(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user id in context",
		})
	}

	userRoleVal := c.Get("user_role")
	userRole, ok := userRoleVal.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user role in context",
		})
	}

	idParam := c.Param("id")
	resID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Invalid reservation ID",
			Errors:  "ID must be an integer",
		})
	}

	err = h.svc.CancelReservation(userID, userRole, uint(resID))
	if err != nil {
		if err == ErrReservationNotFound {
			return c.JSON(http.StatusNotFound, httpresponse.APIErrorResponse{
				Success: false,
				Message: "Failed to cancel reservation",
				Errors:  err.Error(),
			})
		}
		if err == ErrForbidden {
			return c.JSON(http.StatusForbidden, httpresponse.APIErrorResponse{
				Success: false,
				Message: "Failed to cancel reservation",
				Errors:  err.Error(),
			})
		}
		if err == ErrCannotCancel {
			return c.JSON(http.StatusConflict, httpresponse.APIErrorResponse{
				Success: false,
				Message: "Failed to cancel reservation",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Failed to cancel reservation",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.APIResponse{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}
