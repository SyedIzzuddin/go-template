package middleware

import (
	"go-template/internal/logger"
	"go-template/pkg/response"
	"go-template/pkg/roles"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RequireRoles creates middleware that checks if user has one of the required roles
func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			
			// Get user role from context (set by auth middleware)
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				logger.Warn("User role not found in context", zap.String("request_id", requestID))
				return response.Forbidden(c, "Access denied")
			}

			// Check if user role is in allowed roles
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					logger.Debug("Role access granted", 
						zap.String("user_role", userRole),
						zap.String("required_roles", joinRoles(allowedRoles)),
						zap.String("request_id", requestID))
					return next(c)
				}
			}

			logger.Warn("Access denied - insufficient role", 
				zap.String("user_role", userRole),
				zap.String("required_roles", joinRoles(allowedRoles)),
				zap.String("request_id", requestID))
			
			return response.Forbidden(c, "Insufficient permissions")
		}
	}
}

// RequireAdmin creates middleware that requires admin role
func RequireAdmin() echo.MiddlewareFunc {
	return RequireRoles(roles.Admin)
}

// RequireModeratorOrAdmin creates middleware that requires moderator or admin role
func RequireModeratorOrAdmin() echo.MiddlewareFunc {
	return RequireRoles(roles.Admin, roles.Moderator)
}

// Helper function to join roles for logging
func joinRoles(roles []string) string {
	if len(roles) == 0 {
		return ""
	}
	
	result := roles[0]
	for i := 1; i < len(roles); i++ {
		result += ", " + roles[i]
	}
	return result
}