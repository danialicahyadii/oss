package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"oss.kftd.co.id/v2/main/shared/utils"
	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/model"
	"oss.kftd.co.id/v2/user-management/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	req request.LoginRequest,
) (string, error) {

	user, err := s.userRepo.FindByUsername(
		ctx,
		req.Username,
	)
	fmt.Println(user)

	if err != nil {
		return "", errors.New(
			"username or password is incorrect",
		)
	}

	if !user.IsActive {
		return "", errors.New("user is inactive")
	}

	if err := utils.CheckPassword(
		user.Password,
		req.Password,
	); err != nil {
		return "", errors.New(
			"username or password is incorrect",
		)
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Username,
		[]string{},
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(
	ctx context.Context,
	req request.RegisterRequest,
) (*model.User, error) {

	if _, err := s.userRepo.FindByUsername(
		ctx,
		req.Username,
	); err == nil {
		return nil, errors.New("username already exists")
	}

	if _, err := s.userRepo.FindByEmail(
		ctx,
		req.Email,
	); err == nil {
		return nil, errors.New("email already exists")
	}

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		UUID:      uuid.New().String(),
		Username:  req.Username,
		Name:      req.Name,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  password,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Me(
	ctx context.Context,
	userID uint64,
) (*model.User, error) {

	user, err := s.userRepo.FindByID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	return user, nil
}

func (s *AuthService) ChangePassword(
	ctx context.Context,
	userID uint64,
	req request.ChangePasswordRequest,
) error {

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := utils.CheckPassword(
		user.Password,
		req.CurrentPassword,
	); err != nil {
		return errors.New("current password is incorrect")
	}

	if req.CurrentPassword == req.NewPassword {
		return errors.New(
			"new password must be different from current password",
		)
	}

	hashedPassword, err := utils.HashPassword(
		req.NewPassword,
	)
	if err != nil {
		return errors.New(
			"failed to hash password",
		)
	}

	if err := s.userRepo.UpdatePassword(
		ctx,
		userID,
		hashedPassword,
	); err != nil {
		return errors.New(
			"failed to update password",
		)
	}

	return nil
}

func (s *AuthService) RefreshToken(
	tokenString string,
) (string, error) {

	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		return "", errors.New("invalid or expired token")
	}

	return utils.GenerateToken(
		claims.UserID,
		claims.Username,
		claims.Roles,
	)
}
