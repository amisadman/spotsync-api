package parkingzone

import (
	"gorm.io/gorm"
)

type Repository interface {
	CreateZone(zone *ParkingZone) error
	GetAllZones() ([]ParkingZone, error)
	GetZoneByID(id uint) (*ParkingZone, error)
	GetActiveReservationsCount(zoneID uint) (int, error)
}

type parkingZoneRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &parkingZoneRepository{db: db}
}

func (r *parkingZoneRepository) CreateZone(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *parkingZoneRepository) GetAllZones() ([]ParkingZone, error) {
	var zones []ParkingZone
	err := r.db.Find(&zones).Error
	return zones, err
}

func (r *parkingZoneRepository) GetZoneByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone
	err := r.db.First(&zone, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &zone, nil
}

func (r *parkingZoneRepository) GetActiveReservationsCount(zoneID uint) (int, error) {
	var count int64
	err := r.db.Table("reservations").Where("zone_id = ? AND status = ?", zoneID, "active").Count(&count).Error
	return int(count), err
}
