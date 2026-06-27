package reservation

import (
	"errors"
	"spotsync/internal/domain/parkingzone"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrZoneFull = errors.New("parking zone is full")
var ErrZoneNotFound = errors.New("parking zone not found")

type Repository interface {
	CreateReservationAtomic(userID uint, zoneID uint, licensePlate string) (*Reservation, error)
	GetByID(id uint) (*Reservation, error)
	GetByUserID(userID uint) ([]Reservation, error)
	GetAll() ([]Reservation, error)
	UpdateStatus(id uint, status string) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateReservationAtomic(userID uint, zoneID uint, licensePlate string) (*Reservation, error) {
	var res *Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone parkingzone.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrZoneNotFound
			}
			return err
		}

		var activeCount int64
		if err := tx.Model(&Reservation{}).Where("zone_id = ? AND status = ?", zoneID, "active").Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		newRes := Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		if err := tx.Create(&newRes).Error; err != nil {
			return err
		}

		res = &newRes
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := r.db.Preload("Zone").First(res, res.ID).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (r *reservationRepository) GetByID(id uint) (*Reservation, error) {
	var res Reservation
	err := r.db.Preload("Zone").First(&res, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *reservationRepository) GetByUserID(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) GetAll() ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", status).Error
}
