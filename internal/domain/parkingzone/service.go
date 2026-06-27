package parkingzone

import (
	"errors"
	"spotsync/internal/domain/parkingzone/dto"
)

var ErrZoneNotFound = errors.New("parking zone not found")

type Service interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneCreateResponse, error)
	GetAllZones() ([]dto.ZoneDetailsResponse, error)
	GetZoneByID(id uint) (*dto.ZoneDetailsResponse, error)
	GetActiveCount(zoneID uint) (int, error)
}

type parkingZoneService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &parkingZoneService{repo: repo}
}

func (s *parkingZoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneCreateResponse, error) {
	zone := ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.CreateZone(&zone); err != nil {
		return nil, err
	}

	return &dto.ZoneCreateResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}, nil
}

func (s *parkingZoneService) GetAllZones() ([]dto.ZoneDetailsResponse, error) {
	zones, err := s.repo.GetAllZones()
	if err != nil {
		return nil, err
	}

	res := make([]dto.ZoneDetailsResponse, 0, len(zones))
	for _, z := range zones {
		activeCount, err := s.repo.GetActiveReservationsCount(z.ID)
		if err != nil {
			activeCount = 0
		}

		available := z.TotalCapacity - activeCount
		if available < 0 {
			available = 0
		}

		res = append(res, dto.ZoneDetailsResponse{
			ID:             z.ID,
			Name:           z.Name,
			Type:           z.Type,
			TotalCapacity:  z.TotalCapacity,
			AvailableSpots: available,
			PricePerHour:   z.PricePerHour,
			CreatedAt:      z.CreatedAt,
		})
	}

	return res, nil
}

func (s *parkingZoneService) GetZoneByID(id uint) (*dto.ZoneDetailsResponse, error) {
	z, err := s.repo.GetZoneByID(id)
	if err != nil {
		return nil, err
	}
	if z == nil {
		return nil, ErrZoneNotFound
	}

	activeCount, err := s.repo.GetActiveReservationsCount(z.ID)
	if err != nil {
		activeCount = 0
	}

	available := z.TotalCapacity - activeCount
	if available < 0 {
		available = 0
	}

	return &dto.ZoneDetailsResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt,
	}, nil
}

func (s *parkingZoneService) GetActiveCount(zoneID uint) (int, error) {
	return s.repo.GetActiveReservationsCount(zoneID)
}
