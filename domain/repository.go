package domain

import (
	"context"

	"github.com/google/uuid"
)

// BaseReader defines common read operations for any entity
// Following OCP: new entity types can implement this interface without modifying existing code
type BaseReader[T any] interface {
	Get(ctx context.Context, id uuid.UUID) (*T, error)
	List(ctx context.Context, limit, offset int) ([]*T, error)
}

// BaseWriter defines common write operations for any entity
// Following OCP: new entity types can implement this interface without modifying existing code
type BaseWriter[T any] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// BaseRepository combines read and write operations for any entity
// Following OCP: provides a generic pattern that can be extended
type BaseRepository[T any] interface {
	BaseReader[T]
	BaseWriter[T]
}

// UserReader defines read operations for users
// Following ISP: clients that only need to read users don't depend on write operations
type UserReader interface {
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// UserWriter defines write operations for users
// Following ISP: clients that only need to write users don't depend on read operations
type UserWriter interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// UserRepository combines read and write operations
// Following ISP: full repository interface for clients that need both operations
type UserRepository interface {
	UserReader
	UserWriter
}

// ItemReader defines read operations for items
type ItemReader interface {
	GetItem(ctx context.Context, id uuid.UUID) (*Item, error)
	ListItems(ctx context.Context, limit, offset int) ([]*Item, error)
}

// ItemWriter defines write operations for items
type ItemWriter interface {
	CreateItem(ctx context.Context, item *Item) (*Item, error)
	UpdateItem(ctx context.Context, item *Item) (*Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

// ItemRepository combines read and write operations for items
type ItemRepository interface {
	ItemReader
	ItemWriter
}
