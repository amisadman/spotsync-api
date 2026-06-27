package user

import (
	"fmt"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")
var ErrEmailAlreadyExists = fmt.Errorf("email already exists")

type Service interface {
	CreateUser(req dto.CreateRequest) (*dto.RegisterResponseData, error)
	LoginUser(req dto.LoginRequest) (*dto.LoginResponseData, error)
	GetUserByID(id uint) (*User, error)
}

type userService struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) Service {
	return &userService{repo: repo, jwtService: jwtService}
}

func (s *userService) CreateUser(req dto.CreateRequest) (*dto.RegisterResponseData, error) {
	existing, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	role := req.Role
	if role == "" {
		role = "driver"
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  role,
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &dto.RegisterResponseData{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) LoginUser(req dto.LoginRequest) (*dto.LoginResponseData, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.checkPassword(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponseData{
		Token: token,
		User: dto.UserShortResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *userService) GetUserByID(id uint) (*User, error) {
	return s.repo.GetUserByID(id)
}
