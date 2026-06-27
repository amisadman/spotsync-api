package reservation

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

	authGroup := e.Group("/api/v1/reservations", middleware.AuthMiddleware(jwtService))
	authGroup.POST("", handler.CreateReservation)
	authGroup.GET("/my-reservations", handler.GetMyReservations)
	authGroup.DELETE("/:id", handler.CancelReservation)

	adminGroup := e.Group("/api/v1/reservations", middleware.AuthMiddleware(jwtService), middleware.RequireRole("admin"))
	adminGroup.GET("", handler.GetAllReservations)
}
