package repository

import (
	"context"
	"database/sql"

	db "go-template/db/sqlc"

	"go-template/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, name, email string) (*entity.User, error)
	CreateWithPassword(ctx context.Context, name, email, passwordHash, role string) (*entity.User, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByEmailWithPassword(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, id int, name string) (*entity.User, error)
	UpdateRole(ctx context.Context, id int, role string) (*entity.User, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByRole(ctx context.Context, role string) ([]entity.User, error)
}

type userRepository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewUserRepository(dbConn *sql.DB) UserRepository {
	return &userRepository{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

func (r *userRepository) Create(ctx context.Context, name, email string) (*entity.User, error) {
	createdUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           int(createdUser.ID),
		Name:         createdUser.Name,
		Email:        createdUser.Email,
		PasswordHash: createdUser.PasswordHash,
		Role:         createdUser.Role,
		CreatedAt:    createdUser.CreatedAt.Time,
		UpdatedAt:    createdUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := r.queries.GetUser(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           int(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           int(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, id int, name string) (*entity.User, error) {
	updatedUser, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:   int32(id),
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           int(updatedUser.ID),
		Name:         updatedUser.Name,
		Email:        updatedUser.Email,
		PasswordHash: updatedUser.PasswordHash,
		Role:         updatedUser.Role,
		CreatedAt:    updatedUser.CreatedAt.Time,
		UpdatedAt:    updatedUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteUser(ctx, int32(id))
}

func (r *userRepository) CreateWithPassword(ctx context.Context, name, email, passwordHash, role string) (*entity.User, error) {
	createdUser, err := r.queries.CreateUserWithPassword(ctx, db.CreateUserWithPasswordParams{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           int(createdUser.ID),
		Name:         createdUser.Name,
		Email:        createdUser.Email,
		PasswordHash: createdUser.PasswordHash,
		Role:         createdUser.Role,
		CreatedAt:    createdUser.CreatedAt.Time,
		UpdatedAt:    createdUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByEmailWithPassword(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmailWithPassword(ctx, email)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           int(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	// Get all users without pagination
	userList, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, len(userList))
	for i, dbUser := range userList {
		users[i] = entity.User{
			ID:           int(dbUser.ID),
			Name:         dbUser.Name,
			Email:        dbUser.Email,
			PasswordHash: dbUser.PasswordHash,
			Role:         dbUser.Role,
			CreatedAt:    dbUser.CreatedAt.Time,
			UpdatedAt:    dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}

func (r *userRepository) UpdateRole(ctx context.Context, id int, role string) (*entity.User, error) {
	updatedUser, err := r.queries.UpdateUserRole(ctx, db.UpdateUserRoleParams{
		ID:   int32(id),
		Role: role,
	})
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           int(updatedUser.ID),
		Name:         updatedUser.Name,
		Email:        updatedUser.Email,
		PasswordHash: updatedUser.PasswordHash,
		Role:         updatedUser.Role,
		CreatedAt:    updatedUser.CreatedAt.Time,
		UpdatedAt:    updatedUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByRole(ctx context.Context, role string) ([]entity.User, error) {
	userList, err := r.queries.GetUsersByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, len(userList))
	for i, dbUser := range userList {
		users[i] = entity.User{
			ID:           int(dbUser.ID),
			Name:         dbUser.Name,
			Email:        dbUser.Email,
			PasswordHash: dbUser.PasswordHash,
			Role:         dbUser.Role,
			CreatedAt:    dbUser.CreatedAt.Time,
			UpdatedAt:    dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}
