package user

import (
	"context"

	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/SoraDaibu/go-clean-starter/internal/repository"
	"github.com/SoraDaibu/go-clean-starter/internal/repository/common"
	"github.com/SoraDaibu/go-clean-starter/internal/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userRepository implements domain.UserRepository
// Following DIP: depends on abstractions (domain interfaces) not concrete implementations
// Following composition: uses BaseRepository for common functionality
type userRepository struct {
	*repository.BaseRepository
}

// NewUserRepository creates a new user repository implementation
// Following DIP: returns domain interface, not concrete type
func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{
		BaseRepository: repository.NewBaseRepository(pool),
	}
}

// GetUser implements domain.UserReader
func (r *userRepository) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, err := r.GetQueries(ctx).GetUser(ctx, common.UUIDToPgtype(id))
	if err != nil {
		return nil, err
	}

	userID, err := common.PgtypeToUUID(u.ID)
	if err != nil {
		return nil, err
	}

	return domain.UserFromSource(userID, u.Name, u.Email), nil
}

// ListUsers implements domain.UserReader
// Note: The current sqlc query doesn't support limit/offset, so we apply manual pagination
func (r *userRepository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := r.GetQueries(ctx).ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	// Apply manual pagination since sqlc query doesn't support it
	start := offset
	end := offset + limit
	if start >= len(users) {
		return []*domain.User{}, nil
	}
	if end > len(users) {
		end = len(users)
	}

	result := make([]*domain.User, end-start)
	for i, u := range users[start:end] {
		userID, err := common.PgtypeToUUID(u.ID)
		if err != nil {
			return nil, err
		}
		result[i] = domain.UserFromSource(userID, u.Name, u.Email)
	}

	return result, nil
}

// GetUserByEmail implements domain.UserReader
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.GetQueries(ctx).GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	userID, err := common.PgtypeToUUID(u.ID)
	if err != nil {
		return nil, err
	}

	return domain.UserFromSource(userID, u.Name, u.Email), nil
}

// CreateUser implements domain.UserWriter
func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	u, err := r.GetQueries(ctx).CreateUser(ctx, sqlc.CreateUserParams{
		ID:       common.UUIDToPgtype(user.ID()),
		Name:     user.Name(),
		Email:    user.Email(),
		Password: string(user.Password()),
	})

	if err != nil {
		return nil, err
	}

	userID, err := common.PgtypeToUUID(u.ID)
	if err != nil {
		return nil, err
	}

	return domain.UserFromSource(userID, u.Name, u.Email), nil
}

// UpdateUser implements domain.UserWriter
func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	u, err := r.GetQueries(ctx).UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:   common.UUIDToPgtype(user.ID()),
		Name: user.Name(),
	})

	if err != nil {
		return nil, err
	}

	userID, err := common.PgtypeToUUID(u.ID)
	if err != nil {
		return nil, err
	}

	return domain.UserFromSource(userID, u.Name, u.Email), nil
}

// DeleteUser implements domain.UserWriter
func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.GetQueries(ctx).DeleteUser(ctx, common.UUIDToPgtype(id))
}
