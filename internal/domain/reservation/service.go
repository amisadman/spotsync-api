package reservation

import (
	"errors"
	"spotsync/internal/domain/reservation/dto"
)

var (
	ErrForbidden           = errors.New("forbidden: you cannot perform this action")
	ErrReservationNotFound = errors.New("reservation not found")
	ErrCannotCancel        = errors.New("cannot cancel reservation that is not active")
)

type Service interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	GetAllReservations() ([]dto.AdminReservationResponse, error)
	CancelReservation(userID uint, userRole string, reservationID uint) error
}

type reservationService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &reservationService{repo: repo}
}

func (s *reservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	res, err := s.repo.CreateReservationAtomic(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		res = append(res, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneShortDetails{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt:    r.CreatedAt,
		})
	}

	return res, nil
}

func (s *reservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]dto.AdminReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		res = append(res, dto.AdminReservationResponse{
			ID:           r.ID,
			UserID:       r.UserID,
			User: dto.UserShortDetails{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			ZoneID:       r.ZoneID,
			Zone: dto.ZoneShortDetails{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}

	return res, nil
}

func (s *reservationService) CancelReservation(userID uint, userRole string, reservationID uint) error {
	res, err := s.repo.GetByID(reservationID)
	if err != nil {
		return err
	}
	if res == nil {
		return ErrReservationNotFound
	}

	if userRole != "admin" && res.UserID != userID {
		return ErrForbidden
	}

	if res.Status != "active" {
		return ErrCannotCancel
	}

	return s.repo.UpdateStatus(reservationID, "cancelled")
}
