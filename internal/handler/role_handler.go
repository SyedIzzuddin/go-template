package handler

import (
	"strconv"

	"go-template/internal/dto"
	"go-template/internal/logger"
	"go-template/internal/service"
	"go-template/pkg/response"
	"go-template/pkg/roles"
	"go-template/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RoleHandler struct {
	roleService service.RoleService
	validator   *validator.Validator
}

func NewRoleHandler(roleService service.RoleService, validator *validator.Validator) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		validator:   validator,
	}
}

func (h *RoleHandler) UpdateUserRole(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("UpdateUserRole request started", zap.String("request_id", requestID))

	// Get user ID from path parameter
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid user ID", err.Error())
	}

	// Bind request
	var req dto.UpdateUserRoleRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind role update request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Role update validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	// Update user role
	updatedUser, err := h.roleService.UpdateUserRole(c.Request().Context(), userID, req.Role)
	if err != nil {
		logger.Error("Failed to update user role", zap.Error(err), zap.String("request_id", requestID))
		
		switch err {
		case service.ErrInvalidRole:
			return response.BadRequest(c, "Invalid role", err.Error())
		case service.ErrSameRole:
			return response.BadRequest(c, "User already has this role", err.Error())
		default:
			return response.InternalServerError(c, "Failed to update user role", err.Error())
		}
	}

	logger.Info("UpdateUserRole request completed", zap.String("request_id", requestID))
	return response.Success(c, "User role updated successfully", updatedUser)
}

func (h *RoleHandler) GetUsersByRole(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("GetUsersByRole request started", zap.String("request_id", requestID))

	// Get role from query parameter
	role := c.QueryParam("role")
	if role == "" {
		logger.Error("Missing role query parameter", zap.String("request_id", requestID))
		return response.BadRequest(c, "Role query parameter is required", nil)
	}

	// Get users by role
	users, err := h.roleService.GetUsersByRole(c.Request().Context(), role)
	if err != nil {
		logger.Error("Failed to get users by role", zap.Error(err), zap.String("request_id", requestID))
		
		switch err {
		case service.ErrInvalidRole:
			return response.BadRequest(c, "Invalid role", err.Error())
		default:
			return response.InternalServerError(c, "Failed to get users by role", err.Error())
		}
	}

	logger.Info("GetUsersByRole request completed", zap.String("request_id", requestID))
	return response.Success(c, "Users retrieved successfully", users)
}

func (h *RoleHandler) GetAvailableRoles(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("GetAvailableRoles request started", zap.String("request_id", requestID))

	availableRoles := roles.GetAllRoles()
	
	logger.Info("GetAvailableRoles request completed", zap.String("request_id", requestID))
	return response.Success(c, "Available roles retrieved successfully", map[string]interface{}{
		"roles": availableRoles,
	})
}