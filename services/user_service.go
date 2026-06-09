package services

import (
	"errors"
	"gk-capital-backend/dto"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
	"gk-capital-backend/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) FindAll(page, pageSize int, search string) (*dto.UserListResponse, error) {
	users, total, err := s.userRepo.FindAll(page, pageSize, search)
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponse
	for _, u := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	// If userResponses is nil, make it an empty slice instead of null
	if userResponses == nil {
		userResponses = []dto.UserResponse{}
	}

	return &dto.UserListResponse{
		Users: userResponses,
		Total: total,
	}, nil
}

func (s *UserService) Create(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	existing, _ := s.userRepo.FindByEmailRaw(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: *req.IsActive,
	}

	err = s.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) Update(id uint, req dto.UpdateUserRequest, currentUserID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Email != req.Email {
		existing, _ := s.userRepo.FindByEmailRaw(req.Email)
		if existing != nil {
			return nil, errors.New("email already registered")
		}
	}

	// Security check for SUPER_ADMIN role deactivation or role change
	if user.Role == "SUPER_ADMIN" && (req.Role != "SUPER_ADMIN" || !*req.IsActive) {
		superAdminCount, err := s.userRepo.CountSuperAdmins()
		if err != nil {
			return nil, err
		}
		if superAdminCount <= 1 {
			return nil, errors.New("cannot deactivate or change role of the last active super admin")
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	user.IsActive = *req.IsActive

	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) Delete(id uint, currentUserID uint) error {
	if id == currentUserID {
		return errors.New("cannot delete your own account")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role == "SUPER_ADMIN" && user.IsActive {
		superAdminCount, err := s.userRepo.CountSuperAdmins()
		if err != nil {
			return err
		}
		if superAdminCount <= 1 {
			return errors.New("cannot delete the last active super admin")
		}
	}

	return s.userRepo.Delete(id)
}
