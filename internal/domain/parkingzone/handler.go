package parkingzone

import (
	"net/http"
	"strconv"
	"spotsync/internal/domain/parkingzone/dto"
	"spotsync/internal/httpresponse"

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

	res, err := h.svc.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Failed to create parking zone",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.APIResponse{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    res,
	})
}

func (h *Handler) GetAllZones(c *echo.Context) error {
	res, err := h.svc.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Failed to retrieve parking zones",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.APIResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    res,
	})
}

func (h *Handler) GetZoneByID(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Invalid zone ID",
			Errors:  "ID must be an integer",
		})
	}

	res, err := h.svc.GetZoneByID(uint(id))
	if err != nil {
		if err == ErrZoneNotFound {
			return c.JSON(http.StatusNotFound, httpresponse.APIErrorResponse{
				Success: false,
				Message: "Parking zone not found",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.APIErrorResponse{
			Success: false,
			Message: "Failed to retrieve parking zone",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.APIResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    res,
	})
}
// wait, line 22 is httpcall, line 68 is httpresponse. I must fix them.
