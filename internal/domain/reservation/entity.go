package reservation

import (
	"spotsync/internal/domain/parkingzone"
	"spotsync/internal/domain/user"
	"time"
)

type Reservation struct {
	ID           uint                    `json:"id" gorm:"primaryKey"`
	UserID       uint                    `json:"user_id" gorm:"not null"`
	User         user.User               `json:"user" gorm:"foreignKey:UserID"`
	ZoneID       uint                    `json:"zone_id" gorm:"not null"`
	Zone         parkingzone.ParkingZone `json:"zone" gorm:"foreignKey:ZoneID"`
	LicensePlate string                  `json:"license_plate" gorm:"type:varchar(15);not null"`
	Status       string                  `json:"status" gorm:"type:varchar(20);default:'active';not null"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}
