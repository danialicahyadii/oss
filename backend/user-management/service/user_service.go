package service

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"oss.kftd.co.id/v2/main/shared/utils"
	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/dto/response"
	"oss.kftd.co.id/v2/user-management/model"
	"oss.kftd.co.id/v2/user-management/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUsers(
	ctx context.Context,
	req request.UserListRequest,
) (*response.UserListResponse, error) {

	page := req.Page

	if page <= 0 {
		page = 1
	}

	limit := req.Limit

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, total, err := s.userRepo.FindAll(
		ctx,
		req.Search,
		req.IsActive,
		offset,
		limit,
	)

	if err != nil {
		return nil, err
	}

	data := make([]response.UserResponse, 0, len(users))

	for _, user := range users {
		data = append(data, response.UserResponse{
			UUID:     user.UUID,
			Username: user.Username,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			IsActive: user.IsActive,
		})
	}

	totalPages := int(math.Ceil(
		float64(total) / float64(limit),
	))

	return &response.UserListResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUser(
	ctx context.Context,
	userUUID string,
) (*response.UserResponse, error) {

	user, err := s.userRepo.FindByUUID(
		ctx,
		userUUID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &response.UserResponse{
		// ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		IsActive: user.IsActive,
	}, nil
}

func (s *UserService) CreateUser(
	ctx context.Context,
	req request.CreateUserRequest,
) (*response.UserResponse, error) {

	// Check username
	if _, err := s.userRepo.FindByUsername(
		ctx,
		req.Username,
	); err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check email
	if _, err := s.userRepo.FindByEmail(
		ctx,
		req.Email,
	); err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	password, err := utils.HashPassword(
		req.Password,
	)

	if err != nil {
		return nil, err
	}

	isActive := true

	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user := &model.User{
		UUID:     uuid.New().String(),
		Username: req.Username,
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Password: password,
		IsActive: isActive,
	}

	if err := s.userRepo.Create(
		ctx,
		user,
	); err != nil {
		return nil, err
	}

	return &response.UserResponse{
		// ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		IsActive: user.IsActive,
	}, nil
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	userUUID string,
	req request.UpdateUserRequest,
) (*response.UserResponse, error) {

	user, err := s.userRepo.FindByUUID(
		ctx,
		userUUID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Email != "" && req.Email != user.Email {

		existing, err := s.userRepo.FindByEmail(
			ctx,
			req.Email,
		)

		if err == nil {
			if existing.UUID != user.UUID {
				return nil, errors.New("email already exists")
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		updates["email"] = req.Email
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := s.userRepo.Update(
			ctx,
			userUUID,
			updates,
		); err != nil {
			return nil, err
		}
	}

	user, err = s.userRepo.FindByUUID(
		ctx,
		userUUID,
	)

	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		// ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		IsActive: user.IsActive,
	}, nil
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	userUUID string,
) error {

	_, err := s.userRepo.FindByUUID(
		ctx,
		userUUID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}

		return err
	}

	return s.userRepo.Delete(
		ctx,
		userUUID,
	)
}
