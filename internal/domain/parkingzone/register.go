package parkingzone

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	jwtService := auth.NewJWTService(cfg.JwtSecret)

	adminGroup := e.Group("/api/v1/zones", middleware.AuthMiddleware(jwtService), middleware.RequireRole("admin"))
	adminGroup.POST("", handler.CreateZone)

	publicGroup := e.Group("/api/v1/zones")
	publicGroup.GET("", handler.GetAllZones)
	publicGroup.GET("/:id", handler.GetZoneByID)
}
