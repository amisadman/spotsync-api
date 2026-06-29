package server

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"spotsync/internal/config"
	"spotsync/internal/domain/parkingzone"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate (i interface{})error {
	return cv.validator.Struct(i)
}


func Start(db *gorm.DB, cfg *config.Config){
	db.AutoMigrate(&user.User{}, &parkingzone.ParkingZone{}, &reservation.Reservation{})

	e:=  echo.New();
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())
	
	e.GET("/", func(c *echo.Context) error {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		duration := time.Since(startTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		uptimeStr := fmt.Sprintf("%d hours %d minutes", hours, minutes)

		clientIP := c.RealIP()
		accessedAt := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

		type ClientDetails struct {
			ClientIP   string `json:"clientIP"`
			AccessedAt string `json:"accessedAt"`
		}

		type ServerDetails struct {
			Hostname string `json:"hostname"`
			Platform string `json:"platform"`
			Uptime   string `json:"uptime"`
		}

		type WelcomeResponse struct {
			Success       bool          `json:"success"`
			Message       string        `json:"message"`
			Version       string        `json:"version"`
			ClientDetails ClientDetails `json:"clientDetails"`
			ServerDetails ServerDetails `json:"serverDetails"`
		}

		return c.JSON(http.StatusOK, WelcomeResponse{
			Success: true,
			Message: "Welcome to Spotsync",
			Version: "1.0.0",
			ClientDetails: ClientDetails{
				ClientIP:   clientIP,
				AccessedAt: accessedAt,
			},
			ServerDetails: ServerDetails{
				Hostname: hostname,
				Platform: runtime.GOOS,
				Uptime:   uptimeStr,
			},
		})
	})

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})
	user.RegisterRoutes(e, db, cfg)
	parkingzone.RegisterRoutes(e, db, cfg)
	reservation.RegisterRoutes(e, db, cfg)
	
	e.Start(":"+ cfg.Port)
}