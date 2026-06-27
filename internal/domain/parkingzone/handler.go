package parkingzone

import (
	"net/http"
	"strconv"
	"spotsync/internal/domain/parkingzone/dto"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateZone(c *echo.Context) error {
	var req dto.CreateZoneRequest
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

	res, err := h.svc.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to create parking zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"message": "Parking zone created successfully",
		"data":    res,
	})
}

func (h *Handler) GetAllZones(c *echo.Context) error {
	res, err := h.svc.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to retrieve parking zones",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Parking zones retrieved successfully",
		"data":    res,
	})
}

func (h *Handler) GetZoneByID(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Invalid zone ID",
			"errors":  "ID must be an integer",
		})
	}

	res, err := h.svc.GetZoneByID(uint(id))
	if err != nil {
		if err == ErrZoneNotFound {
			return c.JSON(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Parking zone not found",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to retrieve parking zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Parking zone retrieved successfully",
		"data":    res,
	})
}
