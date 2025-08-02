package router

import (
	"net/http"
	"time"

	"go-template/internal/database"
	"go-template/internal/handler"
	"go-template/internal/middleware"
	"go-template/internal/repository"
	"go-template/pkg/jwt"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *database.DB, userHandler *handler.UserHandler, fileHandler *handler.FileHandler, authHandler *handler.AuthHandler, jwtManager *jwt.JWTManager, userRepo repository.UserRepository, roleHandler *handler.RoleHandler) {
	api := e.Group("/api/v1")

	// Health check (public)
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"database":  db.Health(),
			"timestamp": time.Now().Unix(),
		})
	})

	// Authentication routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	
	// Protected auth routes
	authProtected := auth.Group("", middleware.AuthMiddleware(jwtManager, userRepo))
	authProtected.GET("/me", authHandler.GetProfile)

	// Protected user routes (most require admin access)
	users := api.Group("/users", middleware.AuthMiddleware(jwtManager, userRepo))
	users.POST("", userHandler.CreateUser, middleware.RequireAdmin())
	users.GET("", userHandler.GetAllUsers, middleware.RequireModeratorOrAdmin())
	users.GET("/:id", userHandler.GetUser, middleware.RequireModeratorOrAdmin())
	users.PUT("/:id", userHandler.UpdateUser, middleware.RequireAdmin())
	users.DELETE("/:id", userHandler.DeleteUser, middleware.RequireAdmin())

	// Role management routes (admin only)
	roles := api.Group("/roles", middleware.AuthMiddleware(jwtManager, userRepo), middleware.RequireAdmin())
	roles.PUT("/users/:id/role", roleHandler.UpdateUserRole)
	roles.GET("/users", roleHandler.GetUsersByRole) // Query param: ?role=admin|moderator|user
	roles.GET("/available", roleHandler.GetAvailableRoles)

	// Protected file routes
	files := api.Group("/files", middleware.AuthMiddleware(jwtManager, userRepo))
	files.POST("/upload", fileHandler.UploadFile) // All authenticated users can upload
	files.GET("", fileHandler.GetAllFiles, middleware.RequireModeratorOrAdmin()) // Only mods/admins see all files
	files.GET("/my", fileHandler.GetMyFiles) // Users can see their own files
	files.GET("/:id", fileHandler.GetFile) // Users can see file details (we'll handle ownership in handler)
	files.PUT("/:id", fileHandler.UpdateFile) // Users can update their own files (we'll handle ownership in handler)
	files.DELETE("/:id", fileHandler.DeleteFile) // Users can delete their own files (we'll handle ownership in handler)
	files.GET("/:id/download", fileHandler.DownloadFile) // Users can download files (we'll handle access in handler)

	// Static file serving (public)
	e.Static("/uploads", "uploads")
	e.GET("/files/:filename", fileHandler.ServeFile)
}
