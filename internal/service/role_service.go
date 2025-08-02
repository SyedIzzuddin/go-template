package service

import (
	"context"
	"database/sql"
	"errors"
	"go-template/internal/dto"
	"go-template/internal/logger"
	"go-template/internal/repository"
	"go-template/pkg/roles"

	"go.uber.org/zap"
)

var (
	ErrInvalidRole = errors.New("invalid role")
	ErrSameRole    = errors.New("user already has this role")
)

type RoleService interface {
	UpdateUserRole(ctx context.Context, userID int, newRole string) (*dto.UserResponse, error)
	GetUsersByRole(ctx context.Context, role string) ([]dto.UserResponse, error)
}

type roleService struct {
	userRepo repository.UserRepository
}

func NewRoleService(userRepo repository.UserRepository) RoleService {
	return &roleService{
		userRepo: userRepo,
	}
}

func (s *roleService) UpdateUserRole(ctx context.Context, userID int, newRole string) (*dto.UserResponse, error) {
	logger.Info("Updating user role", zap.Int("user_id", userID), zap.String("new_role", newRole))

	// Validate role
	if !roles.IsValidRole(newRole) {
		logger.Warn("Invalid role provided", zap.String("role", newRole))
		return nil, ErrInvalidRole
	}

	// Get current user to check existing role
	currentUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("User not found for role update", zap.Int("user_id", userID))
			return nil, errors.New("user not found")
		}
		logger.Error("Failed to get user for role update", zap.Error(err))
		return nil, errors.New("failed to get user")
	}

	// Check if user already has this role
	if currentUser.Role == newRole {
		logger.Warn("User already has the specified role", zap.Int("user_id", userID), zap.String("role", newRole))
		return nil, ErrSameRole
	}

	// Update user role
	updatedUser, err := s.userRepo.UpdateRole(ctx, userID, newRole)
	if err != nil {
		logger.Error("Failed to update user role", zap.Error(err))
		return nil, errors.New("failed to update user role")
	}

	logger.Info("User role updated successfully", zap.Int("user_id", userID), zap.String("old_role", currentUser.Role), zap.String("new_role", newRole))

	return &dto.UserResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Role:      updatedUser.Role,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

func (s *roleService) GetUsersByRole(ctx context.Context, role string) ([]dto.UserResponse, error) {
	logger.Debug("Getting users by role", zap.String("role", role))

	// Validate role
	if !roles.IsValidRole(role) {
		logger.Warn("Invalid role provided for user lookup", zap.String("role", role))
		return nil, ErrInvalidRole
	}

	users, err := s.userRepo.GetByRole(ctx, role)
	if err != nil {
		logger.Error("Failed to get users by role", zap.Error(err))
		return nil, errors.New("failed to get users by role")
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	logger.Debug("Retrieved users by role", zap.String("role", role), zap.Int("count", len(userResponses)))

	return userResponses, nil
}